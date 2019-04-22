package main

import (
	"flag"
	"strconv"
	"time"

	"github.com/perlin-network/noise/cipher/aead"
	"github.com/perlin-network/noise/handshake/ecdh"

	"github.com/perlin-network/noise/skademlia"

	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/log"
	"github.com/perlin-network/noise/protocol"
)

/** Define Message **/
var (
	opcodeBlock noise.Opcode
	_           noise.Message = (*Block)(nil)
)

func main() {
	container, err := InitializeContainer()

	if err != nil {
		panic(err)
	}

	container.db.Close()

	portNumber := flag.Uint("p", 3000, "port to listen to peer")
	flag.Parse()

	// Instantiate a default set of node parameters.
	params := noise.DefaultParams()
	params.Port = uint16(*portNumber)
	params.Keys = skademlia.RandomKeys()

	// Instantiate a new node that listens for peers on portNumber.
	node, err := noise.NewNode(params)
	if err != nil {
		panic(err)
	}

	defer node.Kill()

	p := protocol.New()
	p.Register(ecdh.New())
	p.Register(aead.New())
	p.Register(skademlia.New())
	p.Enforce(node)

	// Set up the node for listening
	setup(node)

	// Start listening for incoming peers.
	go node.Listen()

	log.Info().Msgf("Listening for peers on port %d.", node.ExternalPort())

	if len(flag.Args()) > 0 {
		for _, address := range flag.Args() {

			attempts := 0
			for {

				attempts++

				if attempts > 5 {
					log.Fatal().Msg("Connection timeout on " + address)
				}

				peer, err := node.Dial(address)
				if err != nil {
					time.Sleep(5 * time.Second)
					continue
				}

				skademlia.WaitUntilAuthenticated(peer)
				break

			}
		}

		peers := skademlia.FindNode(node, protocol.NodeID(node).(skademlia.ID), skademlia.BucketSize(), 8)
		log.Info().Msgf("Bootstrapped with peers: %+v", peers)
	}

	node.Fence()

}

func sendTest() {
	time.Sleep(20 * time.Second)

	// reader := bufio.NewReader(os.Stdin)
	// for {

	// 	_, err := reader.ReadString('\n')

	// 	if err != nil && err.Error() != "EOF" {
	// 		panic(err)
	// 	}

	// 	block := Block{
	// 		Amount:       1000,
	// 		From:         "Trace",
	// 		Hash:         "laksjdflkadn",
	// 		Index:        1,
	// 		PreviousHash: "pwjeanalcoiaen123",
	// 		Timestamp:    time.Now().String(),
	// 		To:           "You"}

	// 	skademlia.BroadcastAsync(node, block)
	// }
}

func setup(node *noise.Node) {
	opcodeBlock = noise.RegisterMessage(noise.Opcode(16), (*Block)(nil))
	node.OnPeerInit(func(node *noise.Node, peer *noise.Peer) error {
		peer.OnConnError(func(node *noise.Node, peer *noise.Peer, err error) error {
			log.Info().Msgf("Got an error: %v", err)

			return nil
		})

		peer.OnDisconnect(func(node *noise.Node, peer *noise.Peer) error {
			ip := peer.RemoteIP().String()
			port := strconv.Itoa(int(peer.RemotePort()))
			log.Info().Msgf("Peer %v has disconnected.", ip+":"+port)

			return nil
		})

		go func() {
			for {
				select {
				case msg := <-peer.Receive(opcodeBlock):
					log.Info().Msgf("[%s] : %s", protocol.PeerID(peer), msg)
				case <-time.After(10 * time.Second):
					log.Debug().Msg("Killing the node")
					node.Kill()
				}
			}
		}()

		return nil
	})
}

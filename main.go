package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

type Message struct {
	data string
}

func main() {
	container, err := InitializeContainer()

	if err != nil {
		panic(err)
	}

	container.db.Close()

	portNumber := flag.Uint("p", 3000, "port to listen to peer")
	host := flag.Bool("h", false, "Whether to make this a host or not")

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

	if *host {
		r := mux.NewRouter()

		r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			data := Block{}
			bytes, err := ioutil.ReadAll(request.Body)

			defer request.Body.Close()
			fmt.Println(string(bytes))
			if err != nil {
				log.Fatal().Msg(err.Error())
			}
			json.Unmarshal(bytes, &data)
			fmt.Println(data)

			skademlia.BroadcastAsync(node, data)

		}).Methods("POST")

		http.ListenAndServe(":8080", r)
	}

	node.Fence()

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
				}
			}
		}()

		return nil
	})
}

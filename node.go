package main

import (
	"os"
	"strconv"

	"github.com/perlin-network/noise"
)

func makeNode() {

	// Get the port number
	port := os.Args[1]
	portNumber, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		panic(err)
	}

	// Instantiate a default set of node parameters.
	params := noise.DefaultParams()
	params.Port = uint16(portNumber)

	// Instantiate a new node that listens for peers on portNumber.
	node, err := noise.NewNode(params)
	if err != nil {
		panic(err)
	}

	// Start listening for incoming peers.
	go node.Listen()

	// Dial peer at address 127.0.0.1:3001.
	peer, err := node.Dial("127.0.0.1:3001")
	if err != nil {
		panic("failed to dial peer located at 127.0.0.1:3001!")
	}

	message := BlockMessage{}
	message.block = []byte("trace")
	mistake := <-peer.SendMessageAsync(message)
	if mistake != nil {
		panic(mistake)
	}

	select {}
}

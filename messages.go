package main

import (
	"encoding/json"
	"fmt"

	"github.com/bobg/scp"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/payload"
)

type NominateMessage struct {
	Counter   int
	MessageID int
}

type Payload struct {
	X int32
	Y int32
}

func (message NominateMessage) Read(reader payload.Reader) (noise.Message, error) {
	bytes, err := reader.ReadBytes()
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &message)
	if err != nil {
		panic(err)
	}
	return message, err
}
func (message NominateMessage) Write() []byte {
	fmt.Print(message)
	bytes, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	return payload.NewWriter(nil).WriteBytes(bytes).Bytes()
}

type NominatePrepareTopic scp.Msg
type PrepareTopic scp.PrepTopic

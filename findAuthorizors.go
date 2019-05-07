package main

import (
	"encoding/json"

	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/payload"
)

type FindAuthorizors struct {
	Total     float64
	Remaining float64
}

func (findAuths FindAuthorizors) Read(reader payload.Reader) (noise.Message, error) {
	bytes, err := reader.ReadBytes()
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &findAuths)
	if err != nil {
		panic(err)
	}

	return findAuths, err
}

func (findAuths FindAuthorizors) Write() []byte {
	bytes, err := json.Marshal(findAuths)
	if err != nil {
		panic(err)
	}
	return payload.NewWriter(nil).WriteBytes(bytes).Bytes()
}

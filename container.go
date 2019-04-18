package main

import (
	"encoding/json"
)

// Container represent an IoC Container
type Container struct {
	db IDatastore
}

func ProvideContainer(db *Db) *Container {

	return &Container{db: db}
}

func (container *Container) writeBlock(block *Block) {
	bytes, err := json.Marshal(&block)
	err = container.db.Put([]byte(block.Hash), bytes)
	if err != nil {
		panic(err)
	}
}

func (container *Container) readBlock(key string) *Block {
	blockBytes, _ := container.db.Get([]byte(key))

	var newBlock Block
	json.Unmarshal(*blockBytes, &newBlock)

	return &newBlock
}

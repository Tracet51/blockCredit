package main

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
)

func ProvideLevelDb() *leveldb.DB {

	db, err := leveldb.OpenFile("./data/", nil)
	if err != nil {
		panic(err)
	}
	return db
}

func ProvideDb(db *leveldb.DB) *Db {
	return &Db{db}
}

// IDatastore is an interface that wraps the levelDB api to provide IoC
type IDatastore interface {
	WriteBlock(block *Block) error
	ReadBlock(key string) (*Block, error)
	Close() error
}

// Db implements the IDatastore Interface
type Db struct {
	db *leveldb.DB
}

// WriteBlock write the block to the datastore
func (store Db) WriteBlock(block *Block) error {
	bytes, err := json.Marshal(&block)
	err = store.db.Put([]byte(block.Hash), bytes, nil)
	return err
}

// ReadBlock returns a pointer to a Block object with the given key
func (store Db) ReadBlock(key string) (*Block, error) {

	var newBlock Block
	blockBytes, err := store.db.Get([]byte(key), nil)
	if err != nil {
		return &newBlock, err
	}
	err = json.Unmarshal(blockBytes, &newBlock)

	return &newBlock, err
}

func (store Db) Close() error {
	err := store.db.Close()

	return err
}

package main

import (
	"github.com/syndtr/goleveldb/leveldb"
)

func LoadDb() (IDatastore, error) {
	db, err := leveldb.OpenFile("./data/", nil)
	if err != nil {
		panic(err)
	}
	return &Db{db}, nil
}

func ProvideLevelDb() *leveldb.DB {

	db, err := leveldb.OpenFile("./data/", nil)
	if err != nil {
		panic(err)
	}

	return db
}

func ProvideDb(db *leveldb.DB) Db {
	return Db{db}
}

// IDatastore is an interface that wraps the levelDB api to provide IoC
type IDatastore interface {
	Put([]byte, []byte) error
	Get([]byte) (*[]byte, error)
	Close()
}

// Db the IDatastore Interface
type Db struct {
	db *leveldb.DB
}

// Put stores the value with key
func (store *Db) Put(key []byte, value []byte) error {
	store.db.Put(key, value, nil)
	return nil
}

// Get gets the value with given key
func (store *Db) Get(value []byte) (*[]byte, error) {
	block, err := store.db.Get(value, nil)
	return &block, err
}

// Close closes the connection to the database
func (store *Db) Close() {
	store.db.Close()
}

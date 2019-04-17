//+build wireinject

package main

import (
	"github.com/google/wire"
)

func InitializeDB() (Db, error) {

	wire.Build(ProvideDb, ProvideLevelDb)

	return Db{}, nil
}

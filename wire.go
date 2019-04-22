//+build wireinject

package main

import (
	"github.com/google/wire"
)

func InitializeContainer() (*Container, error) {

	// _ := wire.NewSet(Db{}, wire.Bind(new(IDatastore), new(Db)))
	wire.Build(ProvideContainer, ProvideDb, ProvideLevelDb)
	return &Container{}, nil
}

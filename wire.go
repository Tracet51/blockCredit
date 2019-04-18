//+build wireinject

package main

import (
	"github.com/google/wire"
)

func InitializeContainer() (*Container, error) {

	wire.Build(ProvideContainer, ProvideDb, ProvideLevelDb)
	return &Container{}, nil
}

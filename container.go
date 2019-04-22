package main

// Container represent an IoC Container
type Container struct {
	db IDatastore
}

func ProvideContainer(db *Db) *Container {

	return &Container{db: db}
}

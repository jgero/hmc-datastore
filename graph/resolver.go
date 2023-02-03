package graph

import "github.com/jgero/hmc-datastore/graph/repository"

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	Repo repository.Repository
}

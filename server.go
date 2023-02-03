package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jgero/hmc-datastore/graph"
	"github.com/jgero/hmc-datastore/graph/repository"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	ctx, cancel := context.WithCancel(context.Background())
	repo := repository.NewNeo4jRepo("neo4j://localhost:7687", "neo4j", "neo4jneo4j")
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Repo: repo,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	repo.Close(ctx)
	cancel()
}

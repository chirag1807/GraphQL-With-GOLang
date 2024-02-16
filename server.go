package main

import (
	// "articlewithgraphql/api/middleware"
	"articlewithgraphql/api/middleware"
	"articlewithgraphql/config"
	"articlewithgraphql/constants"
	"articlewithgraphql/db"
	"articlewithgraphql/graph"
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	config.LoadEnv()

	conn, err := db.DBConnection()
	if err != nil {
		return
	}
	defer conn.Close(context.Background())

	resolverConfig := graph.Config{Resolvers: &graph.Resolver{}}

	var srv http.Handler = handler.NewDefaultServer(graph.NewExecutableSchema(resolverConfig))

	middleware := middleware.SetDBConnection(conn)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", middleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", constants.PORT_NO)
	log.Fatal(http.ListenAndServe(constants.PORT_NO, nil))
}

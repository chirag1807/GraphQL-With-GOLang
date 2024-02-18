package main

import (
	"articlewithgraphql/api/middleware"
	"articlewithgraphql/config"
	"articlewithgraphql/constants"
	"articlewithgraphql/dataloader"
	"articlewithgraphql/db"
	"articlewithgraphql/graph"
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.LoadEnv()

	conn, err := db.DBConnection()
	if err != nil {
		return
	}
	defer conn.Close(context.Background())

	router := chi.NewRouter()
	router.Use(middleware.SetDBConnection(conn))

	resolver := &graph.Resolver{
		DB: conn,
	}

	resolverConfig := graph.Config{Resolvers: resolver}

	resolverConfig.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		ctx, err := middleware.AuthenticateUser(ctx)
		if err != nil {
			return "", err
		}
		return next(ctx)
	}
 
	resolverConfig.Directives.IsAdmin = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		err := middleware.AuthorizeAdmin(ctx)
		if err != nil {
			return "", err
		}
		return next(ctx)
	}

	var srv http.Handler = handler.NewDefaultServer(graph.NewExecutableSchema(resolverConfig))
	srv = dataloader.DataLoaderMiddleware(conn, srv)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", constants.PORT_NO)
	log.Fatal(http.ListenAndServe(constants.PORT_NO, router))
}

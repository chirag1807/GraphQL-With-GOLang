package db

import (
	"articlewithgraphql/config"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DBConnection() (conn *pgxpool.Pool, err error) {
	DATABASE_URL := "postgresql://" + config.DatabaseConfig.DATABASE_USERNAME + ":" + config.DatabaseConfig.DATABASE_PASSWORD + "@127.0.0.1:" + config.DatabaseConfig.DATABASE_PORT + "/" + config.DatabaseConfig.DATABASE_NAME + "?sslmode=" + config.DatabaseConfig.DATABASE_SSLMODE
	conn, err = pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn, nil
}

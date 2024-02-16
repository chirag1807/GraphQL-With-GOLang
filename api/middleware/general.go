package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func SetDBConnection(conn *pgx.Conn) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("1")
			token := r.Header.Get("Authorization")
			ctx := context.WithValue(r.Context(), "conn", conn)
			ctx = context.WithValue(ctx, "token", token)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthenticateUser() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Middleware: Before handling the request")

			// Call the next handler in the chain
			handler.ServeHTTP(w, r)

			fmt.Println("Middleware: After handling the request")
		})
	}
}

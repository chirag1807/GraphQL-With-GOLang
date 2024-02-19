package middleware

import (
	"articlewithgraphql/error"
	"articlewithgraphql/utils"
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetDBConnection(conn *pgxpool.Pool) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			ctx := context.WithValue(r.Context(), "conn", conn)
			ctx = context.WithValue(ctx, "token", token)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthenticateUser(ctx context.Context) (context.Context, error) {
	token := ctx.Value("token").(string)
	if token != "" {
		id, isadmin, err := utils.VerifyToken(token)
		if err != nil {
			return ctx, errorhandling.AccessTokenExpired
		}
		ctx = context.WithValue(ctx, "id", id)
		ctx = context.WithValue(ctx, "isadmin", isadmin)
		return ctx, nil
	} else {
		return ctx, errorhandling.TokenNotFound
	}
}

func AuthorizeAdmin(ctx context.Context) error {
	if isadmin, ok := ctx.Value("isadmin").(bool); ok {
		if isadmin {
			return nil
		} else {
			return errorhandling.UnauthorizedError
		}
	} else {
		return errorhandling.TokenNotFound
	}
}

package utils

// import (
// 	"articlewithgraphql/error"
// 	"context"
// 	"fmt"

// 	"github.com/99designs/gqlgen/graphql"
// 	"github.com/vektah/gqlparser/v2/gqlerror"
// )

// func ErrorGenerator(err error) (string, error) {
// 	if error, ok := err.(errorhandling.CustomError); ok {
// 		return "", errorhandling.CustomError{
// 			Message:    error.Message,
// 			StatusCode: error.StatusCode,
// 		}
// 	} else {
// 		return "", errorhandling.CustomError{
// 			Message:    "Internal Server Error",
// 			StatusCode: 500,
// 		}
// 	}
// }

package errorhandling

import (
	"net/http"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

var (
	EmailvalidationError    = CreateCustomError("Email Validation Failed, Please Provide Valid Email.", http.StatusBadRequest)
	DuplicateEmailFound     = CreateCustomError("Duplicate Email Found.", http.StatusConflict)
	RegistrationFailedError = CreateCustomError("User Registration Failed.", http.StatusInternalServerError)
	LoginFailedError        = CreateCustomError("User Login Failed.", http.StatusUnauthorized)
	AccessTokenExpired      = CreateCustomError("Access Token is Expired, Please Regenrate It.", http.StatusUnauthorized)
	UnauthorizedError       = CreateCustomError("You are Not Authorized to Perform this Action.", http.StatusUnauthorized)
	NoUserFound             = CreateCustomError("No User Found for This Request.", http.StatusNotFound)
	TokenNotFound           = CreateCustomError("Access Token Not Found.", http.StatusUnauthorized)
	PasswordNotMatch        = CreateCustomError("Password is Incorrect.", http.StatusUnauthorized)
)

func CreateCustomError(message string, statusCode int) *gqlerror.Error {
	return &gqlerror.Error{
			Message: message,
			Extensions: map[string]interface{}{
				"StatusCode": statusCode,
			},
		}
}

package http

import (
	"fmt"
	"net/http"
)

// ErrorType is the base type for http error
type ErrorType interface {
	Error() string
}

type errorType struct{}

func (errorType) Error() string {
	return ""
}

// Error is a base error to facilitate `errors.Is(xxx, http.Error)`
var Error = errorType{}

// ClientError facilitate `errors.Is(xxx, http.ClientError)`
var ClientError = fmt.Errorf("%w", Error)

// ServerError facilitate `errors.Is(xxx, http.ServerError)`
var ServerError = fmt.Errorf("%w", Error)

var NotFoundError = fmt.Errorf("%w%s", ClientError, http.StatusText(http.StatusNotFound))

// errorFromStatusCode returns an error of the corresponding status code that we maybe interested in checking later using `errors.Is(err, http.XXXError)`
func errorFromStatusCode(code int) ErrorType {
	// TODO: use go generate to create errors for all http status and return them here, for now we just manually create them where needed
	switch true {
	case code == http.StatusNotFound:
		return NotFoundError

	case code > http.StatusInternalServerError:
		return ServerError
	case code > http.StatusBadRequest:
		return ClientError
	default:
		return nil
	}
}

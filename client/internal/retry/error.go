package retry

import (
	"fmt"
)

// ErrorType is the base error interface of the retry package
type ErrorType interface {
	Error() string
}

// errorType is the implementing class
type errorType struct{}

func (err errorType) Error() string {
	return "retry error"
}

// Error facilitates `errors.Is(err, retry.Error)`
var Error = errorType{}

// TimeoutError facilitates `errors.Is(err, retry.TimeoutError)`
var TimeoutError = fmt.Errorf("%w: timeout", Error)

// newTimeoutErrorf is used within this package to create new timeout error
func newTimeoutErrorf(format string, a ...any) ErrorType {
	return fmt.Errorf("%w: %w", TimeoutError, fmt.Errorf(format, a...))
}

// RetriesExceededError facilitates `errors.Is(err, retry.RetriesExceededError)`
var RetriesExceededError = fmt.Errorf("%w: retries exceeded", Error)

// newRetriesExceededErrorf is used within this package to create new RetriesExceededErrorType error
func newRetriesExceededErrorf(format string, a ...any) ErrorType {
	return fmt.Errorf("%w: %w", RetriesExceededError, fmt.Errorf(format, a...))
}

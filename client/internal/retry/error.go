package retry

import (
	"fmt"
)

// ErrorType is the base error interface of the retry package.
type ErrorType interface {
	Error() string
}

// errorType is the implementing class.
type errorType struct{}

func (err errorType) Error() string {
	return "retry error"
}

// Error facilitates `errors.Is(err, retry.Error)`.
var Error = errorType{}

// TimeoutError facilitates `errors.Is(err, retry.TimeoutError)`.
var TimeoutError = fmt.Errorf("%w: timeout", Error)

// newTimeoutErrorf is used within this package to create new TimeoutError.
func newTimeoutErrorf(format string, a ...any) ErrorType {
	return fmt.Errorf("%w: %w", TimeoutError, fmt.Errorf(format, a...))
}

// RetriesExceededError facilitates `errors.Is(err, retry.RetriesExceededError)`.
var RetriesExceededError = fmt.Errorf("%w: retries exceeded", Error)

// newRetriesExceededErrorf is used within this package to create new RetriesExceededError.
func newRetriesExceededErrorf(format string, a ...any) ErrorType {
	return fmt.Errorf("%w: %w", RetriesExceededError, fmt.Errorf(format, a...))
}

// ContextCancelledError facilitates `errors.Is(err, retry.RetriesExceededError)`.
var ContextCancelledError = fmt.Errorf("%w: context cancelled", Error)

// newContextCancelledErrorf is used within this package to create new RetriesExceededError.
func newContextCancelledErrorf(format string, a ...any) ErrorType {
	return fmt.Errorf("%w: %w", ContextCancelledError, fmt.Errorf(format, a...))
}

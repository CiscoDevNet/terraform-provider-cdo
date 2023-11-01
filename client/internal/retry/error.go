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
	return "error"
}

// Error facilitates `errors.Is(err, retry.Error)`.
var Error = errorType{}

// TimeoutError facilitates `errors.Is(err, retry.TimeoutError)`.
var TimeoutError = fmt.Errorf("%w: timeout", Error)

// newTimeoutErrorf is used within this package to create new TimeoutError.
func newTimeoutErrorf(title, format string, a ...any) ErrorType {
	return fmt.Errorf("retrying: %s\n%w: %w", title, TimeoutError, fmt.Errorf(format, a...))
}

// RetriesExceededError facilitates `errors.Is(err, retry.RetriesExceededError)`.
var RetriesExceededError = fmt.Errorf("%w: retries exceeded", Error)

// newRetriesExceededErrorf is used within this package to create new RetriesExceededError.
func newRetriesExceededErrorf(title, format string, a ...any) ErrorType {
	return fmt.Errorf("retrying: %s\n%w: %w", title, RetriesExceededError, fmt.Errorf(format, a...))
}

// ContextCancelledError facilitates `errors.Is(err, retry.ContextCancelledError)`.
var ContextCancelledError = fmt.Errorf("%w: context cancelled", Error)

// newContextCancelledErrorf is used within this package to create new ContextCancelledError.
func newContextCancelledErrorf(title, format string, a ...any) ErrorType {
	return fmt.Errorf("retrying: %s\n%w: %w", title, ContextCancelledError, fmt.Errorf(format, a...))
}

// FuncError facilitates `errors.Is(err, retry.FuncError)`.
var FuncError = fmt.Errorf("%w", Error)

// newFuncErrorf is used within this package to create new FuncError.
func newFuncErrorf(title, format string, a ...any) ErrorType {
	return fmt.Errorf("retrying: %s\n%w: %w", title, FuncError, fmt.Errorf(format, a...))
}

// AttemptError facilitates `errors.Is(err, retry.AttemptError)`.
var AttemptError = fmt.Errorf("%w", Error)

// newAttemptErrorf is used within this package to create new FuncError.
func newAttemptErrorf(format string, a ...any) ErrorType {
	return fmt.Errorf("%w: %w", AttemptError, fmt.Errorf(format, a...))
}

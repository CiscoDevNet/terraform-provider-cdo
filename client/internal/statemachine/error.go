package statemachine

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/cdo"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine"
)

type ErrorType interface {
	Error() string
}

type errorType struct{}

func (err errorType) Error() string {
	return "workflow error"
}

// Error facilitates `errors.Is(err, statemachine.Error)`.
var Error ErrorType = errorType{}

var NotFoundError = fmt.Errorf("%w: state machine not found", Error)

var MoreThanOneRunningError = fmt.Errorf("%w: multiple running instances found, this is not expected, please report this issue at: %s", Error, cdo.TerraformProviderCDOIssuesUrl)

var UnknownError = fmt.Errorf("%w: unknown error", Error)

func NewWorkflowErrorf(format string, a ...any) ErrorType {
	return NewWorkflowError(fmt.Errorf(format, a...))
}

func NewWorkflowError(err error) ErrorType {
	return fmt.Errorf("%w: %w", Error, err)
}

func NewWorkflowErrorFromDetails(details statemachine.Details) ErrorType {
	if details.LastError == nil {
		return UnknownError
	} else {
		return NewWorkflowErrorf("message=%s, action=%s, identifier=%s", details.LastError.ErrorMessage, details.LastError.ActionIdentifier, details.LastError.StateMachineIdentifier)
	}
}

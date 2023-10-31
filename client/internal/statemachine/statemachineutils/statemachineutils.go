package statemachineutils

import (
	internalStateMachine "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/statemachine"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine"
)

func ExtractError(details statemachine.Details) internalStateMachine.ErrorType {
	if details.LastError == nil {
		return nil
	} else {
		return internalStateMachine.NewWorkflowErrorf(details.LastError.ErrorMessage)
	}
}

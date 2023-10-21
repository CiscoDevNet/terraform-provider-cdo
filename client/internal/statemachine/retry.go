package statemachine

import (
	"context"
	"errors"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
)

// UntilStarted keeps polling until it finds the state machine with given identifier or error that is not a not found error
func UntilStarted(ctx context.Context, client http.Client, deviceUid string, stateMachineIdentifier string) retry.Func {
	return func() (ok bool, err error) {
		res, err := ReadInstanceByDeviceUid(ctx, client, NewReadInstanceByDeviceUidInput(deviceUid))
		if err != nil {
			if errors.Is(err, NotFoundError) {
				// state machine not found, probably because we are calling too early, and it has not started yet, continue polling
				return false, nil
			}
			// other error, not valid
			return false, err
		}
		if res.StateMachineIdentifier == stateMachineIdentifier {
			// found it, done!
			return true, nil
		}
		// other state machine is running, continue polling
		return false, nil
	}
}

// UntilDoneByIdentifier polls a state machine by its identifier until it has state done, or ends early due to error.
// Note if you know a state machine with the same identifier has been run before, it may be checking the old state machine.
func UntilDoneByIdentifier(ctx context.Context, client http.Client, stateMachineName string) retry.Func {
	return func() (bool, error) {
		sm, err := ReadInstanceByName(ctx, client, NewReadInstanceByNameInput(stateMachineName))
		if err != nil {
			return false, err
		}
		client.Logger.Println(fmt.Sprintf("state machine state=%s", sm.StateMachineInstanceCondition))
		if sm.StateMachineInstanceCondition == state.ERROR {
			if sm.StateMachineDetails.LastError != nil {
				return false, fmt.Errorf(fmt.Sprintf("state machine errored: %s", sm.StateMachineDetails.LastError.ErrorMessage))
			} else {
				return false, fmt.Errorf("state machine errored (no error message)")
			}
		}
		return sm.StateMachineInstanceCondition == state.DONE, nil
	}
}

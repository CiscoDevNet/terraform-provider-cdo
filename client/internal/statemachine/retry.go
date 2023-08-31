package statemachine

import (
	"context"
	"errors"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
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

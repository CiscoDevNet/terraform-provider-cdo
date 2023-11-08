package iosconfig

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/statemachine"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
)

func UntilState(ctx context.Context, client http.Client, specificUid string, expectedState state.Type) retry.Func {

	// create ios config read request
	readReq := NewReadRequest(ctx, client, *NewReadInput(
		specificUid,
	))

	var readOutp ReadOutput

	return func() (bool, error) {
		err := readReq.Send(&readOutp)
		if err != nil {
			return false, err
		}

		client.Logger.Printf("ios config expectedState=%s\n", readOutp.State)
		if readOutp.State == expectedState {
			return true, nil
		}
		if readOutp.State == state.ERROR {
			return false, statemachine.NewWorkflowErrorFromDetails(readOutp.StateMachineDetails)
		}
		if readOutp.State == state.BAD_CREDENTIALS {
			return false, statemachine.NewWorkflowErrorf("Bad Credentials")
		}
		return false, nil
	}
}

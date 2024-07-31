package cloudfmc

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/statemachine"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
)

func UntilStateDone(ctx context.Context, client http.Client, uid string) retry.Func {

	client.Logger.Println("waiting for cdfmc to be done")

	return func() (bool, error) {
		fmcReadSpecificRes, err := ReadSpecific(ctx, client, NewReadSpecificInput(uid))
		if err != nil {
			return false, err
		}
		client.Logger.Printf("fmcReadSpecificRes.State=%s\n", fmcReadSpecificRes.State)
		if fmcReadSpecificRes.State == state.DONE {
			return true, nil
		} else if fmcReadSpecificRes.State == state.ERROR {
			return false, statemachine.NewWorkflowErrorFromDetails(fmcReadSpecificRes.StateMachineDetails)
		}
		return false, nil
	}
}

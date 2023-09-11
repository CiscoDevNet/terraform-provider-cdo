package cloudftd

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
)

func UntilGeneratedCommandAvailable(ctx context.Context, client http.Client, uid string, metadata *Metadata) retry.Func {

	return func() (bool, error) {
		client.Logger.Println("checking if FTD generated command is available")

		readOutp, err := ReadByUid(ctx, client, NewReadByUidInput(uid))
		if err != nil {
			return false, err
		}

		client.Logger.Printf("device metadata=%v\n", readOutp.Metadata)

		if readOutp.Metadata.GeneratedCommand != "" {
			*metadata = readOutp.Metadata
			return true, nil
		} else {
			return false, fmt.Errorf("generated command not found in metadata: %+v", readOutp.Metadata)
		}
	}
}

func UntilStateDone(ctx context.Context, client http.Client, inp ReadByUidInput) retry.Func {
	return func() (bool, error) {
		client.Logger.Println("check FTD state")

		readOutp, err := ReadByUid(ctx, client, inp)
		if err != nil {
			return false, err
		}

		if readOutp.State == state.DONE {
			return true, nil
		} else if readOutp.State == state.ERROR {
			return false, fmt.Errorf("workflow ended in error")
		} else {
			return false, fmt.Errorf("generated command not found in metadata: %+v", readOutp.Metadata)
		}
	}
}

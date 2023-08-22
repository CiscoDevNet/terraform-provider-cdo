package iosconfig

import (
	"context"
	"fmt"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
)

func UntilState(ctx context.Context, client http.Client, specificUid string, expectedState string) retry.Func {

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
		if strings.EqualFold(readOutp.State, expectedState) {
			return true, nil
		}
		if strings.EqualFold(readOutp.State, state.ERROR) {
			return false, fmt.Errorf("workflow ended in %s", state.ERROR)
		}
		if strings.EqualFold(readOutp.State, state.BAD_CREDENTIALS) {
			return false, fmt.Errorf("bad credentials")
		}
		return false, nil
	}
}

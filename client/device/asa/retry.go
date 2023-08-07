package asa

import (
	"context"
	"fmt"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
)

func UntilStateDoneAndConnectivityOk(ctx context.Context, client http.Client, uid string) retry.Func {

	return func() (bool, error) {
		readOutp, err := Read(ctx, client, *NewReadInput(uid))
		if err != nil {
			return false, err
		}

		client.Logger.Printf("device state=%s\n", readOutp.State)

		if strings.EqualFold(readOutp.State, "DONE") && strings.EqualFold(readOutp.Status, "IDLE") {

			if readOutp.ConnectivityState <= 0 {
				return false, fmt.Errorf("connectivity error: %s", readOutp.ConnectivityError)
			}

			return true, nil
		}
		if strings.EqualFold(readOutp.State, "ERROR") {
			return false, fmt.Errorf("workflow ended in ERROR")
		}
		if strings.EqualFold(readOutp.State, "BAD_CREDENTIALS") {
			return false, fmt.Errorf("bad credentials")
		}
		return false, nil
	}
}

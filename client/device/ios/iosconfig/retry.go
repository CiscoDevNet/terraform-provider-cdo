package iosconfig

import (
	"context"
	"fmt"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/internal/retry"
)

const (
	IosConfigStateDone           = "DONE"
	IosConfigStateError          = "ERROR"
	IosConfigStateBadCredentials = "BAD_CREDENTIALS"
)

func UntilState(ctx context.Context, client http.Client, specificUid string, state string) retry.Func {

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

		client.Logger.Printf("ios config state=%s\n", readOutp.State)
		if strings.EqualFold(readOutp.State, state) {
			return true, nil
		}
		if strings.EqualFold(readOutp.State, IosConfigStateError) {
			return false, fmt.Errorf("workflow ended in %s", IosConfigStateError)
		}
		if strings.EqualFold(readOutp.State, IosConfigStateBadCredentials) {
			return false, fmt.Errorf("bad credentials")
		}
		return false, nil
	}
}

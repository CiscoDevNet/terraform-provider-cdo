package asaconfig

import (
	"context"
	"fmt"
	"strings"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/retry"
)

// TODO: Create AsaConfigState type
const (
	AsaConfigStateDone           = "DONE"
	AsaConfigStateError          = "ERROR"
	AsaConfigStateBadCredentials = "BAD_CREDENTIALS"
)

func UntilStateDone(ctx context.Context, client http.Client, specificUid string) retry.Func {

	// create asa config read request
	readReq := NewReadRequest(ctx, client, *NewReadInput(
		specificUid,
	))

	var readOutp ReadOutput

	return func() (bool, error) {
		err := readReq.Send(&readOutp)
		if err != nil {
			return false, err
		}

		client.Logger.Printf("asa config state=%s\n", readOutp.State)
		if strings.EqualFold(readOutp.State, AsaConfigStateDone) {
			return true, nil
		}
		if strings.EqualFold(readOutp.State, AsaConfigStateError) {
			return false, fmt.Errorf("workflow ended in ERROR")
		}
		if strings.EqualFold(readOutp.State, AsaConfigStateBadCredentials) {
			return false, fmt.Errorf("bad credentials")
		}
		return false, nil
	}
}

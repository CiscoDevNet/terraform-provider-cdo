package asaconfig

import (
	"context"
	"fmt"
	"strings"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/retry"
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
		if strings.EqualFold(readOutp.State, "DONE") {
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

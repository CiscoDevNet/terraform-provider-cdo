package iosconfig

import (
	"context"
	"fmt"
	"strings"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/retry"
)

func UntilState(ctx context.Context, client http.Client, specificUid string, state string) retry.Func {

	// create ios config read request
	readReq := NewReadRequest(ctx, client, *NewReadInput(
		specificUid,
	))

	var readOutp ReadOutput

	return func() (bool, error) {
		readReq.Send(&readOutp)

		client.Logger.Printf("ios config state=%s\n", readOutp.State)
		if strings.EqualFold(readOutp.State, state) {
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

package ftdc

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
)

func UntilGeneratedCommandAvailable(ctx context.Context, client http.Client, uid string) retry.Func {

	return func() (bool, error) {
		readOutp, err := ReadByUid(ctx, client, NewReadByUidInput(uid))
		if err != nil {
			return false, err
		}

		client.Logger.Printf("device metadata=%v\n", readOutp.Metadata)

		if readOutp.Metadata.GeneratedCommand != "" {
			return true, nil
		} else {
			return false, fmt.Errorf("generated command not found in metadata: %+v", readOutp.Metadata)
		}
	}
}

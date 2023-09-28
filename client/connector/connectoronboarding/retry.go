package connectoronboarding

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/status"
)

func UntilLarStatusIsActive(ctx context.Context, client http.Client, readInp connector.ReadByNameInput) retry.Func {
	return func() (bool, error) {
		readRes, err := connector.ReadByName(ctx, client, readInp)
		if err != nil {
			return false, err
		}
		if readRes.LarStatus == status.Active {
			return true, nil
		}
		return false, nil
	}
}

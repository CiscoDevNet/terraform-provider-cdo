package connectoronboarding

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"time"
)

type CreateInput struct {
	Name string
}

func NewCreateInput(name string) CreateInput {
	return CreateInput{
		Name: name,
	}
}

type CreateOutput struct {
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	// wait for larStatus to be "Active"
	err := retry.Do(
		UntilLarStatusIsActive(ctx, client, *connector.NewReadByNameInput(createInp.Name)),
		retry.NewOptionsBuilder().
			Timeout(15*time.Minute).
			Retries(-1).
			Delay(2*time.Second).
			Logger(client.Logger).
			EarlyExitOnError(false).
			Build(),
	)

	if err != nil {
		return nil, err
	} else {
		return &CreateOutput{}, nil
	}
}

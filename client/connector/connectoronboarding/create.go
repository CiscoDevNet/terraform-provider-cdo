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

type CreateOutput = connector.ReadOutput

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	// wait for connector status to be "Active"
	var readOutp connector.ReadOutput
	err := retry.Do(
		ctx,
		UntilConnectorStatusIsActive(ctx, client, *connector.NewReadByNameInput(createInp.Name), &readOutp),
		retry.NewOptionsBuilder().
			Message("until connector active").
			Timeout(15*time.Minute). // usually takes ~3 minutes
			Retries(-1).
			Delay(2*time.Second).
			Logger(client.Logger).
			EarlyExitOnError(false).
			Build(),
	)

	if err != nil {
		return nil, err
	} else {
		return &readOutp, nil
	}
}

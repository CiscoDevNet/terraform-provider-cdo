package seconboarding

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sec"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"time"
)

type CreateInput struct {
	Name string
}

type CreateOutput = sec.ReadOutput

func Create(ctx context.Context, client http.Client, createInput CreateInput) (*CreateOutput, error) {

	client.Logger.Printf("creating SEC Onboarding")

	var createOutp CreateOutput
	err := retry.Do(
		ctx,
		func() (bool, error) {
			readOutp, err := sec.ReadByName(ctx, client, sec.NewReadByNameInputBuilder().Name(createInput.Name).Build())
			if err != nil {
				return false, err
			}
			createOutp = *readOutp
			client.Logger.Printf("SEC status=%s\n", readOutp.EsStatus)
			if readOutp.EsStatus == "ACTIVE" {
				return true, nil
			}
			return false, nil
		},
		retry.NewOptionsBuilder().
			Message("Waiting for SEC to be ACTIVE...").
			Retries(-1).
			Delay(3*time.Second).
			Timeout(15*time.Minute). // usually takes 5-15 minutes
			Logger(client.Logger).
			EarlyExitOnError(true).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	return &createOutp, nil
}

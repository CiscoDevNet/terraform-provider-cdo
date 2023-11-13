package sec

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/statemachine"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"time"
)

type CreateInput struct{}

type CreateOutput struct {
	Uid              string
	Name             string
	SecBootstrapData string
	CdoBoostrapData  string
}

type createSecBody struct {
	QueueTriggerState string `json:"queueTriggerState"`
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	// 1. create new sec
	createUrl := url.CreateSec(client.BaseUrl())
	createBody := createSecBody{QueueTriggerState: "ONBOARD_EVENT_STREAMER"}
	req := client.NewPost(ctx, createUrl, createBody)
	var createReqOutput CreateOutput
	if err := req.Send(&createReqOutput); err != nil {
		return nil, err
	}

	// 2. wait for state machine finish
	err := retry.Do(
		ctx,
		statemachine.UntilDone(ctx, client, createReqOutput.Uid, "eventingPushRequest"),
		retry.NewOptionsBuilder().
			Message("Waiting for SEC to be created...").
			Logger(client.Logger).
			EarlyExitOnError(true).
			Timeout(5*time.Minute).
			Retries(-1).
			Delay(time.Second).
			Build(),
	)
	if err != nil {
		return nil, err
	}
	// 3. get sec bootstrap data
	readOutput, err := Read(ctx, client, NewReadInputBuilder().Uid(createReqOutput.Uid).Build())
	if err != nil {
		return nil, err
	}
	secBootstrapData := readOutput.BootStrapData

	// 4. generate cdo bootstrap data
	cdoBootstrapData, err := generateBootstrapData(ctx, client, readOutput.Name)
	if err != nil {
		return nil, err
	}

	// 5. re-read the sec until its name is updated, no idea why the name is empty some time...
	var readOut ReadOutput
	err = retry.Do(
		ctx,
		func() (bool, error) {
			out, err := Read(ctx, client, NewReadInputBuilder().Uid(createReqOutput.Uid).Build())
			if err != nil {
				return false, err
			}
			readOut = *out
			return out.Name != "", nil
		},
		retry.NewOptionsBuilder().
			Message("Waiting for SEC to finalize...").
			Logger(client.Logger).
			EarlyExitOnError(true).
			Timeout(1*time.Minute).
			Retries(-1).
			Delay(500*time.Millisecond).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	// done, create output
	createOutput := CreateOutput{
		Uid:              createReqOutput.Uid,
		Name:             readOut.Name,
		SecBootstrapData: secBootstrapData,
		CdoBoostrapData:  cdoBootstrapData,
	}

	return &createOutput, nil
}

package connector

import (
	"context"
	"errors"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"time"
)

type DeleteInput struct {
	Uid string `json:"-"`
}

func NewDeleteInput(uid string) DeleteInput {
	return DeleteInput{
		Uid: uid,
	}
}

type DeleteOutput struct {
}

func Delete(ctx context.Context, client http.Client, inp DeleteInput) (*DeleteOutput, error) {

	var deleteOutp DeleteOutput

	// retry delete until connector is not present
	// we cant just delete normally because sometimes it can fail to delete for some unknown reason in CI
	err := retry.Do(
		ctx,
		func() (bool, error) {
			if connectorPresent(ctx, client, inp.Uid) {
				err := sendDeleteConnectorRequest(ctx, client, inp.Uid, &deleteOutp)
				if errors.Is(err, http.NotFoundError) {
					return false, nil // delete in progress, ignore not found error
				}
				return false, err // return other error, if any
			} else {
				return true, nil // connector no longer present, delete works just fine
			}
		},
		retry.NewOptionsBuilder().
			Retries(-1).
			Message("Waiting for connector to be deleted...").
			Logger(client.Logger).
			Timeout(3*time.Minute).
			EarlyExitOnError(true).
			Delay(500*time.Millisecond).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	return &deleteOutp, nil
}

func connectorPresent(ctx context.Context, client http.Client, uid string) bool {
	connectors, err := ReadAll(ctx, client, *NewReadAllInput())
	if err != nil {
		return false
	}
	present := false
	for _, outp := range *connectors {
		if outp.Uid == uid {
			present = true
		}
	}
	return present
}

func sendDeleteConnectorRequest(ctx context.Context, client http.Client, uid string, deleteOutp *DeleteOutput) error {
	deleteUrl := url.DeleteConnector(client.BaseUrl(), uid)
	deleteReq := client.NewDelete(ctx, deleteUrl)
	if err := deleteReq.Send(&deleteOutp); err != nil {
		return err
	}
	return nil
}

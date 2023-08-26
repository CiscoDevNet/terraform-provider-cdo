package statemachine

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine"
)

type ReadInstanceByDeviceUidInput struct {
	Uid string // Uid of the device that runs the state machine
}

func NewReadInstanceByDeviceUidInput(deviceUid string) ReadInstanceByDeviceUidInput {
	return ReadInstanceByDeviceUidInput{
		Uid: deviceUid,
	}
}

type ReadInstanceByDeviceUidOutput = statemachine.Instance

func ReadInstanceByDeviceUid(ctx context.Context, client http.Client, readInp ReadInstanceByDeviceUidInput) (*ReadInstanceByDeviceUidOutput, error) {

	readUrl := url.ReadStateMachineInstance(client.BaseUrl())
	req := client.NewGet(ctx, readUrl)
	req.QueryParams.Add("limit", "1")
	req.QueryParams.Add("q", fmt.Sprintf("objectReference.uid:%s", readInp.Uid))
	req.QueryParams.Add("sort", "lastActiveDate:desc")

	var readRes []ReadInstanceByDeviceUidOutput
	if err := req.Send(&readRes); err != nil {
		return nil, err
	}
	if len(readRes) == 0 {
		return nil, StateMachineNotFoundError
	}

	// TODO: this can happen, no idea why, limit 1 does not seems to work
	//if len(readRes) > 1 {
	//	return nil, MoreThanOneStateMachineRunningError
	//}

	return &readRes[0], nil
}

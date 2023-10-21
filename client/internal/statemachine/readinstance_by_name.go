package statemachine

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine"
)

type ReadInstanceByNameInput struct {
	Name string
}

func NewReadInstanceByNameInput(name string) ReadInstanceByNameInput {
	return ReadInstanceByNameInput{
		Name: name,
	}
}

type ReadInstanceByNameOutput = statemachine.Instance

var NewReadInstanceByNameOutputBuilder = statemachine.NewInstanceBuilder

func ReadInstanceByName(ctx context.Context, client http.Client, readInp ReadInstanceByNameInput) (*ReadInstanceByNameOutput, error) {

	readUrl := url.ReadStateMachineInstance(client.BaseUrl())
	req := client.NewGet(ctx, readUrl)
	req.QueryParams.Add("limit", "1")
	req.QueryParams.Add("q", fmt.Sprintf("stateMachineIdentifier:%s", readInp.Name))
	req.QueryParams.Add("sort", "lastActiveDate:desc")

	var readRes []ReadInstanceByDeviceUidOutput
	if err := req.Send(&readRes); err != nil {
		return nil, err
	}
	if len(readRes) == 0 {
		return nil, NotFoundError
	}

	if len(readRes) > 1 {
		return nil, MoreThanOneRunningError
	}

	return &readRes[0], nil
}

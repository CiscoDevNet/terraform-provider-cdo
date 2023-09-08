package fmcconfig

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmcconfig"
)

type ReadTaskStatusInput struct {
	FmcDomainUid string
	TaskId       string
}

func NewReadTaskStatusInput(fmcDomainUid string, taskId string) ReadTaskStatusInput {
	return ReadTaskStatusInput{
		FmcDomainUid: fmcDomainUid,
		TaskId:       taskId,
	}
}

type ReadTaskStatusOutput = fmcconfig.TaskStatus

func ReadTaskStatus(ctx context.Context, client http.Client, readTaskStatusInp ReadTaskStatusInput) (*ReadTaskStatusOutput, error) {

	readUrl := url.ReadFmcTaskStatus(client.BaseUrl(), readTaskStatusInp.FmcDomainUid, readTaskStatusInp.TaskId)
	req := client.NewGet(ctx, readUrl)

	var readOutp ReadTaskStatusOutput
	if err := req.Send(&readOutp); err != nil {
		return nil, err
	}

	return &readOutp, nil
}

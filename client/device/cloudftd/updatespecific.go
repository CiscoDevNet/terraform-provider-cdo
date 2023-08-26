package cloudftd

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type UpdateSpecificFtdInput struct {
	SpecificUid       string
	QueueTriggerState string
}

func NewUpdateSpecificFtdInput(specificUid string, queueTriggerState string) UpdateSpecificFtdInput {
	return UpdateSpecificFtdInput{
		SpecificUid:       specificUid,
		QueueTriggerState: queueTriggerState,
	}
}

type UpdateSpecificFtdOutput struct {
	SpecificUid string `json:"uid"`
}

type updateSpecificRequestBody struct {
	QueueTriggerState string `json:"queueTriggerState"`
}

func UpdateSpecific(ctx context.Context, client http.Client, updateInp UpdateSpecificFtdInput) (*UpdateSpecificFtdOutput, error) {

	updateUrl := url.UpdateSpecificCloudFtd(client.BaseUrl(), updateInp.SpecificUid)

	updateBody := updateSpecificRequestBody{
		QueueTriggerState: updateInp.QueueTriggerState,
	}
	req := client.NewPut(ctx, updateUrl, updateBody)
	var updateOutp UpdateSpecificFtdOutput
	if err := req.Send(&updateOutp); err != nil {
		return nil, err
	}

	return &updateOutp, nil
}

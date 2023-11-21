package fmcappliance

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type UpdateInput struct {
	FmcApplianceUid     string
	QueueTriggerState   string
	StateMachineContext map[string]string
}

func NewUpdateInput(FmcApplianceUid, queueTriggerState string, stateMachineContext map[string]string) UpdateInput {
	return UpdateInput{
		FmcApplianceUid:     FmcApplianceUid,
		QueueTriggerState:   queueTriggerState,
		StateMachineContext: stateMachineContext,
	}
}

type UpdateOutput struct {
	Uid       string `json:"uid"`
	State     string `json:"state"`
	DomainUid string `json:"domainUid"`
}

type updateRequestBody struct {
	QueueTriggerState   string            `json:"queueTriggerState"`
	StateMachineContext map[string]string `json:"stateMachineContext"`
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {
	updateUrl := url.UpdateFmcAppliance(client.BaseUrl(), updateInp.FmcApplianceUid)
	updateBody := newUpdateRequestBodyBuilder().
		QueueTriggerState(updateInp.QueueTriggerState).
		StateMachineContext(updateInp.StateMachineContext).
		Build()
	req := client.NewPut(ctx, updateUrl, updateBody)
	var updateOup UpdateOutput
	if err := req.Send(&updateOup); err != nil {
		return nil, err
	}

	return &updateOup, nil
}

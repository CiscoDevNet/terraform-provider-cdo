package application

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/application/applicationstatus"
)

type ReadInput struct {
}

type ReadOutput struct {
	Uid                string                 `json:"uid"`
	Name               string                 `json:"name"`
	Version            int                    `json:"version"`
	ApplicationType    string                 `json:"applicationType"`
	ApplicationStatus  applicationstatus.Type `json:"applicationStatus"`
	ApplicationContent ApplicationContent     `json:"applicationContent"`
}

type ApplicationContent struct {
	Type                       string      `json:"@type"`
	FmceDeviceUid              interface{} `json:"fmceDeviceUid"`
	DevicesCount               int         `json:"devicesCount"`
	SfcnDevicesCount           int         `json:"sfcnDevicesCount"`
	FmcApplianceUid            interface{} `json:"fmcApplianceUid"`
	RequestedDevicesCount      int         `json:"requestedDevicesCount"`
	EstimatedDevicesCountRange string      `json:"estimatedDevicesCountRange"`
}

func Read(ctx context.Context, client http.Client, readInp ReadInput) (*ReadOutput, error) {
	// create request
	readUrl := url.ReadApplication(client.BaseUrl())
	readReq := client.NewGet(ctx, readUrl)

	// send request & map response
	var readApplicationOutput []ReadOutput
	err := readReq.Send(&readApplicationOutput)
	if err != nil {
		return nil, err
	}

	// check and return
	if len(readApplicationOutput) < 1 {
		return nil, NotFoundError
	}
	return &readApplicationOutput[0], nil
}

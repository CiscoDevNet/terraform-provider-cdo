package application

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/status"
)

type ReadInput struct {
}

type ReadOutput struct {
	Uid                string             `json:"uid"`
	Name               string             `json:"name"`
	Version            int                `json:"version"`
	ApplicationType    string             `json:"applicationType"`
	ApplicationStatus  string             `json:"applicationStatus"`
	ApplicationContent ApplicationContent `json:"applicationContent"`
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

func Read(ctx context.Context, client http.Client, readInp ReadInput) (ReadOutput, error) {
	readUrl := url.ReadApplication(client.BaseUrl())
	readReq := client.NewGet(ctx, readUrl)
	var readApplicationOutput []ReadOutput
	err = readReq.Send(&readApplicationOutput)
	if err != nil {
		return false, err
	}
	if len(readApplicationOutput) < 1 {
		return false, nil // fmc not yet present, should come up soon
	}
	fmc := readApplicationOutput[0]
	if fmc.ApplicationStatus != status.Active {
		return false, nil
	}
}

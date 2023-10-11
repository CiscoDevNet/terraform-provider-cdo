package cloudfmc

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/goutil"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/status"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"
)

type CreateInput struct {
}

func NewCreateInput() ReadInput {
	return ReadInput{}
}

type createApplicationBody struct {
	ApplicationType    string             `json:"applicationType"`
	ApplicationStatus  string             `json:"applicationStatus"`
	ApplicationContent applicationContent `json:"applicationContent"`
}

type applicationContent struct {
	Type string `json:"@type"`
}

type CreateOutput = device.ReadOutput

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*ReadOutput, error) {

	client.Logger.Println("creating cloud FMC")

	// 1. POST /aegis/rest/v1/services/targets/applications
	createApplicationUrl := url.CreateApplication(client.BaseUrl())
	createApplicationReq := client.NewPost(
		ctx,
		createApplicationUrl,
		createApplicationBody{
			ApplicationType:   "FMCE",
			ApplicationStatus: "REQUESTED",
			ApplicationContent: applicationContent{
				Type: "FmceApplicationContent",
			},
		},
	)
	var createApplicationOutp device.CreateOutput
	err := createApplicationReq.Send(&createApplicationOutp)
	if err != nil {
		return nil, err
	}

	// 2. POST /aegis/rest/v1/services/targets/devices
	createDeviceRes, err := device.Create(ctx, client, device.NewCreateInputBuilder().
		Name("FMC").
		DeviceType(devicetype.CloudFmc).
		Model(false).
		ConnectorType("CDG").
		IgnoreCertificate(goutil.NewBoolPointer(true)).
		EnableOobDetection(goutil.NewBoolPointer(false)).
		Build(),
	)
	if err != nil {
		return nil, err
	}

	// https://ci.dev.lockhart.io/aegis/rest/v1/services/targets/applications
	err := retry.Do(func() (bool, error) {
		readApplicationUrl := url.ReadApplication(client.BaseUrl())
		readApplicationReq := client.NewGet(ctx, readApplicationUrl)
		var readApplicationOutput []device.ReadOutput
		err = readApplicationReq.Send(&readApplicationOutput)
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
	},
		*retry.NewOptionsWithLogger(client.Logger),
	)

	return nil, nil
}

// {"name":"FMC","deviceType":"FMCE","larType":"CDG","ignoreCertificate":true,"model":false,"enableOobDetection":true}

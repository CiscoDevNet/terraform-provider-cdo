package cloudfmc

import (
	"context"
	"errors"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/application"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/goutil"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/application/applicationstatus"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"
	"time"
)

type CreateInput struct {
}

func NewCreateInput() CreateInput {
	return CreateInput{}
}

type createApplicationBody struct {
	ApplicationType    string             `json:"applicationType"`
	ApplicationStatus  string             `json:"applicationStatus"`
	ApplicationContent applicationContent `json:"applicationContent"`
}

type applicationContent struct {
	Type string `json:"@type"`
}

type CreateOutput = device.CreateOutput

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	client.Logger.Println("Creating application object for cdFMC")

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

	client.Logger.Println("creating cloud FMC device")

	// 2. POST /aegis/rest/v1/services/targets/devices
	_, err = device.Create(ctx, client, device.NewCreateInputBuilder().
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

	client.Logger.Println("waiting for fmce state machine to be done")

	err = retry.Do(
		ctx,
		untilApplicationActive(ctx, client),
		retry.NewOptionsBuilder().
			Message("Waiting for cdFMC to be created...").
			Retries(-1).
			Timeout(30*time.Minute). // usually takes about 15-20 minutes
			Delay(3*time.Second).
			EarlyExitOnError(true).
			Logger(client.Logger).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	client.Logger.Println("re-reading cdFMC to get latest info")
	// re-read the cdfmc so that we get updated output, e.g. updated hostname
	readOutp, err := Read(ctx, client, NewReadInput())
	if err != nil {
		return nil, err
	}

	client.Logger.Println("cloud FMC application successfully created")

	return readOutp, nil
}

func untilApplicationActive(ctx context.Context, client http.Client) retry.Func {
	var unreachable bool
	var initialUnreachableTime time.Time
	return func() (bool, error) {
		fmc, err := application.Read(ctx, client, application.ReadInput{})
		if err != nil {
			if !errors.Is(err, application.NotFoundError) {
				// maybe the application is not created yet, and hopefully this is temporarily, ignoring
				return false, nil
			}
			return false, err
		}
		if fmc.ApplicationStatus == applicationstatus.Unreachable {
			// initial unreachable is possibly caused by https://jira-eng-rtp3.cisco.com/jira/browse/LH-71821
			// wait for some time to confirm it is actually unreachable
			if unreachable {
				if initialUnreachableTime.Add(time.Minute * 5).After(time.Now()) {
					// if long enough time has passed, and we are still unreachable, treat it as actual error
					return false, err
				}
			} else {
				unreachable = true
				initialUnreachableTime = time.Now()
				return false, nil
			}
		}
		return fmc.ApplicationStatus == applicationstatus.Active, nil
	}
}

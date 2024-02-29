package cloudftd

import (
	"context"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcappliance"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcplatform"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
)

type UpdateInput struct {
	Uid      string
	Name     string
	Tags     tags.Type
	Licenses []license.Type
}

func NewUpdateInput(uid, name string, tags tags.Type, licenses []license.Type) UpdateInput {
	return UpdateInput{
		Uid:      uid,
		Name:     name,
		Tags:     tags,
		Licenses: licenses,
	}
}

type updateRequestBody struct {
	Name string    `json:"name"`
	Tags tags.Type `json:"tags"`
}

type UpdateOutput = ReadOutput

var NewUpdateOutputBuilder = NewReadOutputBuilder

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating FTD")

	client.Logger.Println("updating CDO settings")

	// update CDO settings
	updateUrl := url.UpdateDevice(client.BaseUrl(), updateInp.Uid)
	updateBody := updateRequestBody{
		Name: updateInp.Name,
		Tags: updateInp.Tags,
	}
	req := client.NewPut(ctx, updateUrl, updateBody)
	var updateOutp UpdateOutput
	if err := req.Send(&updateOutp); err != nil {
		return nil, err
	}

	// read FMC host
	client.Logger.Println("reading FMC host")
	fmcRes, err := cloudfmc.Read(ctx, client, cloudfmc.NewReadInput())
	if err != nil {
		return nil, err
	}

	// update FTD license through FMC api
	client.Logger.Println("updating FTD licenses")
	_, err = fmcplatform.UpdateDeviceLicenses(
		ctx,
		client,
		fmcplatform.NewUpdateDeviceLicensesInputBuilder().
			FmcHost(fmcRes.Host).
			LicenseTypes(updateInp.Licenses).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	// trigger oob detection in cdo to sync license changes

	// read FMC that manages this cloud FTD for its uid, so that we can get its specific uid later
	fmcReadRes, err := cloudfmc.Read(ctx, client, cloudfmc.NewReadInput())
	if err != nil {
		return nil, err
	}

	// read FMC specific device, i.e. the actual FMC for its device uid
	fmcReadSpecificRes, err := cloudfmc.ReadSpecific(ctx, client, cloudfmc.NewReadSpecificInput(fmcReadRes.Uid))
	if err != nil {
		return nil, err
	}

	// trigger oob detection
	_, err = fmcappliance.Update(
		ctx,
		client,
		fmcappliance.NewUpdateInputBuilder().
			FmcApplianceUid(fmcReadSpecificRes.SpecificUid).
			QueueTriggerState("PENDING_OOB_DETECTION").
			Build(),
	)
	if err != nil {
		return nil, err
	}

	// waiting for oob to be done
	err = retry.Do(
		ctx,
		cloudfmc.UntilStateDone(ctx, client, fmcReadRes.Uid),
		retry.NewOptionsBuilder().
			Message("Waiting for FTD to be updated...").
			Retries(-1).
			Timeout(5*time.Minute).
			Logger(client.Logger).
			EarlyExitOnError(true).
			Delay(3*time.Second).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	client.Logger.Println("re-reading FTD for latest info")

	// re-read the FTD for new license change
	ftdReadRes, err := ReadByUid(ctx, client, NewReadByUidInput(updateInp.Uid))
	if err != nil {
		return nil, err
	}

	return ftdReadRes, nil
}

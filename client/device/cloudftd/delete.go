package cloudftd

import (
	"context"
	"errors"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcappliance"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcplatform"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/statemachine"
	"strings"
	"time"
)

type DeleteInput struct {
	Uid string
}

func NewDeleteInput(uid string) DeleteInput {
	return DeleteInput{
		Uid: uid,
	}
}

type DeleteOutput = device.DeleteOutput

func Delete(ctx context.Context, client http.Client, deleteInp DeleteInput) (*DeleteOutput, error) {

	// 1. read FMC that manages this cloud FTD
	fmcReadRes, err := cloudfmc.Read(ctx, client, cloudfmc.NewReadInput())
	if err != nil {
		return nil, err
	}

	// 2. read FMC specific device, i.e. the actual FMC
	fmcReadSpecificRes, err := cloudfmc.ReadSpecific(ctx, client, cloudfmc.NewReadSpecificInput(fmcReadRes.Uid))
	if err != nil {
		return nil, err
	}

	// 2.5 wait for any FTD deployment to finish, otherwise backend fmceDeleteFtdcStateMachine will fail
	// in order to check for FTD deployment status, we need to read FMC host and domainUid
	// 2.5.2 read fmc domain info in its specific/appliance device, fmcDomainUid is in the domain info
	readFmcDomainRes, err := fmcplatform.ReadFmcDomainInfo(ctx, client, fmcplatform.NewReadDomainInfoInput(fmcReadRes.Host))
	if err != nil {
		return nil, err
	}
	if len(readFmcDomainRes.Items) == 0 {
		return nil, fmt.Errorf("%w: failed to read fmc domain uid, fmc domain info not found", http.NotFoundError)
	}
	fmcDomainUid := readFmcDomainRes.Items[0].Uuid

	// now we find the FTD device record, to do that we need to find all device records
	// then we find the FTD device record with the same name as the CDO FTD
	// 2.5.1 read FTD name from input uid
	readFtdOutp, err := ReadByUid(ctx, client, NewReadByUidInput(deleteInp.Uid))
	if err != nil {
		return nil, err
	}
	// 2.5.2 read all device records
	allDeviceRecords, err := fmcconfig.ReadAllDeviceRecords(ctx, client, fmcconfig.NewReadAllDeviceRecordsInput(fmcDomainUid, fmcReadRes.Host))
	if err != nil {
		return nil, err
	}
	var ftdRecordId string
	// 2.5.3 check if FTD name is present in device records, logic: same name + both are FTDs = found
	client.Logger.Printf("looking for existing FTD with id=%s and name=%s\n", deleteInp.Uid, fmcReadRes.Name)
	for _, record := range allDeviceRecords.Items {
		if record.Name != readFtdOutp.Name {
			// different name, ignore
			continue
		}
		// the allDeviceRecords only contains the name, so we need to make another call to retrieve the details of the device to check whether this is a FTD
		// potentially we will be making a lot of network calls and cause this loop to run for long time if
		// we have many device records with the same name, I suppose that rarely happens
		deviceRecord, err := fmcconfig.ReadDeviceRecord(ctx, client, fmcconfig.NewReadDeviceRecordInput(fmcDomainUid, fmcReadRes.Host, record.Id))
		if err != nil {
			return nil, err
		}
		if strings.Contains(deviceRecord.Model, "Firepower Threat Defense") { // Question: is there a better way to check? Does this check cover all cases?
			// found
			ftdRecordId = record.Id
		} // else not a FTD, just some other device with the same name, ignore
	}

	if ftdRecordId != "" {
		// FTD record found in the cdFMC device records
		// now use the FTD device record id and to read the device record, until the FTD deployment status is DEPLOYED

		// 2.5.4 re-read FMC's FTD deployment status until it is DEPLOYED
		err = retry.Do(
			ctx,
			func() (bool, error) {
				ftdDeviceRecord, err := fmcconfig.ReadDeviceRecord(ctx, client, fmcconfig.NewReadDeviceRecordInput(fmcDomainUid, fmcReadRes.Host, ftdRecordId))
				if err != nil {
					return false, err
				}
				client.Logger.Printf("current deployment status=%s\n", ftdDeviceRecord.DeploymentStatus)
				if ftdDeviceRecord.DeploymentStatus == "DEPLOYED" || ftdDeviceRecord.DeploymentStatus == "" { // no idea what this deployment status could be, so it is a string
					return true, nil
					// TODO: check for error here: like if ftdDeviceRecord.DeploymentStatus == "DEPLOY_ERROR" {, not sure what is the error deployment status so I did not do it for now
				} else {
					return false, nil
				}
			},
			retry.NewOptionsBuilder().
				Message("Waiting for FTD deployment to finish...").
				Timeout(15*time.Minute). // usually 5-10 minutes
				Delay(3*time.Second).
				Logger(client.Logger).
				Retries(-1).
				EarlyExitOnError(true).
				Build(),
		)
		if err != nil {
			// error during retry, maybe the deployment failed, we fail the delete as well
			return nil, err
		}

	} // else
	// the ftd device record is not found in cdFMC if reach here, this is unexpected,
	// but since we are performing delete operation, it is going to be deleted anyway, so ignore it

	// now we checked the FTD is ready to be deleted

	// 3. delete FTD in cdFMC, schedule a state machine for cloud fmc to delete the cloud FTD
	_, err = fmcappliance.Update(
		ctx,
		client,
		fmcappliance.NewUpdateInputBuilder().
			FmcApplianceUid(fmcReadSpecificRes.SpecificUid).
			QueueTriggerState("PENDING_DELETE_FTDC").
			StateMachineContext(&map[string]string{"ftdCDeviceIDs": deleteInp.Uid}).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	// 4. wait until the delete cloud FTD state machine is done
	err = retry.Do(
		ctx,
		statemachine.UntilDone(ctx, client, fmcReadSpecificRes.SpecificUid, "fmceDeleteFtdcStateMachine"),
		retry.NewOptionsBuilder().
			Message("Waiting for FTD deletion to finish...").
			Retries(retry.DefaultRetries).
			Delay(retry.DefaultDelay).
			Logger(client.Logger).
			EarlyExitOnError(true).
			Timeout(retry.DefaultTimeout).
			Build(),
	)
	// skip 404 errors, we are doing deletion anyway
	if err != nil {
		if !errors.Is(err, http.NotFoundError) {
			return nil, err
		}
	}

	// 5. delete FTD in CDO as well,
	ftdDeleteOutput, err := device.Delete(ctx, client, *device.NewDeleteInput(deleteInp.Uid))
	// similarly, skip 404 errors, we are doing deletion anyway
	if err != nil {
		if !errors.Is(err, http.NotFoundError) {
			return nil, err
		}
	}

	// done!
	return ftdDeleteOutput, nil

}

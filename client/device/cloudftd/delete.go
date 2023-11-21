package cloudftd

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcappliance"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/statemachine"
)

type DeleteInput struct {
	Uid string
}

func NewDeleteInput(uid string) DeleteInput {
	return DeleteInput{
		Uid: uid,
	}
}

type DeleteOutput struct {
}

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

	// 3. schedule a state machine for cloud fmc to delete the cloud FTD
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

	// 4. wait until the delete cloud FTD state machine has started
	err = retry.Do(
		ctx,
		statemachine.UntilStarted(ctx, client, fmcReadSpecificRes.SpecificUid, "fmceDeleteFtdcStateMachine"),
		retry.NewOptionsBuilder().
			Message("Waiting for FTD deletion to begin...").
			Retries(retry.DefaultRetries).
			Delay(retry.DefaultDelay).
			Logger(client.Logger).
			EarlyExitOnError(true).
			Timeout(retry.DefaultTimeout).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	// 5. we are not waiting for it to finish, like the CDO UI

	// done!
	return &DeleteOutput{}, nil

}

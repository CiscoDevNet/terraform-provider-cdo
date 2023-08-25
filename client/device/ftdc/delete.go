package ftdc

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cdfmc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cdfmc/fmcappliance"
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

	// 1. read FMC that manages this FTDc
	cdfmcReadRes, err := cdfmc.Read(ctx, client, cdfmc.NewReadInput())
	if err != nil {
		return nil, err
	}

	// 2. read FMC specific device, i.e. the actual FMC
	cdfmcReadSpecificRes, err := cdfmc.ReadSpecific(ctx, client, cdfmc.NewReadSpecificInput(cdfmcReadRes.Uid))
	if err != nil {
		return nil, err
	}

	// 3. schedule a state machine for fmc to delete the FTDc
	_, err = fmcappliance.Update(ctx, client, fmcappliance.NewUpdateInput(
		cdfmcReadSpecificRes.SpecificUid,
		"PENDING_DELETE_FTDC",
		map[string]string{
			"ftdCDeviceIDs": deleteInp.Uid,
		},
	))
	if err != nil {
		return nil, err
	}

	// 4. wait until the delete FTDc state machine has started
	err = retry.Do(statemachine.UntilStarted(ctx, client, cdfmcReadSpecificRes.SpecificUid, "fmceDeleteFtdcStateMachine"), *retry.NewOptionsWithLoggerAndRetries(client.Logger, 3))
	if err != nil {
		return nil, err
	}

	// 5. we are not waiting for it to finish, like the CDO UI

	// done!
	return &DeleteOutput{}, nil

}

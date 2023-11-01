package cloudftdonboarding

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcplatform"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"
)

type CreateInput struct {
	FtdUid string
}

func NewCreateInput(ftdId string) CreateInput {
	return CreateInput{
		FtdUid: ftdId,
	}
}

type CreateOutput = fmcconfig.CreateDeviceRecordOutput

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	client.Logger.Println("creating cloud ftd onboarding")

	// 1. read fmc domain uid
	client.Logger.Println("retrieving fmc domain uid")

	// 1.1 read fmc
	fmcRes, err := cloudfmc.Read(ctx, client, cloudfmc.NewReadInput())
	if err != nil {
		return nil, err
	}
	// 1.2 read fmc domain info in its specific/appliance device
	readFmcDomainRes, err := fmcplatform.ReadFmcDomainInfo(ctx, client, fmcplatform.NewReadDomainInfoInput(fmcRes.Host))
	if err != nil {
		return nil, err
	}
	if len(readFmcDomainRes.Items) == 0 {
		return nil, fmt.Errorf("failed to read fmc domain uid, fmc domain info not found")
	}
	fmcDomainUid := readFmcDomainRes.Items[0].Uuid

	// 1.5 check device already registered
	// 1.5.1 read FTD name
	readFtdOutp, err := cloudftd.ReadByUid(ctx, client, cloudftd.NewReadByUidInput(createInp.FtdUid))
	if err != nil {
		return nil, err
	}
	// 1.5.2 read all device records
	allDeviceRecords, err := fmcconfig.ReadAllDeviceRecords(ctx, client, fmcconfig.NewReadAllDeviceRecordsInput(fmcDomainUid, fmcRes.Host))
	if err != nil {
		return nil, err
	}
	// 1.5.3 check if FTD name is present in device records, logic: same name + both are FTDs = duplicate
	client.Logger.Printf("checking if FTD already exists with id=%s and name=%s\n", createInp.FtdUid, fmcRes.Name)
	for _, record := range allDeviceRecords.Items {
		if record.Name != readFtdOutp.Name {
			// different name, ignore
			continue
		}
		// the allDeviceRecords only contains the name, so we need to make another call to retrieve the details of the device to check whether this is a FTD
		// potentially we will be making a lot of network calls and cause this loop to run for long time if
		// we have many device records with the same name, I suppose that rarely happens
		deviceRecord, err := fmcconfig.ReadDeviceRecord(ctx, client, fmcconfig.NewReadDeviceRecordInput(fmcDomainUid, fmcRes.Host, record.Id))
		if err != nil {
			return nil, err
		}
		if strings.Contains(deviceRecord.Model, "Firepower Threat Defense") { // Question: is there a better way to check? Does this check cover all cases?
			return nil, fmt.Errorf("FTD with id=%s and name=%s is already registered", createInp.FtdUid, fmcRes.Name)
		} else {
			// not a FTD, just some other device with the same name, ignore
		}
	}
	client.Logger.Printf("FTD with id=%s and name=%s is not registered, proceeding\n", createInp.FtdUid, fmcRes.Name)

	// 2. get a system token for creating FTD device record in FMC
	// CDO token does not work, we will get a 405 method not allowed if we do that
	client.Logger.Println("getting a system token for creating FTD device record in FMC")

	// 2.1 get tenant uid from API token
	readOutp, err := user.GetTokenInfo(ctx, client, user.NewGetTokenInfoInput())
	if err != nil {
		return nil, err
	}
	tenantUid := readOutp.UserAuthentication.Details.TenantUid
	// 2.2 get a system token with scope of the tenant uid of the API token
	createTokenOutp, err := user.CreateSystemToken(ctx, client, user.NewCreateSystemTokenInput(tenantUid))
	if err != nil {
		return nil, err
	}

	// 3. post device record
	client.Logger.Println("creating FTD device record in FMC")

	// 3.1 read ftd metadata
	// 3.1.5 handle license
	licenseCaps, err := license.DeserializeAllFromCdo(readFtdOutp.Metadata.LicenseCaps)
	if err != nil {
		return nil, err
	}
	// 3.2 create ftd device
	createDeviceInp := fmcconfig.NewCreateDeviceRecordInputBuilder().
		Type("Device").
		NatId(readFtdOutp.Metadata.NatID).
		Name(readFtdOutp.Name).
		AccessPolicyUid(readFtdOutp.Metadata.AccessPolicyUid).
		LicenseCaps(&licenseCaps).
		PerformanceTier(readFtdOutp.Metadata.PerformanceTier).
		RegKey(readFtdOutp.Metadata.RegKey).
		FmcDomainUid(fmcDomainUid).
		FmcHostname(fmcRes.Host).
		SystemApiToken(createTokenOutp.AccessToken).
		Build()
	var createOutp fmcconfig.CreateDeviceRecordOutput
	err = retry.Do(
		ctx,
		fmcconfig.UntilCreateDeviceRecordSuccess(ctx, client, createDeviceInp, &createOutp),
		retry.NewOptionsBuilder().
			Message("Waiting for FTD device record to be created on cdFMC...").
			Retries(-1).
			Delay(3*time.Second).
			Timeout(1*time.Hour). // it can take 15-20 minutes for FTD to come up + 10 minutes to create device record
			Logger(client.Logger).
			EarlyExitOnError(false).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	// 4. trigger FTD onboarding state machine
	client.Logger.Println("triggering ftdcRegisterStateMachine")

	// 4.1 get ftd specific device
	ftdSpecificOutp, err := cloudftd.ReadSpecific(ctx, client, cloudftd.NewReadSpecificInputBuilder().Uid(readFtdOutp.Uid).Build())
	if err != nil {
		return nil, err
	}
	// 4.2 trigger register state machine
	_, err = cloudftd.UpdateSpecific(ctx, client,
		cloudftd.NewUpdateSpecificFtdInput(
			ftdSpecificOutp.SpecificUid,
			"INITIATE_FTDC_REGISTER",
		),
	)
	if err != nil {
		return nil, err
	}
	// 4.3 wait until state machine done
	err = retry.Do(
		ctx,
		cloudftd.UntilSpecificStateDone(
			ctx,
			client,
			cloudftd.NewReadSpecificInputBuilder().
				Uid(readFtdOutp.Uid).
				Build(),
		),
		retry.NewOptionsBuilder().
			Message("Waiting for FTD to be onboarded to CDO...").
			Retries(-1).
			Delay(1*time.Second).
			Timeout(20*time.Minute). // usually done in less than 5 minutes because we already registered in FTDc
			Logger(client.Logger).
			EarlyExitOnError(false).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	return &createOutp, nil
}

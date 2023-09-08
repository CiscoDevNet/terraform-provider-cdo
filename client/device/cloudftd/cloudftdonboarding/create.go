package cloudftdonboarding

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcplatform"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"
	"time"
)

type CreateInput struct {
	FtdName string
}

func NewCreateInput(ftdName string) CreateInput {
	return CreateInput{
		FtdName: ftdName,
	}
}

type CreateOutput struct {
}

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
	readFmcDomainRes, err := fmcplatform.ReadFmcDomainInfo(ctx, client, fmcplatform.NewReadDomainInfo(fmcRes.Host))
	if err != nil {
		return nil, err
	}
	if len(readFmcDomainRes.Items) == 0 {
		return nil, fmt.Errorf("failed to read fmc domain uid, fmc domain info not found")
	}
	fmcDomainInfo := readFmcDomainRes.Items[0] // contains uid

	// 2. get a system token for creating FTD device record in FMC
	// CDO token does not work, otherwise 405 method not allowed
	client.Logger.Println("getting a system token for creating FTD device record in FMC")

	// 2.1 get tenant context => tenant uid
	readOutp, err := user.GetTokenInfo(ctx, client, user.NewGetTokenInfoInput())
	if err != nil {
		return nil, err
	}
	tenantUid := readOutp.UserAuthentication.Details.TenantUid
	// 2.2 get a system token with scope of tenant uid
	createTokenOutp, err := user.CreateSystemToken(ctx, client, user.NewCreateSystemTokenInput(tenantUid))
	if err != nil {
		return nil, err
	}

	// 3. post device record
	client.Logger.Println("creating FTD device record in FMC")

	// 3.1 read ftd metadata
	readFtdOutp, err := cloudftd.ReadByName(ctx, client, cloudftd.NewReadByNameInput(createInp.FtdName))
	if err != nil {
		return nil, err
	}
	// 3.2 create ftd device
	createDeviceInp := fmcconfig.NewCreateDeviceRecordInputBuilder().
		Type("Device").
		NatId(readFtdOutp.Metadata.NatID).
		Name(readFtdOutp.Name).
		AccessPolicyUid(readFtdOutp.Metadata.AccessPolicyUid).
		LicenseCaps(readFtdOutp.Metadata.LicenseCaps).
		PerformanceTier(readFtdOutp.Metadata.PerformanceTier).
		RegKey(readFtdOutp.Metadata.RegKey).
		FmcDomainUid(fmcDomainInfo.Uuid).
		SystemApiToken(createTokenOutp.AccessToken).
		Build()
	err = retry.Do(
		fmcconfig.UntilCreateDeviceRecordSuccess(ctx, client, createDeviceInp),
		retry.NewOptionsBuilder().
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
	client.Logger.Println("re-triggering FTD onboarding state machine")

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
	// TODO: wait until state done

	return nil, nil
}

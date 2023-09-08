package cloudftdonboarding

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcplatform"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"
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

	// TODO: loop here
	// 3. post device record
	client.Logger.Println("creating FTD device record in FMC")
	// 3.1 read ftd metadata
	readFtdOutp, err := cloudftd.ReadByName(ctx, client, cloudftd.NewReadByNameInput(createInp.FtdName))
	if err != nil {
		return nil, err
	}
	// 3.2 create ftd device
	// TODO: use system token
	createDeviceInp := fmcconfig.NewCreateDeviceRecordInputBuilder().
		Type("Device").
		NatId(readFtdOutp.Metadata.NatID).
		Name(readFtdOutp.Name).
		AccessPolicyUid(readFtdOutp.Metadata.AccessPolicyUid).
		LicenseCaps(readFtdOutp.Metadata.LicenseCaps).
		PerformanceTier(readFtdOutp.Metadata.PerformanceTier).
		RegKey(readFtdOutp.Metadata.RegKey).
		FmcDomainUid(fmcDomainInfo.Uuid).
		Build()
	createOutp, err := fmcconfig.CreateDeviceRecord(ctx, client, createDeviceInp)
	if err != nil {
		return nil, err
	}

	// 4 read task until success
	// TODO: loop here
	readTaskOutp, err := fmcconfig.ReadTaskStatus(ctx, client, fmcconfig.NewReadTaskStatusInput(fmcDomainInfo.Uuid, createOutp.Id))
	if err != nil {
		return nil, err
	}
	if readTaskOutp.Status == "RUNNING" {
		// TODO: continue
	}

	// 5. trigger FTD onboarding state machine
	// TODO

	return nil, nil
}

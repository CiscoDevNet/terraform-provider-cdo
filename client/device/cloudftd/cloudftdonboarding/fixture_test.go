package cloudftdonboarding_test

import (
	"encoding/json"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcplatform"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	fmcconfig2 "github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmcconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmcconfig/fmctaskstatus"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmcdomain"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/user/auth"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/user/auth/role"
)

const (
	baseUrl = "https://unit-test.cdo.cisco.com"

	fmcName      = "unit-test-device-name"
	fmcUid       = "unit-test-uid"
	fmcDomainUid = "unit-test-domain-uid"
	fmcHost      = "unit-test-fmc-host.com"
	fmcPort      = 1234

	fmcDomainLinks  = "unit-test-links"
	fmcDomainCount  = 123
	fmcDomainOffset = 234
	fmcDomainLimit  = 345
	fmcDomainPages  = 456
	fmcDomainName   = "unit-test-name"
	fmcDomainType_  = "unit-test-type"

	ftdName = "unit-test-ftdName"
	ftdUid  = "unit-test-ftdUid"

	ftdGeneratedCommand   = "unit-test-ftdGeneratedCommand"
	ftdAccessPolicyName   = "unit-test-access-policy-item-name"
	ftdNatID              = "unit-test-ftdNatID"
	ftdCloudManagerDomain = "unit-test-ftdCloudManagerDomain.com"
	ftdRegKey             = "unit-test-ftdRegKey"

	tenantUid = "unit-test-tenant-uid"

	systemToken      = "unit-test-system-token"
	systemTokenScope = tenantUid

	fmcCreateDeviceTaskId = "unit-test-task-id"

	ftdSpecificUid = "unit-test-ftd-specific-id"
)

var (
	ftdLicenseCaps     = &[]license.Type{license.Base, license.Carrier}
	ftdPerformanceTier = tier.FTDv5
)

var (
	validReadFtdSpecificOutput = cloudftd.NewReadSpecificOutputBuilder().
					SpecificUid(ftdSpecificUid).
					State(state.DONE).
					Build()

	validUpdateSpecificUidOutput = cloudftd.NewUpdateSpecificFtdOutputBuilder().
					SpecificUid(ftdSpecificUid).
					Build()

	validReadSpecificOutput = cloudftd.NewReadSpecificOutputBuilder().
				SpecificUid(ftdSpecificUid).
				Type(string(devicetype.CloudFtd)).
				State(state.DONE).
				Build()

	validReadTaskOutput = fmcconfig.ReadTaskStatusOutput{
		Status: fmctaskstatus.Success,
	}

	validCreateFmcDeviceRecordOutput = fmcconfig.CreateDeviceRecordOutput{
		Metadata: fmcconfig2.Metadata{
			Task: fmcconfig2.Task{
				Name: "",
				Id:   fmcCreateDeviceTaskId,
				Type: "",
			},
			IsPartOfContainer: false,
			IsMultiInstance:   false,
		},
	}

	validReadFtdOutput = cloudftd.NewReadOutputBuilder().
				Uid(ftdUid).
				Name(ftdName).
				Metadata(cloudftd.NewMetadataBuilder().
					LicenseCaps(ftdLicenseCaps).
					GeneratedCommand(ftdGeneratedCommand).
					AccessPolicyName(ftdAccessPolicyName).
					PerformanceTier(&ftdPerformanceTier).
					NatID(ftdNatID).
					CloudManagerDomain(ftdCloudManagerDomain).
					RegKey(ftdRegKey).
					Build()).
				Build()

	validCreateSystemApiTokenOutput = auth.Token{
		TenantUid:    tenantUid,
		TenantName:   "",
		AccessToken:  systemToken,
		RefreshToken: "",
		TokenType:    "",
		Scope:        "",
	}

	validReadFmcOutput = []device.ReadOutput{
		device.NewReadOutputBuilder().
			AsCloudFmc().
			WithName(fmcName).
			WithUid(fmcUid).
			WithLocation(fmcHost, fmcPort).
			Build(),
	}

	validReadDomainInfo = &fmcplatform.ReadDomainInfoOutput{
		Links:  fmcdomain.NewLinks(fmcDomainLinks),
		Paging: fmcdomain.NewPaging(fmcDomainCount, fmcDomainOffset, fmcDomainLimit, fmcDomainPages),
		Items: []fmcdomain.Item{
			fmcdomain.NewItem(fmcDomainUid, fmcDomainName, fmcDomainType_),
		},
	}

	validReadApiTokenInfo = auth.Info{UserAuthentication: auth.Authentication{
		Authorities: []auth.Authority{
			{Authority: role.Admin},
		},
		Details: auth.Details{
			TenantUid:              tenantUid,
			TenantName:             "",
			SseTenantUid:           "",
			TenantOrganizationName: "",
			TenantDbFeatures:       "",
			TenantUserRoles:        "",
			TenantDatabaseName:     "",
			TenantPayType:          "",
		},
		Authenticated: false,
		Principle:     "",
		Name:          "",
	}}
)

func init() {
	fmt.Println("json.Marshal(validReadApiTokenInfo)")
	t, _ := json.Marshal(validReadApiTokenInfo)
	fmt.Println(string(t))
}

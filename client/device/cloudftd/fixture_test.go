package cloudftd_test

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/accesspolicies"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmcdomain"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
)

const (
	baseUrl = "https://unit-test.net"

	deviceName              = "unit-test-device-name"
	deviceUid               = "unit-test-uid"
	deviceHost              = "https://unit-test.com"
	devicePort              = 1234
	deviceCloudConnectorUId = "unit-test-device-connector-uid"

	fmcDomainUid = "unit-test-domain-uid"
	fmcHost      = "unit-test-fmc-host.com"
	fmcLink      = "unit-test-fmc-link"

	fmcDomainPages    = 123
	fmcDomainCount    = 123
	fmcDomainOffset   = 123
	fmcDomainLimit    = 123
	fmcDomainItemName = "unit-test-fmcDomainItemName"
	fmcDomainItemType = "unit-test-fmcDomainItemType"
	fmcDomainItemUid  = "unit-test-fmcDomainItemUid"

	fmcAccessPolicyPages    = 123
	fmcAccessPolicyCount    = 123
	fmcAccessPolicyOffset   = 123
	fmcAccessPolicyLimit    = 123
	fmcAccessPolicyItemName = "unit-test-access-policy-item-name"
	fmcAccessPolicyItemType = "unit-test-access-policy-item-type"
	fmcAccessPolicyItemUid  = "unit-test-access-policy-item-uid"

	ftdName = "unit-test-ftdName"
	ftdUid  = "unit-test-ftdUid"

	ftdGeneratedCommand   = "unit-test-ftdGeneratedCommand"
	ftdAccessPolicyName   = "unit-test-ftdAccessPolicyName"
	ftdNatID              = "unit-test-ftdNatID"
	ftdCloudManagerDomain = "unit-test-ftdCloudManagerDomain"
	ftdRegKey             = "unit-test-ftdRegKey"
)

var (
	ftdLicenseCaps     = []license.Type{license.Base, license.Carrier}
	ftdPerformanceTier = tier.FTDv5
)

var (
	validReadFmcOutput = device.NewReadOutputBuilder().
				AsCloudFmc().
				WithName(deviceName).
				WithUid(deviceUid).
				WithLocation(deviceHost, devicePort).
				OnboardedUsingCloudConnector(deviceCloudConnectorUId).
				Build()

	validReadFmcDomainInfoOutput = fmcdomain.NewInfoBuilder().
					Links(fmcdomain.NewLinks(fmcLink)).
					Paging(fmcdomain.NewPaging(fmcDomainCount, fmcDomainOffset, fmcDomainLimit, fmcDomainPages)).
					Items([]fmcdomain.Item{
			fmcdomain.NewItem(
				fmcDomainItemUid,
				fmcDomainItemName,
				fmcDomainItemType,
			),
		}).
		Build()

	validReadAccessPoliciesOutput = accesspolicies.NewAccessPoliciesBuilder().
					Links(accesspolicies.NewLinks(fmcLink)).
					Paging(accesspolicies.NewPaging(
			fmcAccessPolicyPages,
			fmcAccessPolicyCount,
			fmcAccessPolicyOffset,
			fmcAccessPolicyLimit,
		)).
		Items([]accesspolicies.Item{accesspolicies.NewItem(
			fmcAccessPolicyItemUid,
			fmcAccessPolicyItemName,
			fmcAccessPolicyItemType,
			accesspolicies.NewLinks(fmcLink),
		)}).
		Build()

	validCreateFtdOutput = cloudftd.NewCreateOutputBuilder().
				Name(ftdName).
				Uid(ftdUid).
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

	validUpdateSpecificFtdOutput = cloudftd.NewUpdateSpecificFtdOutputBuilder().
					SpecificUid(ftdSpecificUid).
					Build()
)

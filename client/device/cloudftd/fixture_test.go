package cloudftd_test

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcappliance"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/statemachine"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/accesspolicies"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmcconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmcdomain"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
)

const (
	baseUrl = "https://unit-test.net"

	fmcName      = "unit-test-device-name"
	fmcUid       = "unit-test-uid"
	fmcDomainUid = "unit-test-domain-uid"
	fmcHost      = "unit-test-fmc-host.com"
	fmcPort      = 1234
	fmcLink      = "unit-test-fmc-link"

	fmcDomainPages    = 123
	fmcDomainCount    = 123
	fmcDomainOffset   = 123
	fmcDomainLimit    = 123
	fmcDomainItemName = "unit-test-fmcDomainItemName"
	fmcDomainItemType = "unit-test-fmcDomainItemType"
	fmcDomainItemUid  = fmcDomainUid

	fmcAccessPolicyPages    = 123
	fmcAccessPolicyCount    = 123
	fmcAccessPolicyOffset   = 123
	fmcAccessPolicyLimit    = 123
	fmcAccessPolicyItemName = "unit-test-access-policy-item-name"
	fmcAccessPolicyItemType = "unit-test-access-policy-item-type"
	fmcAccessPolicyItemUid  = "unit-test-access-policy-item-uid"

	fmcApplianceUid    = "unit-test-fmc-specific-uid"
	fmcApplianceStatus = "unit-test-fmc-specific-status"
	fmcApplianceState  = "unit-test-fmc-specific-state"

	ftdName    = "unit-test-ftdName"
	ftdUid     = "unit-test-ftdUid"
	ftdVirtual = true

	ftdGeneratedCommand   = "unit-test-ftdGeneratedCommand"
	ftdAccessPolicyName   = fmcAccessPolicyItemName
	ftdNatID              = "unit-test-ftdNatID"
	ftdCloudManagerDomain = "unit-test-ftdCloudManagerDomain.com"
	ftdRegKey             = "unit-test-ftdRegKey"

	ftdSpecificUid       = "unit-test-ftd-specific-uid"
	ftdSpecificType      = "unit-test-ftd-specific-type"
	ftdSpecificState     = "unit-test-ftd-specific-state"
	ftdSpecificNamespace = "unit-test-ftd-specific-namespace"

	ftdDeviceRecordId = "unit-test-ftd-device-record-id"
)

var (
	ftdLicenseCaps     = &[]license.Type{license.Base, license.Carrier}
	ftdPerformanceTier = tier.FTDv5
	ftdTags            = tags.New("tags1", "tags2", "tags3")
)

var (
	validReadFmcOutput = []device.ReadOutput{
		device.NewReadOutputBuilder().
			AsCloudFmc().
			WithName(fmcName).
			WithUid(fmcUid).
			WithLocation(fmcHost, fmcPort).
			Build(),
	}

	validFmcDomainItem = fmcdomain.NewItem(
		fmcDomainItemUid,
		fmcDomainItemName,
		fmcDomainItemType,
	)

	validReadFmcDomainInfoOutput = fmcdomain.NewInfoBuilder().
					Links(fmcdomain.NewLinks(fmcLink)).
					Paging(fmcdomain.NewPaging(fmcDomainCount, fmcDomainOffset, fmcDomainLimit, fmcDomainPages)).
					Items([]fmcdomain.Item{validFmcDomainItem}).
					Build()

	validReadAccessPoliciesOutput = accesspolicies.Builder().
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

	validFtdMetadata = cloudftd.NewMetadataBuilder().
				LicenseCaps(ftdLicenseCaps).
				AccessPolicyName(ftdAccessPolicyName).
				PerformanceTier(&ftdPerformanceTier).
				CloudManagerDomain(ftdCloudManagerDomain).
				Build()

	validCreateFtdOutput = cloudftd.NewCreateOutputBuilder().
				Name(ftdName).
				Uid(ftdUid).
				Metadata(validFtdMetadata).
				Build()

	validReadDeviceRecordsOutput = fmcconfig.NewAllDeviceRecordsBuilder().
					Items([]fmcconfig.Item{fmcconfig.NewItem(ftdDeviceRecordId, ftdName, fmcDomainItemType, fmcconfig.NewLinks(""))}).
					Build()

	validReadFtdOutput = cloudftd.NewReadOutputBuilder().
				Name(ftdName).
				Uid(ftdUid).
				Metadata(validFtdMetadata).
				Build()

	validReadFtdDeviceRecordOutput = fmcconfig.DeviceRecord{
		Id:               ftdDeviceRecordId,
		Model:            "Firepower Threat Defense",
		DeploymentStatus: "DEPLOYED",
	}

	validUpdateSpecificFtdOutput = cloudftd.NewUpdateSpecificFtdOutputBuilder().
					SpecificUid(ftdSpecificUid).
					Build()

	validReadFtdSpecificDeviceOutput = cloudftd.NewReadSpecificOutputBuilder().
						SpecificUid(ftdSpecificUid).
						Type(ftdSpecificType).
						State(ftdSpecificState).
						Namespace(ftdSpecificNamespace).
						Build()

	validReadFtdGeneratedCommandOutput = cloudftd.NewCreateOutputBuilder().
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

	validReadSpecificOutput = cloudfmc.ReadSpecificOutput{
		SpecificUid: fmcApplianceUid,
		DomainUid:   fmcDomainUid,
		State:       state.DONE,
		Status:      fmcApplianceStatus,
	}

	validUpdateFmcSpecificOutput = fmcappliance.NewUpdateOutputBuilder().
					Uid(fmcApplianceUid).
					State(fmcApplianceState).
					DomainUid(fmcDomainUid).
					Build()

	validReadStateMachineOutput = statemachine.NewReadInstanceByDeviceUidOutputBuilder().
					StateMachineIdentifier("fmceDeleteFtdcStateMachine").
					StateMachineInstanceCondition(state.DONE).
					Build()

	validDeleteOutput = cloudftd.DeleteOutput{}
)

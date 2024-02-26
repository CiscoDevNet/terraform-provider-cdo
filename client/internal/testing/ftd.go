package testing

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd/cloudftdonboarding"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/publicapilabels"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
)

func (m Model) FtdOnboardingInput() cloudftdonboarding.CreateInput {
	return cloudftdonboarding.CreateInput{
		FtdUid: m.FtdUid.String(),
	}
}

func (m Model) FtdCreateInput() cloudftd.CreateInput {
	return cloudftd.CreateInput{
		Name:             m.FtdName,
		AccessPolicyName: m.FtdAccessPolicyName,
		PerformanceTier:  &m.FtdPerformanceTier,
		Virtual:          false,
		Licenses:         nil,
		Labels:           publicapilabels.Empty(),
	}
}

func (m Model) FtdReadOutput() cloudftd.ReadOutput {
	return cloudftd.ReadOutput{
		Uid:  m.FtdUid.String(),
		Name: m.FtdName,
		Metadata: cloudftd.Metadata{
			AccessPolicyName:   m.FtdAccessPolicyName,
			AccessPolicyUid:    m.FtdAccessPolicyUid.String(),
			CloudManagerDomain: "",
			GeneratedCommand:   "",
			LicenseCaps:        "",
			NatID:              "",
			PerformanceTier:    &m.FtdPerformanceTier,
			RegKey:             "",
		},
		State: "",
		Tags:  tags.Empty(),
	}
}

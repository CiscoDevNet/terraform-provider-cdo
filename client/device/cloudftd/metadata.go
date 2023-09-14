package cloudftd

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
)

type Metadata struct {
	AccessPolicyName   string     `json:"accessPolicyName,omitempty"`
	AccessPolicyUid    string     `json:"accessPolicyUuid,omitempty"`
	CloudManagerDomain string     `json:"cloudManagerDomain,omitempty"`
	GeneratedCommand   string     `json:"generatedCommand,omitempty"`
	LicenseCaps        string     `json:"license_caps,omitempty"`
	NatID              string     `json:"natID,omitempty"`
	PerformanceTier    *tier.Type `json:"performanceTier,omitempty"`
	RegKey             string     `json:"regKey,omitempty"`
}

package cloudftd

import (
	"encoding/json"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
)

type Metadata struct {
	AccessPolicyName   string         `json:"accessPolicyName,omitempty"`
	AccessPolicyUid    string         `json:"accessPolicyUuid,omitempty"`
	CloudManagerDomain string         `json:"cloudManagerDomain,omitempty"`
	GeneratedCommand   string         `json:"generatedCommand,omitempty"`
	LicenseCaps        []license.Type `json:"license_caps,omitempty"`
	NatID              string         `json:"natID,omitempty"`
	PerformanceTier    *tier.Type     `json:"performanceTier,omitempty"`
	RegKey             string         `json:"regKey,omitempty"`
}

type internalMetadata struct {
	AccessPolicyName   string     `json:"accessPolicyName,omitempty"`
	AccessPolicyUuid   string     `json:"accessPolicyUuid,omitempty"`
	CloudManagerDomain string     `json:"cloudManagerDomain,omitempty"`
	GeneratedCommand   string     `json:"generatedCommand,omitempty"`
	LicenseCaps        string     `json:"license_caps,omitempty"` // first, unmarshal it into string
	NatID              string     `json:"natID,omitempty"`
	PerformanceTier    *tier.Type `json:"performanceTier,omitempty"`
	RegKey             string     `json:"regKey,omitempty"`
}

// UnmarshalJSON defines custom unmarshal json for metadata, because we need to handle license caps differently,
// it is a string containing command separated values, instead of a json list where it can be parsed directly.
// Note that this method is defined on the *Metadata type, so if you unmarshal or marshal a Metadata without pointer,
// it will not be called.
func (metadata *Metadata) UnmarshalJSON(data []byte) error {
	fmt.Printf("\nunmarshalling Metadata: %s\n", string(data))
	var internalMeta internalMetadata
	err := json.Unmarshal(data, &internalMeta)
	if err != nil {
		return err
	}

	licenseCaps, err := license.DeserializeAll(internalMeta.LicenseCaps) // now parse it into golang type
	if err != nil {
		return err
	}

	(*metadata).AccessPolicyName = internalMeta.AccessPolicyName
	(*metadata).AccessPolicyUid = internalMeta.AccessPolicyUuid
	(*metadata).CloudManagerDomain = internalMeta.CloudManagerDomain
	(*metadata).GeneratedCommand = internalMeta.GeneratedCommand
	(*metadata).NatID = internalMeta.NatID
	(*metadata).PerformanceTier = internalMeta.PerformanceTier
	(*metadata).RegKey = internalMeta.RegKey

	(*metadata).LicenseCaps = licenseCaps // set it as usual

	return nil
}

func (metadata *Metadata) MarshalJSON() ([]byte, error) {
	fmt.Printf("\nmarshalling Metadata: %+v\n", metadata)
	var internalMeta internalMetadata
	internalMeta.AccessPolicyName = metadata.AccessPolicyName
	internalMeta.AccessPolicyUuid = metadata.AccessPolicyUid
	internalMeta.CloudManagerDomain = metadata.CloudManagerDomain
	internalMeta.GeneratedCommand = metadata.GeneratedCommand
	internalMeta.LicenseCaps = license.SerializeAll(metadata.LicenseCaps)
	internalMeta.NatID = metadata.NatID
	internalMeta.PerformanceTier = metadata.PerformanceTier
	internalMeta.RegKey = metadata.RegKey

	return json.Marshal(internalMeta)
}

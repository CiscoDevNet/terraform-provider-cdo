package cloudftd

import (
	"context"
	"encoding/json"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
)

type ReadByUidInput struct {
	Uid string `json:"uid"`
}

func NewReadByUidInput(uid string) ReadByUidInput {
	return ReadByUidInput{
		Uid: uid,
	}
}

type ReadOutput struct {
	Uid      string   `json:"uid"`
	Name     string   `json:"name"`
	Metadata Metadata `json:"metadata,omitempty"`
}

type Metadata struct {
	AccessPolicyName   string         `json:"accessPolicyName,omitempty"`
	AccessPolicyUuid   string         `json:"accessPolicyUuid,omitempty"`
	CloudManagerDomain string         `json:"cloudManagerDomain,omitempty"`
	GeneratedCommand   string         `json:"generatedCommand,omitempty"`
	LicenseCaps        []license.Type `json:"license_caps,omitempty"`
	NatID              string         `json:"natID,omitempty"`
	PerformanceTier    *tier.Type     `json:"performanceTier,omitempty"`
	RegKey             string         `json:"regKey,omitempty"`
}

// custom unmarshal json for metadata, because we need to handle license caps differently,
// it is a string containing command separated values, instead of a json list where it can be parsed directly.
func (metadata *Metadata) UnmarshalJSON(data []byte) error {
	var t struct {
		AccessPolicyName   string     `json:"accessPolicyName,omitempty"`
		AccessPolicyUuid   string     `json:"accessPolicyUuid,omitempty"`
		CloudManagerDomain string     `json:"cloudManagerDomain,omitempty"`
		GeneratedCommand   string     `json:"generatedCommand,omitempty"`
		LicenseCaps        string     `json:"license_caps,omitempty"` // first, unmarshal it into string
		NatID              string     `json:"natID,omitempty"`
		PerformanceTier    *tier.Type `json:"performanceTier,omitempty"`
		RegKey             string     `json:"regKey,omitempty"`
	}
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	licenseCaps, err := license.ParseAll(t.LicenseCaps) // now parse it into golang type
	if err != nil {
		return err
	}

	(*metadata).AccessPolicyName = t.AccessPolicyName
	(*metadata).AccessPolicyUuid = t.AccessPolicyUuid
	(*metadata).CloudManagerDomain = t.CloudManagerDomain
	(*metadata).GeneratedCommand = t.GeneratedCommand
	(*metadata).NatID = t.NatID
	(*metadata).PerformanceTier = t.PerformanceTier
	(*metadata).RegKey = t.RegKey

	(*metadata).LicenseCaps = licenseCaps // set it as usual

	return nil
}

func ReadByUid(ctx context.Context, client http.Client, readInp ReadByUidInput) (*ReadOutput, error) {

	readUrl := url.ReadDevice(client.BaseUrl(), readInp.Uid)
	req := client.NewGet(ctx, readUrl)

	var readOutp ReadOutput
	if err := req.Send(&readOutp); err != nil {
		return nil, err
	}

	return &readOutp, nil
}

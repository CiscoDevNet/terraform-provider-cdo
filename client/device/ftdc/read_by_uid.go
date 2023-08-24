package ftdc

import (
	"context"
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

type ReadByUidOutput struct {
	Uid      string   `json:"uid"`
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

func ReadByUid(ctx context.Context, client http.Client, readInp ReadByUidInput) (*ReadByUidOutput, error) {

	readUrl := url.ReadDevice(client.BaseUrl(), readInp.Uid)
	req := client.NewGet(ctx, readUrl)

	var readOutp ReadByUidOutput
	if err := req.Send(&readOutp); err != nil {
		return nil, err
	}

	return &readOutp, nil
}

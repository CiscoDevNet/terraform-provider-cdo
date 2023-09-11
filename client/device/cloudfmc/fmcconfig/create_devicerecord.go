package fmcconfig

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmcconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
)

type CreateDeviceRecordInput struct {
	FmcDomainUid    string
	FmcHostname     string
	Name            string
	NatId           string
	RegKey          string
	PerformanceTier *tier.Type
	LicenseCaps     *[]license.Type
	AccessPolicyUid string
	Type            string
	SystemApiToken  string // normal cdo token does not work for this request, a cdo system token is needed
}

type createDeviceRecordRequestBody struct {
	Name            string          `json:"name"`
	NatId           string          `json:"natID"`
	RegKey          string          `json:"regKey"`
	PerformanceTier *tier.Type      `json:"performanceTier"`
	Type            string          `json:"type"`
	LicenseCaps     *[]license.Type `json:"license_caps"`
	AccessPolicy    accessPolicy    `json:"accessPolicy"`
}

type accessPolicy struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

func NewCreateDeviceRecordInput(fmcDomainUid string) CreateDeviceRecordInput {
	return CreateDeviceRecordInput{
		FmcDomainUid: fmcDomainUid,
	}
}

type CreateDeviceRecordOutput = fmcconfig.DeviceRecordCreation

func CreateDeviceRecord(ctx context.Context, client http.Client, createInp CreateDeviceRecordInput) (*CreateDeviceRecordOutput, error) {

	createUrl := url.CreateFmcDeviceRecord(client.BaseUrl(), createInp.FmcDomainUid)

	body := createDeviceRecordRequestBody{
		Name:            createInp.Name,
		NatId:           createInp.NatId,
		RegKey:          createInp.RegKey,
		PerformanceTier: createInp.PerformanceTier,
		Type:            createInp.Type,
		LicenseCaps:     createInp.LicenseCaps,
		AccessPolicy: accessPolicy{
			Id:   createInp.AccessPolicyUid,
			Type: "AccessPolicy",
		},
	}
	req := client.NewPost(ctx, createUrl, body)
	req.Header.Add("Fmc-Hostname", createInp.FmcHostname)
	req.OverrideApiToken(createInp.SystemApiToken)

	var createOutp fmcconfig.DeviceRecordCreation
	if err := req.Send(&createOutp); err != nil {
		return nil, err
	}

	return &createOutp, nil
}

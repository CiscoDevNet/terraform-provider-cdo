package ftdc

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cdfmc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/cdo"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/sliceutil"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
	"strings"
)

type CreateInput struct {
	Name             string
	AccessPolicyName string
	PerformanceTier  *tier.Type // ignored if it is physical device
	Virtual          bool
	Licenses         []license.Type
}

type CreateOutput struct {
	Uid      string
	Name     string
	Metadata Metadata
}

func NewCreateInput(
	name string,
	accessPolicyName string,
	performanceTier *tier.Type,
	virtual bool,
	licenses []license.Type,
) CreateInput {
	return CreateInput{
		Name:             name,
		AccessPolicyName: accessPolicyName,
		PerformanceTier:  performanceTier,
		Virtual:          virtual,
		Licenses:         licenses,
	}
}

type createRequestBody struct {
	FmcId      string          `json:"associatedDeviceUid"`
	DeviceType devicetype.Type `json:"deviceType"`
	Metadata   metadata        `json:"metadata"`
	Name       string          `json:"name"`
	State      string          `json:"state"` // TODO: use queueTriggerState?
	Type       string          `json:"type"`
}

type metadata struct {
	AccessPolicyName string     `json:"accessPolicyName"`
	AccessPolicyId   string     `json:"accessPolicyUuid"`
	LicenseCaps      string     `json:"license_caps"`
	PerformanceTier  *tier.Type `json:"performanceTier"`
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	client.Logger.Println("creating ftdc")

	// 1. find cdFMC
	fmcRes, err := cdfmc.Read(ctx, client, cdfmc.NewReadInput())
	if err != nil {
		return nil, err
	}
	// 2. get cdFMC domain id by looking up FMC's specific device
	fmcSpecificRes, err := cdfmc.ReadSpecific(ctx, client, cdfmc.NewReadSpecificInput(fmcRes.Uid))
	if err != nil {
		return nil, err
	}

	// 3. read access policies using cdFMC domain id
	accessPoliciesRes, err := cdfmc.ReadAccessPolicies(
		ctx,
		client,
		cdfmc.NewReadAccessPoliciesInput(fmcSpecificRes.DomainUid, 1000), // 1000 is what CDO UI uses
	)
	if err != nil {
		return nil, err
	}
	selectedPolicy, ok := accessPoliciesRes.Find(createInp.AccessPolicyName)
	if !ok {
		return nil, fmt.Errorf(
			`access policy: "%s" not found, available policies: %s. In rare cases where you have more than 1000 access policies, please raise an issue at: %s`,
			createInp.AccessPolicyName,
			accessPoliciesRes.Items,
			cdo.TerraformProviderCDOIssuesUrl,
		)
	}

	// handle selected license caps
	licenseCaps := sliceutil.Map(createInp.Licenses, func(l license.Type) string { return string(l) })

	// handle performance tier
	var performanceTier *tier.Type = nil // physical is nil
	if createInp.Virtual {
		performanceTier = createInp.PerformanceTier
	}

	// 4. create the ftdc device
	createUrl := url.CreateDevice(client.BaseUrl())
	createBody := createRequestBody{
		Name:       createInp.Name,
		FmcId:      fmcRes.Uid,
		DeviceType: devicetype.Ftdc,
		Metadata: metadata{
			AccessPolicyName: selectedPolicy.Name,
			AccessPolicyId:   selectedPolicy.Id,
			LicenseCaps:      strings.Join(licenseCaps, ","),
			PerformanceTier:  performanceTier,
		},
		State: "NEW",
		Type:  "devices",
	}
	createReq := client.NewPost(ctx, createUrl, createBody)
	var createOup CreateOutput
	if err := createReq.Send(&createOup); err != nil {
		return nil, err
	}

	// 5. read created ftdc's specific device's uid
	readSpecRes, err := device.ReadSpecific(ctx, client, *device.NewReadSpecificInput(createOup.Uid))
	if err != nil {
		return nil, err
	}

	// 6. initiate ftdc onboarding by triggering a weird endpoint using created ftdc's specific uid
	_, err = UpdateSpecificFtd(ctx, client,
		NewUpdateSpecificFtdInput(
			readSpecRes.SpecificUid,
			"INITIATE_FTDC_ONBOARDING",
		),
	)

	// 7. get generate command
	readOutp, err := ReadByUid(ctx, client, NewReadByUidInput(createOup.Uid))
	if err != nil {
		return nil, err
	}

	// done!
	return &CreateOutput{
		Uid:      createOup.Uid,
		Name:     createOup.Name,
		Metadata: readOutp.Metadata,
	}, nil
}

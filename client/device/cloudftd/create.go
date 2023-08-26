package cloudftd

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcplatform"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/cdo"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
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
	Uid      string   `json:"uid"`
	Name     string   `json:"name"`
	Metadata Metadata `json:"metadata"`
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
	Model      bool            `json:"model"`
}

type metadata struct {
	AccessPolicyName string     `json:"accessPolicyName"`
	AccessPolicyId   string     `json:"accessPolicyUuid"`
	LicenseCaps      string     `json:"license_caps"`
	PerformanceTier  *tier.Type `json:"performanceTier"`
}

//const FmcDomainUid = "e276abec-e0f2-11e3-8169-6d9ed49b625f"

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	client.Logger.Println("creating cloudftd")

	// 1. find Cloud FMC
	fmcRes, err := cloudfmc.Read(ctx, client, cloudfmc.NewReadInput())
	if err != nil {
		return nil, err
	}

	// 2. get FMC domain uid by reading Cloud FMC domain info
	readFmcDomainRes, err := fmcplatform.ReadFmcDomainInfo(ctx, client, fmcplatform.NewReadDomainInfo(fmcRes.Host))
	if err != nil {
		return nil, err
	}
	if len(readFmcDomainRes.Items) == 0 {
		return nil, fmt.Errorf("fmc domain info not found")
	}

	// 3. read access policies using Cloud FMC domain id
	accessPoliciesRes, err := cloudfmc.ReadAccessPolicies(
		ctx,
		client,
		cloudfmc.NewReadAccessPoliciesInput(fmcRes.Host, readFmcDomainRes.Items[0].Uuid, 1000), // 1000 is what CDO UI uses
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

	// 4. create the cloud ftd device
	createUrl := url.CreateDevice(client.BaseUrl())
	createBody := createRequestBody{
		Name:       createInp.Name,
		FmcId:      fmcRes.Uid,
		DeviceType: devicetype.CloudFtd,
		Metadata: metadata{
			AccessPolicyName: selectedPolicy.Name,
			AccessPolicyId:   selectedPolicy.Id,
			LicenseCaps:      strings.Join(licenseCaps, ","),
			PerformanceTier:  performanceTier,
		},
		State: "NEW",
		Type:  "devices",
		Model: false,
	}
	createReq := client.NewPost(ctx, createUrl, createBody)
	var createOup CreateOutput
	if err := createReq.Send(&createOup); err != nil {
		return nil, err
	}

	// 5. read created cloudftd's specific device's uid
	readSpecRes, err := device.ReadSpecific(ctx, client, *device.NewReadSpecificInput(createOup.Uid))
	if err != nil {
		return nil, err
	}

	// 6. initiate cloudftd onboarding by triggering a weird endpoint using created cloudftd's specific uid
	_, err = UpdateSpecific(ctx, client,
		NewUpdateSpecificFtdInput(
			readSpecRes.SpecificUid,
			"INITIATE_FTDC_ONBOARDING",
		),
	)

	// 7. get generate command
	err = retry.Do(UntilGeneratedCommandAvailable(ctx, client, createOup.Uid), *retry.NewOptionsWithLoggerAndRetries(client.Logger, 3))
	//readOutp, err := ReadByUid(ctx, client, NewReadByUidInput(createOup.SpecificUid))
	if err != nil {
		return nil, err
	}
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

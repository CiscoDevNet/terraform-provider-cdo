package ftd

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

func Read(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// do read
	inp := cloudftd.NewReadByNameInput(stateData.Name.ValueString())
	res, err := resource.client.ReadCloudFtdByName(ctx, inp)
	if err != nil {
		return err
	}

	// map return struct to model
	stateData.ID = types.StringValue(res.Uid)
	stateData.Name = types.StringValue(res.Name)
	stateData.AccessPolicyName = types.StringValue(res.Metadata.AccessPolicyName)
	stateData.AccessPolicyUid = types.StringValue(res.Metadata.AccessPolicyUid)
	stateData.Virtual = types.BoolValue(res.Metadata.PerformanceTier != nil)
	stateData.Licenses = util.GoStringSliceToTFStringList(strings.Split(res.Metadata.LicenseCaps, ","))
	if res.Metadata.PerformanceTier != nil { // nil means physical cloudftd
		stateData.PerformanceTier = types.StringValue(string(*res.Metadata.PerformanceTier))
	}
	stateData.GeneratedCommand = types.StringValue(res.Metadata.GeneratedCommand)
	stateData.Hostname = types.StringValue(res.Metadata.CloudManagerDomain)
	stateData.NatId = types.StringValue(res.Metadata.NatID)
	stateData.RegKey = types.StringValue(res.Metadata.RegKey)

	return nil
}

func Create(ctx context.Context, resource *Resource, planData *ResourceModel) error {

	// do create
	var performanceTier *tier.Type = nil
	if planData.PerformanceTier.ValueString() != "" {
		t, err := tier.Parse(planData.PerformanceTier.ValueString())
		performanceTier = &t
		if err != nil {
			return err
		}
	}

	licensesGoList := util.TFStringListToGoStringList(planData.Licenses)
	licenses, err := license.DeserializeAllFromCdo(strings.Join(licensesGoList, ","))
	if err != nil {
		return err
	}
	createInp := cloudftd.NewCreateInput(
		planData.Name.ValueString(),
		planData.AccessPolicyName.ValueString(),
		performanceTier,
		planData.Virtual.ValueBool(),
		&licenses,
	)
	res, err := resource.client.CreateCloudFtd(ctx, createInp)
	if err != nil {
		return err
	}

	// map return struct to model
	planData.ID = types.StringValue(res.Uid)
	planData.Name = types.StringValue(res.Name)
	planData.AccessPolicyName = types.StringValue(res.Metadata.AccessPolicyName)
	planData.AccessPolicyUid = types.StringValue(res.Metadata.AccessPolicyUid)
	planData.Licenses = util.GoStringSliceToTFStringList(strings.Split(res.Metadata.LicenseCaps, ","))
	if res.Metadata.PerformanceTier != nil { // nil means physical cloud ftd
		planData.PerformanceTier = types.StringValue(string(*res.Metadata.PerformanceTier))
	}
	planData.GeneratedCommand = types.StringValue(res.Metadata.GeneratedCommand)
	planData.Hostname = types.StringValue(res.Metadata.CloudManagerDomain)
	planData.NatId = types.StringValue(res.Metadata.NatID)
	planData.RegKey = types.StringValue(res.Metadata.RegKey)

	return nil
}

func Update(ctx context.Context, resource *Resource, planData *ResourceModel, stateData *ResourceModel) error {

	// do update
	inp := cloudftd.NewUpdateInput(planData.ID.ValueString(), planData.Name.ValueString())
	res, err := resource.client.UpdateCloudFtd(ctx, inp)
	if err != nil {
		return err
	}

	// map return struct to model
	stateData.Name = types.StringValue(res.Name)

	return nil
}

func Delete(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// do delete
	inp := cloudftd.NewDeleteInput(stateData.ID.ValueString())
	_, err := resource.client.DeleteCloudFtd(ctx, inp)

	return err
}

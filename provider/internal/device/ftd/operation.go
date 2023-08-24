package ftd

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/ftdc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/sliceutil"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Read(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// do read
	inp := ftdc.NewReadByNameInput(stateData.Name.ValueString())
	res, err := resource.client.ReadFtdcByName(ctx, inp)
	if err != nil {
		return err
	}

	// map return struct to sdc model
	stateData.ID = types.StringValue(res.Uid)
	stateData.Name = types.StringValue(res.Name)
	stateData.AccessPolicyName = types.StringValue(res.Metadata.AccessPolicyName)
	stateData.AccessPolicyUid = types.StringValue(res.Metadata.AccessPolicyUuid)
	stateData.Virtual = types.BoolValue(res.Metadata.PerformanceTier != nil)
	stateData.Licenses = util.GoStringSliceToTFStringList(sliceutil.Map(res.Metadata.LicenseCaps, func(l license.Type) string { return string(l) }))
	if res.Metadata.PerformanceTier != nil { // nil means physical ftd
		stateData.PerformanceTier = types.StringValue(string(*res.Metadata.PerformanceTier))
	}
	stateData.GeneratedCommand = types.StringValue(res.Metadata.GeneratedCommand)

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

	createInp := ftdc.NewCreateInput(
		planData.Name.ValueString(),
		planData.AccessPolicyName.ValueString(),
		performanceTier,
		planData.Virtual.ValueBool(),
		sliceutil.Map(util.TFStringListToGoStringList(planData.Licenses), func(s string) license.Type { return license.MustParse(s) }),
	)
	res, err := resource.client.CreateFtdc(ctx, createInp)
	if err != nil {
		return err
	}

	// map return struct to sdc model
	planData.ID = types.StringValue(res.Uid)
	planData.Name = types.StringValue(res.Name)
	planData.AccessPolicyName = types.StringValue(res.Metadata.AccessPolicyName)
	planData.AccessPolicyUid = types.StringValue(res.Metadata.AccessPolicyUuid)
	planData.Virtual = types.BoolValue(res.Metadata.PerformanceTier != nil)
	planData.Licenses = util.GoStringSliceToTFStringList(sliceutil.Map(res.Metadata.LicenseCaps, func(l license.Type) string { return string(l) }))
	if res.Metadata.PerformanceTier != nil { // nil means physical ftd
		planData.PerformanceTier = types.StringValue(string(*res.Metadata.PerformanceTier))
	}
	planData.GeneratedCommand = types.StringValue(res.Metadata.GeneratedCommand)

	return nil
}

func Update(ctx context.Context, resource *Resource, planData *ResourceModel, stateData *ResourceModel) error {
	// TODO: fill me

	// do update

	// map return struct to sdc model

	return nil
}

func Delete(ctx context.Context, resource *Resource, stateData *ResourceModel) error {
	// TODO: fill me

	// do delete

	return nil
}

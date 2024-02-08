package duoadminpanel

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/duoadminpanel"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Read(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	output, err := resource.client.ReadDuoAdminPanel(ctx, duoadminpanel.ReadByUidInput{
		Uid: stateData.Id.ValueString(),
	})
	if err != nil {
		return err
	}
	stateData.Labels = util.GoStringSliceToTFStringSet(output.Tags.Labels)
	stateData.Name = types.StringValue(output.Name)

	return nil
}

func Create(ctx context.Context, resource *Resource, planData *ResourceModel) error {

	// convert tf tags to go tags
	planTags, err := util.TFStringSetToTagLabels(ctx, planData.Labels)
	if err != nil {
		return fmt.Errorf("error while converting terraform tags to go tags, %s", planData.Labels)
	}

	// do create
	output, err := resource.client.CreateDuoAdminPanel(ctx, duoadminpanel.CreateInput{
		Name:           planData.Name.ValueString(),
		Host:           planData.Host.ValueString(),
		IntegrationKey: planData.IntegrationKey.ValueString(),
		SecretKey:      planData.SecretKey.ValueString(),
		Labels:         planTags.Labels,
	})
	if err != nil {
		return err
	}

	// map response to terraform types, it should be unchanged for most parts
	planData.Id = types.StringValue(output.Uid)
	planData.Labels = util.GoStringSliceToTFStringSet(output.Tags.Labels)

	return nil
}

func Update(ctx context.Context, resource *Resource, planData *ResourceModel, stateData *ResourceModel) error {

	// do update
	planTags, err := util.TFStringSetToTagLabels(ctx, planData.Labels)
	if err != nil {
		return fmt.Errorf("error while converting terraform tags to go tags, %s", planData.Labels)
	}

	output, err := resource.client.UpdateDuoAdminPanel(ctx, duoadminpanel.UpdateInput{
		Uid:  stateData.Id.ValueString(),
		Name: planData.Name.ValueString(),
		Tags: planTags,
	})
	if err != nil {
		return err
	}
	stateData.Labels = util.GoStringSliceToTFStringSet(output.Tags.Labels)
	stateData.Name = types.StringValue(output.Name)

	return nil
}

func Delete(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// do delete
	_, err := resource.client.DeleteDuoAdminPanel(ctx, duoadminpanel.DeleteInput{
		Uid: stateData.Id.ValueString(),
	})
	if err != nil {
		return err
	}

	return nil
}

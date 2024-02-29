package duoadminpanel

import (
	"context"
	"errors"
	"fmt"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/duoadminpanel"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/publicapilabels"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
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
	stateData.Labels = util.GoStringSliceToTFStringSet(output.Tags.UngroupedTags())
	stateData.GroupedLabels = util.GoMapToStringSetTFMap(output.Tags.GroupedTags())
	stateData.Name = types.StringValue(output.Name)

	return nil
}

func Create(ctx context.Context, resource *Resource, planData *ResourceModel) error {

	labels, err := publicapiLabelsFromResourceModel(ctx, planData)
	if err != nil {
		return err
	}

	// do create
	output, err := resource.client.CreateDuoAdminPanel(ctx, duoadminpanel.CreateInput{
		Name:           planData.Name.ValueString(),
		Host:           planData.Host.ValueString(),
		IntegrationKey: planData.IntegrationKey.ValueString(),
		SecretKey:      planData.SecretKey.ValueString(),
		Labels:         labels,
	})
	if err != nil {
		return err
	}

	// map response to terraform types, it should be unchanged for most parts
	planData.Id = types.StringValue(output.Uid)
	planData.Labels = util.GoStringSliceToTFStringSet(output.Tags.UngroupedTags())
	planData.GroupedLabels = util.GoMapToStringSetTFMap(output.Tags.GroupedTags())

	return nil
}

func Update(ctx context.Context, resource *Resource, planData *ResourceModel, stateData *ResourceModel) error {

	// do update
	planTags, err := tagsFromResourceModel(ctx, stateData)
	if err != nil {
		return err
	}

	output, err := resource.client.UpdateDuoAdminPanel(ctx, duoadminpanel.UpdateInput{
		Uid:  stateData.Id.ValueString(),
		Name: planData.Name.ValueString(),
		Tags: planTags,
	})
	if err != nil {
		return err
	}
	stateData.Labels = util.GoStringSliceToTFStringSet(output.Tags.UngroupedTags())
	stateData.GroupedLabels = util.GoMapToStringSetTFMap(output.Tags.GroupedTags())
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

func publicapiLabelsFromResourceModel(ctx context.Context, resourceModel *ResourceModel) (publicapilabels.Type, error) {
	ungroupedLabels, groupedLabels, err := extractLabelsFromResourceModel(ctx, resourceModel)
	if err != nil {
		return publicapilabels.Empty(), err
	}

	return publicapilabels.New(ungroupedLabels, groupedLabels), nil
}

func tagsFromResourceModel(ctx context.Context, resourceModel *ResourceModel) (tags.Type, error) {
	ungroupedLabels, groupedLabels, err := extractLabelsFromResourceModel(ctx, resourceModel)
	if err != nil {
		return nil, err
	}

	return tags.New(ungroupedLabels, groupedLabels), nil
}

func extractLabelsFromResourceModel(ctx context.Context, resourceModel *ResourceModel) ([]string, map[string][]string, error) {
	if resourceModel == nil {
		return nil, nil, errors.New("resource model cannot be nil")
	}

	ungroupedLabels, err := util.TFStringSetToGoStringList(ctx, resourceModel.Labels)
	if err != nil {
		return nil, nil, fmt.Errorf("error while converting terraform labels to go slice, %s", resourceModel.Labels)
	}

	groupedLabels, err := util.TFMapToGoMapOfStringSlices(ctx, resourceModel.GroupedLabels)
	if err != nil {
		return nil, nil, fmt.Errorf("error while converting terraform grouped labels to go map, %v", resourceModel.GroupedLabels)
	}

	return ungroupedLabels, groupedLabels, nil
}

package ios

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/publicapilabels"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/ios"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func Read(ctx context.Context, resource *IosDeviceResource, stateData *IosDeviceResourceModel) error {

	readInp := ios.ReadInput{
		Uid: stateData.ID.ValueString(),
	}

	readOutp, err := resource.client.ReadIos(ctx, readInp)
	if err != nil {
		return err
	}

	port, err := strconv.ParseInt(readOutp.Port, 10, 16)
	if err != nil {
		return err
	}

	stateData.Port = types.Int64Value(port)
	stateData.ID = types.StringValue(readOutp.Uid)
	stateData.Name = types.StringValue(readOutp.Name)
	stateData.Ipv4 = types.StringValue(readOutp.SocketAddress)
	stateData.Host = types.StringValue(readOutp.Host)
	stateData.IgnoreCertificate = types.BoolValue(readOutp.IgnoreCertificate)
	stateData.Labels = util.GoStringSliceToTFStringSet(readOutp.Tags.UngroupedTags())
	stateData.GroupedLabels = util.GoMapToStringSetTFMap(readOutp.Tags.GroupedTags())

	return nil
}

func Create(ctx context.Context, resource *IosDeviceResource, planData *IosDeviceResourceModel) error {

	readSdcByNameInp := connector.NewReadByNameInput(
		planData.ConnectorName.ValueString(),
	)

	readSdcOutp, err := resource.client.ReadConnectorByName(ctx, *readSdcByNameInp)
	if err != nil {
		return err
	}

	// convert tf tags to go tags
	planTags, err := labelsFromIosDeviceResourceModel(ctx, planData)
	if err != nil {
		return err
	}

	createInp := ios.NewCreateRequestInput(
		planData.Name.ValueString(),
		readSdcOutp.Uid,
		"SDC",
		planData.Ipv4.ValueString(),
		planData.Username.ValueString(),
		planData.Password.ValueString(),
		planData.IgnoreCertificate.ValueBool(),
		planTags,
	)

	createOutp, err := resource.client.CreateIos(ctx, *createInp)
	tflog.Debug(ctx, fmt.Sprintf("Creation error: %v", err))
	if err != nil {
		return err
	}

	planData.ID = types.StringValue(createOutp.Uid)
	planData.ConnectorName = types.StringValue(planData.ConnectorName.ValueString())
	planData.Name = types.StringValue(createOutp.Name)
	planData.Host = types.StringValue(createOutp.Host)

	port, err := strconv.ParseInt(createOutp.Port, 10, 16)
	if err != nil {
		return fmt.Errorf("failed to parse IOS port, cause=%w", err)
	}
	planData.Port = types.Int64Value(port)
	planData.Labels = util.GoStringSliceToTFStringSet(createOutp.Tags.UngroupedTags())
	planData.GroupedLabels = util.GoMapToStringSetTFMap(createOutp.Tags.GroupedTags())

	return nil
}

func Update(ctx context.Context, resource *IosDeviceResource, planData *IosDeviceResourceModel, stateData *IosDeviceResourceModel) error {

	// convert tf tags to go tags
	planTags, err := tagsFromIosDeviceResourceModel(ctx, planData)
	if err != nil {
		return err
	}

	updateInp := *ios.NewUpdateInput(
		stateData.ID.ValueString(),
		planData.Name.ValueString(),
		planTags,
	)
	updateOutp, err := resource.client.UpdateIos(ctx, updateInp)
	if err != nil {
		return err
	}
	stateData.Name = types.StringValue(updateOutp.Name)
	stateData.Labels = planData.Labels
	stateData.GroupedLabels = planData.GroupedLabels

	return nil
}

func Delete(ctx context.Context, resource *IosDeviceResource, stateData *IosDeviceResourceModel) error {
	deleteInp := ios.NewDeleteInput(stateData.ID.ValueString())
	_, err := resource.client.DeleteIos(ctx, *deleteInp)
	return err
}

func ungroupedAndGroupedLabelsFromIosDeviceResourceModel(ctx context.Context, resourceModel *IosDeviceResourceModel) ([]string, map[string][]string, error) {
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

func tagsFromIosDeviceResourceModel(ctx context.Context, resourceModel *IosDeviceResourceModel) (tags.Type, error) {
	ungroupedLabels, groupedLabels, err := ungroupedAndGroupedLabelsFromIosDeviceResourceModel(ctx, resourceModel)
	if err != nil {
		return nil, err
	}

	return tags.New(ungroupedLabels, groupedLabels), nil
}

func labelsFromIosDeviceResourceModel(ctx context.Context, resourceModel *IosDeviceResourceModel) (publicapilabels.Type, error) {
	ungroupedLabels, groupedLabels, err := ungroupedAndGroupedLabelsFromIosDeviceResourceModel(ctx, resourceModel)
	if err != nil {
		return publicapilabels.Empty(), err
	}

	return publicapilabels.New(ungroupedLabels, groupedLabels), nil
}

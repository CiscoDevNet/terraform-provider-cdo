package ios

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"strconv"

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
	stateData.IgnoreCertifcate = types.BoolValue(readOutp.IgnoreCertifcate)

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

	createInp := ios.NewCreateRequestInput(
		planData.Name.ValueString(),
		readSdcOutp.Uid,
		"SDC",
		planData.Ipv4.ValueString(),
		planData.Username.ValueString(),
		planData.Password.ValueString(),
		planData.IgnoreCertifcate.ValueBool(),
	)

	createOutp, createErr := resource.client.CreateIos(ctx, *createInp)
	tflog.Debug(ctx, fmt.Sprintf("Creation error: %v", createErr))
	if createErr != nil {
		if createErr.CreatedResourceId != nil {
			deleteInp := ios.NewDeleteInput(*createErr.CreatedResourceId)
			_, deletionErr := resource.client.DeleteIos(ctx, *deleteInp)
			if deletionErr != nil {
				tflog.Error(ctx, "Failed to delete iOS device that we failed to create")
				return deletionErr
			}
		}

		return createErr.Err
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

	return nil
}

func Update(ctx context.Context, resource *IosDeviceResource, planData *IosDeviceResourceModel, stateData *IosDeviceResourceModel) error {
	updateInp := *ios.NewUpdateInput(
		stateData.ID.ValueString(),
		planData.Name.ValueString(),
	)
	updateOutp, err := resource.client.UpdateIos(ctx, updateInp)
	if err != nil {
		return err
	}
	stateData.Name = types.StringValue(updateOutp.Name)
	return nil
}

func Delete(ctx context.Context, resource *IosDeviceResource, stateData *IosDeviceResourceModel) error {
	deleteInp := ios.NewDeleteInput(stateData.ID.ValueString())
	_, err := resource.client.DeleteIos(ctx, *deleteInp)
	return err
}

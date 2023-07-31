package ios

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cisco-lockhart/go-client/connector/sdc"
	"github.com/cisco-lockhart/go-client/device/ios"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	stateData.SdcType = types.StringValue(readOutp.LarType)
	stateData.Name = types.StringValue(readOutp.Name)
	stateData.Ipv4 = types.StringValue(readOutp.Ipv4)
	stateData.Host = types.StringValue(readOutp.Host)
	stateData.IgnoreCertifcate = types.BoolValue(readOutp.IgnoreCertifcate)

	return nil
}

func Create(ctx context.Context, resource *IosDeviceResource, planData *IosDeviceResourceModel) error {

	readSdcByNameInp := sdc.NewReadByNameInput(
		planData.SdcName.ValueString(),
	)

	readSdcOutp, err := resource.client.ReadSdcByName(ctx, *readSdcByNameInp)
	if err != nil {
		return err
	}

	createInp := ios.NewCreateRequestInput(
		planData.Name.ValueString(),
		readSdcOutp.Uid,
		planData.SdcType.ValueString(),
		planData.Ipv4.ValueString(),
		planData.Username.ValueString(),
		planData.Password.ValueString(),
		planData.IgnoreCertifcate.ValueBool(),
	)

	createOutp, err := resource.client.CreateIos(ctx, *createInp)
	if err != nil {
		return err
	}

	planData.ID = types.StringValue(createOutp.Uid)
	planData.SdcType = types.StringValue(createOutp.LarType)
	planData.SdcName = types.StringValue(planData.SdcName.ValueString())
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

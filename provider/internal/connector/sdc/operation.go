package sdc

import (
	"context"

	"github.com/cisco-lockhart/go-client/connector/sdc"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Read(ctx context.Context, resource *SdcResource, stateData *SdcResourceModel) error {

	// do read
	readSdcOutp, err := resource.client.ReadSdcByName(ctx, *sdc.NewReadByNameInput(stateData.Name.ValueString()))
	if err != nil {
		return err
	}

	// map return struct to sdc model
	stateData.ID = types.StringValue(readSdcOutp.Uid)
	stateData.Name = types.StringValue(readSdcOutp.Name)

	return nil
}

func Create(ctx context.Context, resource *SdcResource, planData *SdcResourceModel) error {

	// do create
	createSdcOutp, err := resource.client.CreateSdc(ctx, *sdc.NewCreateInput(planData.Name.ValueString()))
	if err != nil {
		return err
	}

	// map return struct to sdc model
	planData.ID = types.StringValue(createSdcOutp.Uid)
	planData.Name = types.StringValue(createSdcOutp.Name)
	planData.BootstrapData = types.StringValue(createSdcOutp.BootstrapData)

	return nil
}

func Update(ctx context.Context, resource *SdcResource, planData *SdcResourceModel, stateData *SdcResourceModel) error {

	// do update
	updateSdcOutp, err := resource.client.UpdateSdc(ctx, sdc.NewUpdateInput(planData.ID.ValueString(), planData.Name.ValueString()))
	if err != nil {
		return err
	}

	// map return struct to sdc model
	planData.ID = types.StringValue(updateSdcOutp.Uid)
	planData.Name = types.StringValue(updateSdcOutp.Name)
	planData.BootstrapData = types.StringValue(updateSdcOutp.BootstrapData) // bootstrap data contains sdc name, it is not fixed

	return nil
}

func Delete(ctx context.Context, resource *SdcResource, stateData *SdcResourceModel) error {

	// do delete
	_, err := resource.client.DeleteSdc(ctx, sdc.NewDeleteInput(stateData.ID.ValueString()))
	if err != nil {
		return err
	}

	return nil
}

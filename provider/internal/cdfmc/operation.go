package cdfmc

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func Read(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// do read
	// e.g. readOutp, err := resource.client.ReadExample(ctx, ...)

	readOut, err := resource.client.ReadCloudFmcDevice(ctx)
	if err != nil {
		return err
	}

	// map response to terraform types
	// e.g. stateData.ID = types.StringValue(readOutp.Uid)
	stateData.Id = types.StringValue(readOut.Uid)
	stateData.Name = types.StringValue(readOut.Name)

	return nil
}

func Create(ctx context.Context, resource *Resource, planData *ResourceModel) error {

	// do create
	// e.g. createOutp, err := resource.client.CreateExample(ctx, ...)

	createOut, err := resource.client.CreateCloudFmcDevice(ctx, cloudfmc.NewCreateInput())
	if err != nil {
		return err
	}

	// map response to terraform types
	// e.g. planData.ID = types.StringValue(createOutp.Uid)
	planData.Id = types.StringValue(createOut.Uid)
	planData.Name = types.StringValue(createOut.Name)

	return nil
}

func Update(ctx context.Context, resource *Resource, planData *ResourceModel, stateData *ResourceModel) error {

	// empty because nothing to update

	return nil
}

func Delete(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	tflog.Warn(ctx, "Delete cdFMC is a noop, please reach out to CDO TAC if you really want to delete a cdFMC.")

	return nil
}

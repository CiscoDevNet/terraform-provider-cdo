package sec

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sec"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Read(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// do read
	readOutp, err := resource.client.ReadSec(ctx, sec.NewReadInputBuilder().Uid(stateData.Id.ValueString()).Build())
	if err != nil {
		return err
	}

	// map response to terraform types
	stateData.Id = types.StringValue(readOutp.Uid)
	stateData.SecBootstrapData = types.StringValue(readOutp.BootStrapData)
	stateData.Name = types.StringValue(readOutp.Name)

	return nil
}

func Create(ctx context.Context, resource *Resource, planData *ResourceModel) error {

	// do create
	createOutp, err := resource.client.CreateSec(ctx, sec.NewCreateInputBuilder().Build())
	if err != nil {
		return err
	}

	// map response to terraform types
	planData.Id = types.StringValue(createOutp.Uid)
	planData.SecBootstrapData = types.StringValue(createOutp.SecBootstrapData)
	planData.CdoBootstrapData = types.StringValue(createOutp.CdoBoostrapData)
	planData.Name = types.StringValue(createOutp.Name)

	return nil
}

func Update(ctx context.Context, resource *Resource, planData *ResourceModel, stateData *ResourceModel) error {

	// do update
	_, err := resource.client.UpdateSec(ctx, sec.NewUpdateInputBuilder().Build())
	if err != nil {
		return err
	}

	return nil
}

func Delete(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// do delete
	_, err := resource.client.DeleteSec(ctx, sec.NewDeleteInputBuilder().Uid(stateData.Id.ValueString()).Build())
	if err != nil {
		return err
	}

	return nil
}

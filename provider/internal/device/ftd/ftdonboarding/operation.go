package ftdonboarding

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd/cloudftdonboarding"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Read(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// do read
	_, err := resource.client.ReadFtdOnboarding(ctx, cloudftdonboarding.NewReadInput())
	if err != nil {
		return err
	}

	return nil
}

func Create(ctx context.Context, resource *Resource, planData *ResourceModel) error {

	// do create
	createOutp, err := resource.client.CreateFtdOnboarding(ctx, cloudftdonboarding.NewCreateInput(planData.FtdId.ValueString()))
	if err != nil {
		return err
	}

	planData.ID = types.StringValue(createOutp.Name)
	
	return nil
}

func Update(ctx context.Context, resource *Resource, planData *ResourceModel, stateData *ResourceModel) error {

	// do update
	_, err := resource.client.UpdateFtdOnboarding(ctx, cloudftdonboarding.NewUpdateInput())
	if err != nil {
		return err
	}

	return nil
}

func Delete(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// do delete
	_, err := resource.client.DeleteFtdOnboarding(ctx, cloudftdonboarding.NewDeleteInput())
	if err != nil {
		return err
	}

	return nil
}

package ftdonboarding

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd/cloudftdonboarding"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Read(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	_, err := resource.client.ReadFtdOnboarding(ctx, cloudftdonboarding.NewReadInput())
	if err != nil {
		return err
	}

	return nil
}

func Create(ctx context.Context, resource *Resource, planData *ResourceModel) error {

	createOutp, err := resource.client.CreateFtdOnboarding(ctx, cloudftdonboarding.NewCreateInput(planData.FtdUid.ValueString()))
	if err != nil {
		return err
	}

	planData.Id = types.StringValue(createOutp.Metadata.RegKey)

	return nil
}

func Update(ctx context.Context, resource *Resource, planData *ResourceModel, stateData *ResourceModel) error {

	_, err := resource.client.UpdateFtdOnboarding(ctx, cloudftdonboarding.NewUpdateInput())
	if err != nil {
		return err
	}

	return nil
}

func Delete(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	_, err := resource.client.DeleteFtdOnboarding(ctx, cloudftdonboarding.NewDeleteInput())
	if err != nil {
		return err
	}

	return nil
}

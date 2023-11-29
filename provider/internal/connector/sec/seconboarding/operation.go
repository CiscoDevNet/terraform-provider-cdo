package seconboarding

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sec/seconboarding"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Read(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// intentional empty: nothing to read

	return nil
}

func Create(ctx context.Context, resource *Resource, planData *ResourceModel) error {

	// do create
	createOutp, err := resource.client.CreateSecOnboarding(ctx, seconboarding.NewCreateInputBuilder().Name(planData.Name.ValueString()).Build())
	if err != nil {
		return err
	}

	// map response to terraform types
	planData.Id = types.StringValue(createOutp.Uid)
	planData.Name = types.StringValue(createOutp.Name)

	return nil
}

func Update(ctx context.Context, resource *Resource, planData *ResourceModel, stateData *ResourceModel) error {

	// intentional empty: nothing to update

	// map response to terraform types
	stateData.Name = planData.Name

	return nil
}

func Delete(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// intentional empty: nothing to delete

	return nil
}

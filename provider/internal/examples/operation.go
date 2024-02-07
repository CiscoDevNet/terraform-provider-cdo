package examples

import (
	"context"
)

func ReadDataSource(ctx context.Context, resource *ExampleDataSource, stateData *ExampleDataSourceModel) error {

	// do read
	// e.g. readOutp, err := resource.client.ReadExample(ctx, ...)

	// map response to terraform types
	// e.g. stateData.ID = types.StringValue(readOutp.Uid)

	return nil
}

func Read(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// do read
	// e.g. readOutp, err := resource.client.ReadExample(ctx, ...)

	// map response to terraform types
	// e.g. stateData.ID = types.StringValue(readOutp.Uid)

	return nil
}

func Create(ctx context.Context, resource *Resource, planData *ResourceModel) error {

	// do create
	// e.g. createOutp, err := resource.client.CreateExample(ctx, ...)

	// map response to terraform types
	// e.g. planData.ID = types.StringValue(createOutp.Uid)

	return nil
}

func Update(ctx context.Context, resource *Resource, planData *ResourceModel, stateData *ResourceModel) error {

	// do update
	// e.g. updateOutp, err := resource.client.UpdateExample(ctx, ...)

	// map response to terraform types
	// stateData.ID = types.StringValue(updateOutp.Uid)

	return nil
}

func Delete(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// do delete
	// _, err := resource.client.DeleteExample(ctx, ...)

	return nil
}

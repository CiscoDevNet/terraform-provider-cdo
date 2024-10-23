package msp_tenant

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/msp/tenants"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Create(ctx context.Context, resource *TenantResource, planData *TenantResourceModel) error {
	createOut, err := resource.client.CreateTenantUsingMspPortal(ctx, tenants.MspCreateTenantInput{
		Name:        planData.Name.ValueString(),
		DisplayName: planData.DisplayName.ValueString(),
	})
	if err != nil {
		return err
	}

	planData.Id = types.StringValue(createOut.Uid)
	planData.Name = types.StringValue(createOut.Name)
	planData.DisplayName = types.StringValue(createOut.DisplayName)

	return nil
}

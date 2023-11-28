package ftd

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ConfigValidator = &performanceTierConfigValidator{}

type performanceTierConfigValidator struct{}

func (c performanceTierConfigValidator) Description(ctx context.Context) string {
	return c.MarkdownDescription(ctx)
}

func (c performanceTierConfigValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintln("Ensure performance tier is set for virtual FTD.")
}

func (c performanceTierConfigValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var configData ResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &configData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if configData virtual is set to true
	if !configData.Virtual.IsUnknown() && !configData.Virtual.IsNull() && configData.Virtual.ValueBool() {
		// and performance tier is not set
		if configData.PerformanceTier.IsNull() || configData.PerformanceTier.IsUnknown() {
			// then we error
			resp.Diagnostics.AddError("Performance Tier is required for virtual FTD.", "You need to select a performance tiers for virtual FTD. Allowed values are: [\"FTDv5\", \"FTDv10\", \"FTDv20\", \"FTDv30\", \"FTDv50\", \"FTDv100\", \"FTDv\"].")
		}
	}
}

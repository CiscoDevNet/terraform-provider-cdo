package validators

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = mspManagedTenantNameValidator{}

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9-_]{1,50}$`)

type mspManagedTenantNameValidator struct {
}

func (v mspManagedTenantNameValidator) Description(ctx context.Context) string {
	return "Ensures that if name is null and api_token is null, fail. If name is not null and api_token is not null, fail. If name is not null and does not match the regex [a-zA-Z0-9-_]{1,50}, fail."
}

func (v mspManagedTenantNameValidator) MarkdownDescription(ctx context.Context) string {
	return "Ensures that if name is null and api_token is null, fail. If name is not null and api_token is not null, fail. If name is not null and does not match the regex [a-zA-Z0-9-_]{1,50}, fail."
}

func (v mspManagedTenantNameValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	var apiTokenAttr attr.Value

	request.Config.GetAttribute(ctx, path.Root("api_token"), &apiTokenAttr)

	if request.ConfigValue.IsNull() && apiTokenAttr.IsNull() {
		response.Diagnostics.AddError(
			"Invalid Configuration",
			"Both name and api_token cannot be null.",
		)
		return
	}

	if !request.ConfigValue.IsNull() && !apiTokenAttr.IsNull() {
		response.Diagnostics.AddError(
			"Invalid Configuration",
			"Both name and api_token cannot be specified at the same time.",
		)
		return
	}

	if !request.ConfigValue.IsNull() && !nameRegex.MatchString(request.ConfigValue.ValueString()) {
		response.Diagnostics.AddError(
			"Invalid Configuration",
			"Name must match the regex `[a-zA-Z0-9-_]{1,50}`.",
		)
	}
}

func NewMspManagedTenantNameValidator() validator.String {
	return mspManagedTenantNameValidator{}
}

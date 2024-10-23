package validators

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

type cdoTenantValidator struct{}

var ValidationString = "The tenant name can only contain alphabets, numbers, -, and _, and is limited to 50 characters."
var TenantNotEmptyString = "The tenant name cannot be empty."

func (cdoTenantValidator) Description(context.Context) string {
	return ValidationString
}

func (cdoTenantValidator) MarkdownDescription(context.Context) string {
	return ValidationString
}

func (c cdoTenantValidator) ValidateString(ctx context.Context, req validator.StringRequest, res *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		res.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			req.Path,
			TenantNotEmptyString,
			"",
		))

		return
	}

	strValue := req.ConfigValue.ValueString()
	re := regexp.MustCompile("^[a-zA-Z0-9-_]{1,50}$")
	if !re.MatchString(strValue) {
		res.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			req.Path,
			ValidationString,
			strValue,
		))
	}
}

func NewCdoTenantValidator() validator.String {
	return cdoTenantValidator{}
}

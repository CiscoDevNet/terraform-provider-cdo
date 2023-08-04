package validators

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = socketAddressValidator{}

var socketAddressRegex = regexp.MustCompile(`.+:\d+`)

// socketAddressValidator validates that the socket address has a host and a port.
type socketAddressValidator struct {
}

func (v socketAddressValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v socketAddressValidator) MarkdownDescription(_ context.Context) string {
	return "value must contain a host and an integer port, in the format `host:port`"
}

func (v socketAddressValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	if !socketAddressRegex.Match([]byte(value.ValueString())) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value.String(),
		))
	}
}

// OneOf checks that the String held in the attribute
// is one of the given `values`.
func ValidateSocketAddress() validator.String {
	return socketAddressValidator{}
}

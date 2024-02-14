package validators

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/settings"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type conflictDetctionIntervalValidator struct{}

func (conflictDetctionIntervalValidator) Description(context.Context) string {
	return "The conflict detection value must be on of: EVERY_10_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_24_HOURS"
}

func (conflictDetctionIntervalValidator) MarkdownDescription(context.Context) string {
	return "The conflict detection value must be on of: `EVERY_10_MINUTES`, `EVERY_HOUR`, `EVERY_6_HOURS`, `EVERY_24_HOURS`"
}

func (v conflictDetctionIntervalValidator) ValidateString(ctx context.Context, req validator.StringRequest, res *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	strValue := req.ConfigValue.ValueString()
	interval := settings.ResolveConflictDetectionInterval(strValue)

	if interval == settings.ConflictDetectionIntervalInvalid {
		res.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			req.Path,
			v.Description(ctx),
			strValue,
		))
	}
}

func NewConflictDetectionIntervalValidator() validator.String {
	return conflictDetctionIntervalValidator{}
}

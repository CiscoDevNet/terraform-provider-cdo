// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validators

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = oneOfValidator{}

// oneOfValidator validates that the value matches one of expected values.
type oneOfValidator struct {
	values []types.String
}

func (v oneOfValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v oneOfValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be one of: %q", v.values)
}

func (v oneOfValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	for _, validValue := range v.values {
		if value.Equal(validValue) {
			return
		}
	}

	response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
		request.Path,
		v.Description(ctx),
		value.String(),
	))
}

// OneOf checks that the String held in the attribute
// is one of the given `values`.
func OneOf(values ...string) validator.String {
	frameworkValues := make([]types.String, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.StringValue(value))
	}

	return oneOfValidator{
		values: frameworkValues,
	}
}

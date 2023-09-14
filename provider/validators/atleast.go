package validators

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ValueStringsAtLeast returns a validator which ensures that at least one configured
// String values passes each String validator.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func ValueStringsAtLeast(elementValidators ...validator.String) validator.List {
	return valueStringsAtLeastValidator{
		elementValidators: elementValidators,
	}
}

var _ validator.List = valueStringsAtLeastValidator{}

// valueStringsAtLeastValidator validates that each List member validates against each of the value validators.
type valueStringsAtLeastValidator struct {
	elementValidators []validator.String
}

// Description describes the validation in plain text formatting.
func (v valueStringsAtLeastValidator) Description(ctx context.Context) string {
	var descriptions []string

	for _, elementValidator := range v.elementValidators {
		descriptions = append(descriptions, elementValidator.Description(ctx))
	}

	return fmt.Sprintf("at least one element value must satisfy all validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v valueStringsAtLeastValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateList performs the validation.
func (v valueStringsAtLeastValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	_, ok := req.ConfigValue.ElementType(ctx).(basetypes.StringTypable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Validator for Element Type",
			"While performing schema-based validation, an unexpected error occurred. "+
				"The attribute declares a String values validator, however its values do not implement types.StringType or the types.StringTypable interface for custom String types. "+
				"Use the appropriate values validator that matches the element type. "+
				"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
				fmt.Sprintf("Path: %s\n", req.Path.String())+
				fmt.Sprintf("Element Type: %T\n", req.ConfigValue.ElementType(ctx)),
		)

		return
	}

	elements := req.ConfigValue.Elements()
	allDiagnostics := diag.Diagnostics{}

	for idx, element := range elements {
		elementPath := req.Path.AtListIndex(idx)

		elementValuable, ok := element.(basetypes.StringValuable)

		// The check above should have prevented this, but raise an error
		// instead of a type assertion panic or skipping the element. Any issue
		// here likely indicates something wrong in the framework itself.
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Validator for Element Value",
				"While performing schema-based validation, an unexpected error occurred. "+
					"The attribute declares a String values validator, however its values do not implement types.StringType or the types.StringTypable interface for custom String types. "+
					"This is likely an issue with terraform-plugin-framework and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Path: %s\n", req.Path.String())+
					fmt.Sprintf("Element Type: %T\n", req.ConfigValue.ElementType(ctx))+
					fmt.Sprintf("Element Value Type: %T\n", element),
			)

			return
		}

		elementValue, diags := elementValuable.ToStringValue(ctx)

		resp.Diagnostics.Append(diags...)

		// Only return early if the new diagnostics indicate an issue since
		// it likely will be the same for all elements.
		if diags.HasError() {
			return
		}

		elementReq := validator.StringRequest{
			Path:           elementPath,
			PathExpression: elementPath.Expression(),
			ConfigValue:    elementValue,
			Config:         req.Config,
		}

		elementDiagnostics := diag.Diagnostics{}

		for _, elementValidator := range v.elementValidators {
			elementResp := &validator.StringResponse{}

			elementValidator.ValidateString(ctx, elementReq, elementResp)

			elementDiagnostics.Append(elementResp.Diagnostics...)
		}

		// this element passed all validators, return directly, dropping any diagnostics from previous elements
		// note this skips warning
		if !elementDiagnostics.HasError() {
			return
		}

		// this element did not pass all validators, append its diagnostics to allDiagnostics
		allDiagnostics.Append(elementDiagnostics...)
	}

	// no element pass all validators, returning all diagnostics
	resp.Diagnostics.Append(allDiagnostics...)
}

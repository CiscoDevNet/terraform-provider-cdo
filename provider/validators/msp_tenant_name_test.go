package validators_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/validators"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"testing"
)

func TestMspTenantNameValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        types.String
		apiToken    attr.Value
		expectError bool
	}

	testCases := map[string]testCase{
		"non-null-name-and-non-null-api-token": {
			name:        types.StringValue("burak-crush-pineapple"),
			apiToken:    types.StringValue("burak-crush-api-token"),
			expectError: true,
		},
		"non-null-name-and-null-api-token": {
			name:        types.StringValue("burak-crush-pineapple"),
			apiToken:    nil,
			expectError: false,
		},
		"null-name-and-null-api-token": {
			name:        types.StringNull(),
			apiToken:    nil,
			expectError: true,
		},
		"null-name-and-non-null-api-token": {
			name:        types.StringNull(),
			apiToken:    types.StringValue("burak-crush-api-token"),
			expectError: false,
		},
	}

	for name, test := range testCases {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{ // nolint
				ConfigValue: test.name,
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"api_token": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"api_token": tftypes.NewValue(tftypes.String, test.apiToken.(types.String).ValueString()), // nolint
					}),
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"api_token": schema.StringAttribute{
								MarkdownDescription: "API token for an API-only user with super-admin privileges on the tenant",
								Required:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(), // Prevent updates to name
								},
							},
						},
					}, // Provide the actual schema if needed
				},
			}
			res := validator.StringResponse{}
			validators.NewMspManagedTenantNameValidator().ValidateString(context.TODO(), req, &res)
			if test.expectError && !res.Diagnostics.HasError() {
				t.Fatalf("expected error, got none")

			}
		})
	}
}

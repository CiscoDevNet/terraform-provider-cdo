package validators_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/validators"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"testing"
)

func TestValidTenantNamesValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        types.String
		expErrors int
	}

	testCases := map[string]testCase{
		"emptyTenantName": {
			in:        types.StringValue(""),
			expErrors: 1,
		},
		"nullTenantName": {
			in:        types.StringNull(),
			expErrors: 1,
		},
		"invalidTenantName": {
			in:        types.StringValue("burak!!!!$%"),
			expErrors: 1,
		},
		"validTenantName": {
			in:        types.StringValue("burak-crush-pineapple"),
			expErrors: 0,
		},
	}

	for name, test := range testCases {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{
				ConfigValue: test.in,
			}
			res := validator.StringResponse{}
			validators.NewCdoTenantValidator().ValidateString(context.TODO(), req, &res)

			if test.expErrors > 0 && !res.Diagnostics.HasError() {
				t.Fatalf("expected %d error(s), got none", test.expErrors)
			}

			if test.expErrors > 0 && test.expErrors != res.Diagnostics.ErrorsCount() {
				t.Fatalf("expected %d error(s), got %d: %v", test.expErrors, res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}

			if test.expErrors == 0 && res.Diagnostics.HasError() {
				t.Fatalf("expected no error(s), got %d: %v", res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}
		})
	}
}

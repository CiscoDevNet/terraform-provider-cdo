package validators_test

import (
	"context"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/validators"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestOneOfRolesValidator(t *testing.T) {
	t.Parallel()

	SuperAdminRoleToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlcyI6WyJST0xFX1NVUEVSX0FETUlOIl19.FBEPzCpXPiYX6esud26WMrJvU3reLDVFdLociGrilVw"
	AdminRoleToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlcyI6WyJST0xFX0FETUlOIl19.Xlk9JQJV9butj8ERu7mZFvRczY66KvmnTklQVnHwoy0"
	BlaRoleToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlcyI6WyJST0xFX0JMQSJdfQ.nx10m1Cb8ZkUbxWKyfyJzsD1Y8rdTtkRbGEC2A4yV2M"
	NoRolesToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlcyI6W119.0Dd6yUeJ4UbCr8WyXOiK3BhqVVwJFk5c53ipJBWenmc"

	type testCase struct {
		in        types.String
		validator validator.String
		expErrors int
	}

	testCases := map[string]testCase{
		"super-admin-role-match": {
			in: types.StringValue(SuperAdminRoleToken),
			validator: validators.OneOfRoles(
				"ROLE_SUPER_ADMIN", "ROLE_ADMIN",
			),
			expErrors: 0,
		},
		"admin-role-match": {
			in: types.StringValue(AdminRoleToken),
			validator: validators.OneOfRoles(
				"ROLE_SUPER_ADMIN", "ROLE_ADMIN",
			),
			expErrors: 0,
		},
		"simple-mismatch": {
			in: types.StringValue(BlaRoleToken),
			validator: validators.OneOfRoles(
				"ROLE_SUPER_ADMIN", "ROLE_ADMIN",
			),
			expErrors: 1,
		},
		"no-roles": {
			in: types.StringValue(NoRolesToken),
			validator: validators.OneOfRoles(
				"ROLE_SUPER_ADMIN", "ROLE_ADMIN",
			),
			expErrors: 1,
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
			test.validator.ValidateString(context.TODO(), req, &res)

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

package validators

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	ErrInvalidTokenFormat = "invalid JWT token format"
	ErrDecodeFailed       = "failed to decode token payload"
	ErrNoRolesFound       = "no roles found in the JWT token"
)

var _ validator.String = oneOfRolesValidator{}

// oneOfRolesValidator validates that the value matches one of expected roles.
type oneOfRolesValidator struct {
	expectedRoles []types.String
}

func (v oneOfRolesValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v oneOfRolesValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("The user must be assigned one of the following CDO roles: %q", v.expectedRoles)
}

func (v oneOfRolesValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	token := request.ConfigValue.String()

	role, err := extractRoleFromToken(token)

	if err != nil {
		fmt.Println("Error:", err)
	}

	for _, expectedRole := range v.expectedRoles {
		if role == expectedRole.ValueString() {
			return
		}
	}

	response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
		request.Path,
		v.Description(ctx),
		role,
	))
}

func extractRoleFromToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", fmt.Errorf("tokenString is nil or empty")
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("%s", ErrInvalidTokenFormat)
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("%s", ErrDecodeFailed)
	}

	var claims jwt.MapClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return "", fmt.Errorf("failed to decode token payload: %v", err)
	}

	if rolesClaim, exists := claims["roles"]; exists { // check if "roles" claim exists
		if roles, ok := rolesClaim.([]interface{}); ok { // convert to a slice of interfaces
			if len(roles) > 0 {
				if role, ok := roles[0].(string); ok {
					return role, nil
				}
			}
		}
	}

	return "", fmt.Errorf("%s", ErrNoRolesFound)
}

// OneOfRoles checks that the JWT token roles String held in the attribute
// is one of the given `roles`.
func OneOfRoles(roles ...string) validator.String {
	frameworkValues := make([]types.String, 0, len(roles))

	for _, value := range roles {
		frameworkValues = append(frameworkValues, types.StringValue(value))
	}

	return oneOfRolesValidator{
		expectedRoles: frameworkValues,
	}
}

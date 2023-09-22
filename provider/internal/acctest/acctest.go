// ci user: terraform-provider-cdo@lockhart.io
package acctest

import (
	"fmt"
	"os"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	apiTokenEnvName    = "ACC_TEST_CISCO_CDO_API_TOKEN"
	apiTokenSecretName = "staging-terraform-provider-cdo-acceptance-test-api-token"
)

var cdoSecretManager = NewCdoSecretManager("us-west-2")

func GetApiToken() (string, error) {
	tokenFromEnv, ok := os.LookupEnv(apiTokenEnvName)
	if ok {
		return tokenFromEnv, nil
	}

	tokenFromSecretManager, err := cdoSecretManager.getCurrentSecretValue(apiTokenSecretName)
	if err == nil {
		return tokenFromSecretManager, nil
	}

	return "", fmt.Errorf("failed to retrieve api token from environment variable and secret manager.\nenvironment variable name=%s\nsecret manager secret token name=%s\nplease set one of them.\ncause=%v", apiTokenEnvName, apiTokenSecretName, err)
}

func PreCheckFunc(t *testing.T) func() {
	return func() {
		_, err := GetApiToken()
		if err != nil {
			t.Fatalf("Precheck failed, cause=%v", err)
		}
	}
}

func ProviderConfig() string {
	token, err := GetApiToken()
	if err != nil {
		panic(fmt.Errorf("failed to retrieve api token, cause=%w", err))
	}
	// TODO: delete
	token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ2ZXIiOiIwIiwic2NvcGUiOlsidHJ1c3QiLCJyZWFkIiwiYzZjZTA1YjEtMmZkMi00ZjJlLTkwZDYtZjgyNzUxNjhmODhkIiwid3JpdGUiXSwicm9sZXMiOlsiUk9MRV9TVVBFUl9BRE1JTiJdLCJhbXIiOm51bGwsImlzcyI6Iml0ZCIsImNsdXN0ZXJJZCI6IjEiLCJpZCI6IjFlMDg0ZDc5LTg2NjgtNDVlYi1hYTA1LTI3Y2M2NzE2N2FkMCIsInN1YmplY3RUeXBlIjoidXNlciIsImp0aSI6IjIxY2YwYTUwLTNlMTgtNGNhZC04ZGIxLWVmZjE0NDc0YzI2NSIsInBhcmVudElkIjoiYzZjZTA1YjEtMmZkMi00ZjJlLTkwZDYtZjgyNzUxNjhmODhkIiwiY2xpZW50X2lkIjoiYXBpLWNsaWVudCJ9.iDGvmfmVioAZP4c5rFX02FdXJEp3Hys77jW82bT5lrtlO9Ev0i53O9IT0ztCegAgFGLTfw53TMP964LHqe0A1Bhdlsm3V5zLXn6BIn61HogcNw9ZOweuO7ZkGsCalcQ6KmhG_uXiX56ML5_d-XLnhPJZIq3oI7Qbd6fgXcl_x3ru32k-q8FW3CFc03dQ7hxdLjXFVZph6DxoSDRkeVTqi10oFtb6VRoFktlnfLV4ks4rDlG5eRcRe1yg1-TTbE7KREsZl645ayJEm14s_HI573h0Ub7rDopDCKWOHVztgDOAPeS82EMP8_IKTyrecbW5tA9ZOs1oEh00WADX6oiBYg"

	return fmt.Sprintf(`
	provider "cdo" {
		api_token = "%s"
		base_url = "https://ci.dev.lockhart.io"
	}
	// New line
	`, token)
}

// ProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"cdo": providerserver.NewProtocol6WithError(provider.New("test")()),
}

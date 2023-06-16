// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const acceptanceTestApiTokenEnvVarName = "ACC_TEST_CISCO_CDO_API_TOKEN"

var acceptanceTestApiToken = os.Getenv(acceptanceTestApiTokenEnvVarName)

var providerConfig = fmt.Sprintf(`
provider "cdo" {
	api_token = "%s"
}
// New line
`, acceptanceTestApiToken)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"cdo": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if acceptanceTestApiToken == "" {
		t.Fatalf("Environment variable: '%s' must be set", acceptanceTestApiTokenEnvVarName)
	}
}

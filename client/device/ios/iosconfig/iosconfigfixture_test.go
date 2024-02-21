package iosconfig_test

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/ios/iosconfig"
	"net/http"

	"github.com/jarcoal/httpmock"
)

const (
	iosConfigUid = "00000000-0000-0000-0000-000000000000"
	baseUrl      = "https://unittest.cdo.cisco.com"
)

func buildIosConfigPath(specificDeviceUid string) string {
	return fmt.Sprintf("/aegis/rest/v1/services/targets/devices/%s", specificDeviceUid)
}

func configureIosConfigReadToSucceedInSubsequentCalls(specificDeviceUid string, outputs []iosconfig.ReadOutput) {
	callCount := 0
	httpmock.RegisterResponder("GET", buildIosConfigPath(specificDeviceUid), func(r *http.Request) (*http.Response, error) {
		defer func() {
			callCount += 1
		}()

		if callCount >= len(outputs) {
			panic("no more configured iosconfig read calls")
		}

		return httpmock.NewJsonResponse(200, outputs[callCount])
	})
}

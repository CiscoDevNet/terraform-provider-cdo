package cloudftd_test

import "github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"

const (
	baseUrl                 = "https://unit-test.net"
	deviceName              = "unit-test-device-name"
	deviceUid               = "unit-test-uid"
	deviceHost              = "https://unit-test.com"
	devicePort              = 1234
	deviceCloudConnectorUId = "unit-test-uid"
)

var (
	validReadFmcOutput = device.NewReadOutputBuilder().
		AsCloudFmc().
		WithName(deviceName).
		WithUid(deviceUid).
		WithLocation(deviceHost, devicePort).
		OnboardedUsingCloudConnector(deviceCloudConnectorUId).
		Build()
)

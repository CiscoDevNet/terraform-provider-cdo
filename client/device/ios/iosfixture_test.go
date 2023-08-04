package ios

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/connector/sdc"
	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/device"
	internalTesting "github.com/CiscoDevNet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/jarcoal/httpmock"
)

const (
	deviceCreatePath = "/aegis/rest/v1/services/targets/devices"
)

func buildDeviceReadSpecificPath(deviceUid string) string {
	return fmt.Sprintf("/aegis/rest/v1/device/%s/specific-device", deviceUid)
}

func buildDevicePath(deviceUid string) string {
	return fmt.Sprintf("/aegis/rest/v1/services/targets/devices/%s", deviceUid)
}

func buildSdcPath(sdcUid string) string {
	return fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", sdcUid)
}

func configureDeviceCreateToRespondSuccessfully(createOutput device.CreateOutput) {
	httpmock.RegisterResponder(
		http.MethodPost,
		deviceCreatePath,
		httpmock.NewJsonResponderOrPanic(200, createOutput),
	)
}

func configureDeviceCreateToRespondWithError() {
	httpmock.RegisterResponder(
		http.MethodPost,
		deviceCreatePath,
		httpmock.NewStringResponder(500, ""),
	)
}

func configureDeviceUpdateToRespondSuccessfully(updateOutput device.UpdateOutput) {
	httpmock.RegisterResponder(
		http.MethodPut,
		buildDevicePath(updateOutput.Uid),
		httpmock.NewJsonResponderOrPanic(200, updateOutput),
	)
}

func configureDeviceUpdateToRespondWithError(deviceUid string) {
	httpmock.RegisterResponder(
		http.MethodPut,
		buildDevicePath(deviceUid),
		httpmock.NewStringResponder(500, ""),
	)
}

func configureDeviceReadToRespondSuccessfully(readOutput device.ReadOutput) {
	httpmock.RegisterResponder(
		http.MethodGet,
		buildDevicePath(readOutput.Uid),
		httpmock.NewJsonResponderOrPanic(200, readOutput),
	)
}

func configureDeviceReadToRespondWithError(deviceUid string) {
	httpmock.RegisterResponder(
		http.MethodGet,
		buildDevicePath(deviceUid),
		httpmock.NewStringResponder(500, ""),
	)
}

func configureDeviceReadSpecificToRespondSuccessfully(deviceUid string, readOutput device.ReadSpecificOutput) {
	httpmock.RegisterResponder(
		http.MethodGet,
		buildDeviceReadSpecificPath(deviceUid),
		httpmock.NewJsonResponderOrPanic(200, readOutput),
	)
}

func configureDeviceReadSpecificToRespondWithError(deviceUid string) {
	httpmock.RegisterResponder(
		http.MethodGet,
		buildDeviceReadSpecificPath(deviceUid),
		httpmock.NewStringResponder(500, ""))
}

func configureDeviceDeleteToRespondSuccessfully(deviceUid string) {
	httpmock.RegisterResponder(
		http.MethodDelete,
		buildDevicePath(deviceUid),
		httpmock.NewStringResponder(200, ""),
	)
}

func configureDeviceDeleteToRespondWithError(deviceUid string) {
	httpmock.RegisterResponder(
		http.MethodDelete,
		buildDevicePath(deviceUid),
		httpmock.NewStringResponder(500, ""),
	)
}

func configureIosConfigReadToSucceedWithSubsequentCalls(specificUid string, responders []httpmock.Responder) {
	count := 0

	httpmock.RegisterResponder(
		http.MethodGet,
		buildDevicePath(specificUid),
		func(r *http.Request) (*http.Response, error) {
			responder := responders[count]
			count += 1

			return responder(r)
		},
	)
}

func configureSdcReadToRespondSuccessfully(readOutput sdc.ReadOutput) {
	httpmock.RegisterResponder(
		http.MethodGet,
		buildSdcPath(readOutput.Uid),
		httpmock.NewJsonResponderOrPanic(200, readOutput),
	)
}

func configureSdcReadToRespondWithError(sdcUid string) {
	httpmock.RegisterResponder(
		http.MethodGet,
		buildSdcPath(sdcUid),
		httpmock.NewStringResponder(500, ""),
	)
}

func assertDeviceCreateWasCalledOnce(t *testing.T) {
	internalTesting.AssertEndpointCalledTimes(http.MethodPost, deviceCreatePath, 1, t)
}

func assertDeviceReadWasCalledTimes(deviceUid string, times int, t *testing.T) {
	internalTesting.AssertEndpointCalledTimes(http.MethodGet, buildDevicePath(deviceUid), times, t)
}

func assertDeviceUpdateWasCalledOnce(deviceUid string, t *testing.T) {
	internalTesting.AssertEndpointCalledTimes(http.MethodPut, buildDevicePath(deviceUid), 1, t)
}

func assertSdcReadByUidWasCalledOnce(sdcUid string, t *testing.T) {
	internalTesting.AssertEndpointCalledTimes(http.MethodGet, buildSdcPath(sdcUid), 1, t)
}

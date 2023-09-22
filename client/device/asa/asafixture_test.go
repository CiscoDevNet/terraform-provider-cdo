package asa_test

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa/asaconfig"
	internalhttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	internalTesting "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/jarcoal/httpmock"
)

const (
	deviceCreatePath = "/aegis/rest/v1/services/targets/devices"
	baseUrl          = "https://unittest.cdo.cisco.com"
)

func buildDeviceReadSpecificPath(deviceUid string) string {
	return fmt.Sprintf("/aegis/rest/v1/device/%s/specific-device", deviceUid)
}

func buildDevicePath(deviceUid string) string {
	return fmt.Sprintf("/aegis/rest/v1/services/targets/devices/%s", deviceUid)
}

func buildAsaConfigPath(specificUid string) string {
	return fmt.Sprintf("/aegis/rest/v1/services/asa/configs/%s", specificUid)
}

func buildConnectorPath(connectorUid string) string {
	return fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", connectorUid)
}

func configureDeviceCreateToRespondSuccessfully(createOutput device.CreateOutput) {
	httpmock.RegisterResponder(
		http.MethodPost,
		deviceCreatePath,
		httpmock.NewJsonResponderOrPanic(200, createOutput),
	)
}

func configureDeviceCreateToRespondSuccessfullyWithNewModel(t *testing.T, createOutput device.CreateOutput) {
	httpmock.RegisterResponder(
		http.MethodPost,
		deviceCreatePath,
		func(req *http.Request) (*http.Response, error) {
			createInp, err := internalhttp.ReadRequestBody[device.CreateInput](req)
			if err != nil {
				return nil, err
			}
			expectedMetadata := &asa.Metadata{IsNewPolicyObjectModel: "true"}
			expectedBytes, err := json.Marshal(expectedMetadata)
			if err != nil {
				return nil, err
			}
			actualBytes, err := json.Marshal(createInp.Metadata)
			if err != nil {
				return nil, err
			}
			expectedMetadataPayload := string(expectedBytes)
			actualMetadataPayload := string(actualBytes)
			assert.Equal(t, expectedMetadataPayload, actualMetadataPayload)
			return httpmock.NewJsonResponse(http.StatusOK, createOutput)
		},
	)
}

func configureDeviceCreateToRespondSuccessfullyWithoutNewModel(t *testing.T, createOutput device.CreateOutput) {
	httpmock.RegisterResponder(
		http.MethodPost,
		deviceCreatePath,
		func(req *http.Request) (*http.Response, error) {
			createInp, err := internalhttp.ReadRequestBody[device.CreateInput](req)
			if err != nil {
				return nil, err
			}
			assert.True(t, reflect.TypeOf(createInp.Metadata).Kind() == reflect.Pointer)
			return httpmock.NewJsonResponse(http.StatusOK, createOutput)
		},
	)
}

func configureDeviceCreateToRespondWithError() {
	httpmock.RegisterResponder(
		http.MethodPost,
		deviceCreatePath,
		httpmock.NewStringResponder(500, ""),
	)
}

func configureDeviceUpdateToRespondSuccessfully(deviceUid string, updateOutput device.UpdateOutput) {
	httpmock.RegisterResponder(
		http.MethodPut,
		buildDevicePath(deviceUid),
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

func configureAsaConfigReadToRespondSuccessfully(specificUid string, readOutput asaconfig.ReadOutput) {
	httpmock.RegisterResponder(
		http.MethodGet,
		buildAsaConfigPath(specificUid),
		httpmock.NewJsonResponderOrPanic(200, readOutput),
	)
}

func configureAsaConfigReadToRespondWithCalls(specificUid string, responders []httpmock.Responder) {
	count := 0

	httpmock.RegisterResponder(
		http.MethodGet,
		buildAsaConfigPath(specificUid),
		func(r *http.Request) (*http.Response, error) {
			responder := responders[count]
			count += 1

			return responder(r)
		},
	)
}

func configureAsaConfigReadToRespondWithError(specificUid string) {
	httpmock.RegisterResponder(
		http.MethodGet,
		buildAsaConfigPath(specificUid),
		httpmock.NewStringResponder(500, ""),
	)
}

func configureAsaConfigUpdateToRespondSuccessfully(specificUid string, updateOutput asaconfig.UpdateOutput) {
	httpmock.RegisterResponder(
		http.MethodPut,
		buildAsaConfigPath(specificUid),
		httpmock.NewJsonResponderOrPanic(200, updateOutput),
	)
}

func configureAsaConfigUpdateToRespondWithError(specificUid string) {
	httpmock.RegisterResponder(
		http.MethodPut,
		buildAsaConfigPath(specificUid),
		httpmock.NewStringResponder(500, ""),
	)
}

func configureConnectorReadToRespondSuccessfully(readOutput connector.ReadOutput) {
	httpmock.RegisterResponder(
		http.MethodGet,
		buildConnectorPath(readOutput.Uid),
		httpmock.NewJsonResponderOrPanic(200, readOutput),
	)
}

func configureConnectorReadToRespondWithError(connectorUid string) {
	httpmock.RegisterResponder(
		http.MethodGet,
		buildConnectorPath(connectorUid),
		httpmock.NewStringResponder(500, ""),
	)
}

func configureReadApiTokenInfoSuccessfully(tokenInfo user.GetTokenInfoOutput) {
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadTokenInfo(baseUrl),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, tokenInfo),
	)
}

func configureReadApiTokenInfoFailed() {
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadTokenInfo(baseUrl),
		httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
	)
}

func assertDeviceCreateWasCalledOnce(t *testing.T) {
	internalTesting.AssertEndpointCalledTimes(http.MethodPost, deviceCreatePath, 1, t)
}

func assertDeviceReadSpecificWasCalledOnce(uid string, t *testing.T) {
	internalTesting.AssertEndpointCalledTimes(http.MethodGet, buildDeviceReadSpecificPath(uid), 1, t)
}

func assertDeviceUpdateWasCalledOnce(deviceUid string, t *testing.T) {
	internalTesting.AssertEndpointCalledTimes(http.MethodPut, buildDevicePath(deviceUid), 1, t)
}

func assertAsaConfigReadWasCalledTimes(specificUid string, times int, t *testing.T) {
	internalTesting.AssertEndpointCalledTimes(http.MethodGet, buildAsaConfigPath(specificUid), times, t)
}

func assertAsaConfigUpdateWasCalledOnce(specificUid string, t *testing.T) {
	internalTesting.AssertEndpointCalledTimes(http.MethodPut, buildAsaConfigPath(specificUid), 1, t)
}

func assertConnectorReadByUidWasCalledOnce(connectorUid string, t *testing.T) {
	internalTesting.AssertEndpointCalledTimes(http.MethodGet, buildConnectorPath(connectorUid), 1, t)
}

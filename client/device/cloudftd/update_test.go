package cloudftd_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcplatform"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/devicelicense"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCloudFtd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testHost := "test-host"
	var testPort uint = 443

	testCloudFtdInput := cloudftd.NewUpdateInputBuilder().
		Uid("test-uid").
		Licenses([]license.Type{license.Essentials}).
		Name("test-name").
		Tags(
			tags.New([]string{"test-tag"}, map[string][]string{
				"grouped-tags": {"grouped-tags-1", "grouped-tags-2"},
			}),
		).
		Build()

	testMetadata := cloudftd.NewMetadataBuilder().
		LicenseCaps(&[]license.Type{}).
		Build()

	successFmcSpecificOutput := cloudfmc.NewReadSpecificOutputBuilder().
		SpecificUid("test-specific-uid").
		State(state.DONE).
		Build()

	successCloudFtdOutput := cloudftd.NewUpdateOutputBuilder().
		Uid(testCloudFtdInput.Uid).
		Metadata(testMetadata).
		Name(testCloudFtdInput.Name).
		Build()

	successCloudFmcOutput := cloudfmc.NewReadOutputBuilder().
		WithUid(successFmcSpecificOutput.SpecificUid).
		WithLocation(testHost, testPort).
		Build()

	successReadDeviceLicenseOutput := fmcplatform.NewReadDeviceLicensesOutputBuilder().
		Id("test-license-id").
		Build()

	successUpdateDeviceLicenseOutput := fmcplatform.NewUpdateDeviceLicensesOutputBuilder().
		Id("test-license-id").
		Build()

	successFmcApplianceOutput := cloudftd.NewUpdateOutputBuilder().
		Uid(successFmcSpecificOutput.SpecificUid).
		Build()

	testCases := []struct {
		testName   string
		input      cloudftd.UpdateInput
		setupFunc  func()
		assertFunc func(output *cloudftd.UpdateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully update Cloud FTD",
			input:    testCloudFtdInput,
			setupFunc: func() {
				updateCdoFtdSettings(baseUrl, successCloudFtdOutput)
				readCloudFmc(baseUrl, successCloudFmcOutput)
				readFtdDeviceLicense(baseUrl, successReadDeviceLicenseOutput)
				updateFtdDeviceLicense(baseUrl, successUpdateDeviceLicenseOutput)
				readCloudFmcSpecific(baseUrl, successFmcSpecificOutput)
				updateCloudFmcAppliance(baseUrl, successFmcApplianceOutput)
				readFtd(baseUrl, successCloudFtdOutput)
			},
			assertFunc: func(output *cloudftd.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, successCloudFtdOutput, *output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := cloudftd.Update(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

func updateCdoFtdSettings(baseUrl string, output cloudftd.UpdateOutput) {
	httpmock.RegisterResponder(
		http.MethodPut,
		url.UpdateDevice(baseUrl, output.Uid),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, output),
	)
}

func readCloudFmc(baseUrl string, output cloudfmc.ReadOutput) {
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadAllDevicesByType(baseUrl),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, []cloudfmc.ReadOutput{output}),
	)
}

func readCloudFmcSpecific(baseUrl string, output cloudfmc.ReadSpecificOutput) {
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadSpecificDevice(baseUrl, output.SpecificUid),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, output),
	)
}

func updateCloudFmcAppliance(baseUrl string, output cloudftd.UpdateOutput) {
	httpmock.RegisterResponder(
		http.MethodPut,
		url.UpdateFmcAppliance(baseUrl, output.Uid),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, output),
	)
}

func readFtdDeviceLicense(baseUrl string, readOutput fmcplatform.ReadDeviceLicensesOutput) {
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadFmcDeviceLicenses(baseUrl),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, devicelicense.DeviceLicense{Items: []fmcplatform.ReadDeviceLicensesOutput{readOutput}}),
	)
}

func updateFtdDeviceLicense(baseUrl string, updateOutput fmcplatform.UpdateDeviceLicensesOutput) {
	httpmock.RegisterResponder(
		http.MethodPut,
		url.UpdateFmcDeviceLicenses(baseUrl, updateOutput.Id),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, updateOutput),
	)
}

func readFtd(baseUrl string, output cloudftd.ReadOutput) {
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadDevice(baseUrl, output.Uid),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, output),
	)
}

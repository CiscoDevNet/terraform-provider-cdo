package cloudftd_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/statemachine"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestDeleteCloudFtd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := []struct {
		testName   string
		input      cloudftd.DeleteInput
		setupFunc  func()
		assertFunc func(output *cloudftd.DeleteOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully delete Cloud FTD, and waited for delete state machine",
			input:    cloudftd.NewDeleteInput(ftdUid),
			setupFunc: func() {
				readFmcIsSuccessful(true)
				readFmcSpecificIsSuccessful(true)
				readFmcDomainInfoIsSuccessful(true)
				readFtdIsSuccessful(true)
				readFmcDeviceRecordsIsSuccessful(true)
				readFtdDeviceRecordIsSuccessful(true)
				triggerFtdDeleteOnFmcIsSuccessful(true)
				waitForFtdDeleteStateMachineEndedSuccessful(true)
				deleteFtdIsSuccessful(true)
			},
			assertFunc: func(output *cloudftd.DeleteOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validDeleteOutput, *output)
			},
		},
		{
			testName: "error when failed to read FMC",
			input:    cloudftd.NewDeleteInput(ftdUid),
			setupFunc: func() {
				readFmcIsSuccessful(false)
				readFmcSpecificIsSuccessful(true)
				readFmcDomainInfoIsSuccessful(true)
				readFtdIsSuccessful(true)
				readFmcDeviceRecordsIsSuccessful(true)
				readFtdDeviceRecordIsSuccessful(true)
				triggerFtdDeleteOnFmcIsSuccessful(true)
				waitForFtdDeleteStateMachineEndedSuccessful(true)
				deleteFtdIsSuccessful(true)
			},
			assertFunc: func(output *cloudftd.DeleteOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "error when failed to read FMC specific",
			input:    cloudftd.NewDeleteInput(ftdUid),
			setupFunc: func() {
				readFmcIsSuccessful(true)
				readFmcSpecificIsSuccessful(false)
				readFmcDomainInfoIsSuccessful(true)
				readFtdIsSuccessful(true)
				readFmcDeviceRecordsIsSuccessful(true)
				readFtdDeviceRecordIsSuccessful(true)
				triggerFtdDeleteOnFmcIsSuccessful(true)
				waitForFtdDeleteStateMachineEndedSuccessful(true)
				deleteFtdIsSuccessful(true)
			},
			assertFunc: func(output *cloudftd.DeleteOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "error when failed to trigger delete FTD state machine",
			input:    cloudftd.NewDeleteInput(ftdUid),
			setupFunc: func() {
				readFmcIsSuccessful(true)
				readFmcSpecificIsSuccessful(true)
				readFmcDomainInfoIsSuccessful(true)
				readFtdIsSuccessful(true)
				readFmcDeviceRecordsIsSuccessful(true)
				readFtdDeviceRecordIsSuccessful(true)
				triggerFtdDeleteOnFmcIsSuccessful(false)
				waitForFtdDeleteStateMachineEndedSuccessful(true)
				deleteFtdIsSuccessful(true)
			},
			assertFunc: func(output *cloudftd.DeleteOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "error when failed to wait for FTD delete state machine starts",
			input:    cloudftd.NewDeleteInput(ftdUid),
			setupFunc: func() {
				readFmcIsSuccessful(true)
				readFmcSpecificIsSuccessful(true)
				readFmcDomainInfoIsSuccessful(true)
				readFtdIsSuccessful(true)
				readFmcDeviceRecordsIsSuccessful(true)
				readFtdDeviceRecordIsSuccessful(true)
				triggerFtdDeleteOnFmcIsSuccessful(true)
				waitForFtdDeleteStateMachineEndedSuccessful(false)
				deleteFtdIsSuccessful(true)
			},
			assertFunc: func(output *cloudftd.DeleteOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "error when failed to read FMC domain info",
			input:    cloudftd.NewDeleteInput(ftdUid),
			setupFunc: func() {
				readFmcIsSuccessful(true)
				readFmcSpecificIsSuccessful(true)
				readFmcDomainInfoIsSuccessful(false)
				readFtdIsSuccessful(true)
				readFmcDeviceRecordsIsSuccessful(true)
				readFtdDeviceRecordIsSuccessful(true)
				triggerFtdDeleteOnFmcIsSuccessful(true)
				waitForFtdDeleteStateMachineEndedSuccessful(true)
				deleteFtdIsSuccessful(true)
			},
			assertFunc: func(output *cloudftd.DeleteOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "error when failed to read FTD",
			input:    cloudftd.NewDeleteInput(ftdUid),
			setupFunc: func() {
				readFmcIsSuccessful(true)
				readFmcSpecificIsSuccessful(true)
				readFmcDomainInfoIsSuccessful(true)
				readFtdIsSuccessful(false)
				readFmcDeviceRecordsIsSuccessful(true)
				readFtdDeviceRecordIsSuccessful(true)
				triggerFtdDeleteOnFmcIsSuccessful(true)
				waitForFtdDeleteStateMachineEndedSuccessful(true)
				deleteFtdIsSuccessful(true)
			},
			assertFunc: func(output *cloudftd.DeleteOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "error when failed to read FMC device records",
			input:    cloudftd.NewDeleteInput(ftdUid),
			setupFunc: func() {
				readFmcIsSuccessful(true)
				readFmcSpecificIsSuccessful(true)
				readFmcDomainInfoIsSuccessful(true)
				readFtdIsSuccessful(true)
				readFmcDeviceRecordsIsSuccessful(false)
				readFtdDeviceRecordIsSuccessful(true)
				triggerFtdDeleteOnFmcIsSuccessful(true)
				waitForFtdDeleteStateMachineEndedSuccessful(true)
				deleteFtdIsSuccessful(true)
			},
			assertFunc: func(output *cloudftd.DeleteOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "error when failed to read FTD device record in FMC",
			input:    cloudftd.NewDeleteInput(ftdUid),
			setupFunc: func() {
				readFmcIsSuccessful(true)
				readFmcSpecificIsSuccessful(true)
				readFmcDomainInfoIsSuccessful(true)
				readFtdIsSuccessful(true)
				readFmcDeviceRecordsIsSuccessful(true)
				readFtdDeviceRecordIsSuccessful(false)
				triggerFtdDeleteOnFmcIsSuccessful(true)
				waitForFtdDeleteStateMachineEndedSuccessful(true)
				deleteFtdIsSuccessful(true)
			},
			assertFunc: func(output *cloudftd.DeleteOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "error when failed to delete FTD in CDO",
			input:    cloudftd.NewDeleteInput(ftdUid),
			setupFunc: func() {
				readFmcIsSuccessful(true)
				readFmcSpecificIsSuccessful(true)
				readFmcDomainInfoIsSuccessful(true)
				readFtdIsSuccessful(true)
				readFmcDeviceRecordsIsSuccessful(true)
				readFtdDeviceRecordIsSuccessful(true)
				triggerFtdDeleteOnFmcIsSuccessful(true)
				waitForFtdDeleteStateMachineEndedSuccessful(true)
				deleteFtdIsSuccessful(false)
			},
			assertFunc: func(output *cloudftd.DeleteOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := cloudftd.Delete(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

func deleteFtdIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodDelete,
			url.DeleteDevice(baseUrl, ftdUid),
			httpmock.NewJsonResponderOrPanic(200, validDeleteOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodDelete,
			url.DeleteDevice(baseUrl, ftdUid),
			httpmock.NewJsonResponderOrPanic(500, "intentional error"),
		)
	}
}

func readFtdDeviceRecordIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcDeviceRecord(baseUrl, fmcDomainUid, ftdDeviceRecordId),
			httpmock.NewJsonResponderOrPanic(200, validReadFtdDeviceRecordOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcDeviceRecord(baseUrl, fmcDomainUid, ftdDeviceRecordId),
			httpmock.NewJsonResponderOrPanic(500, "intentional error"),
		)
	}
}

func readFmcDeviceRecordsIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcAllDeviceRecords(baseUrl, fmcDomainUid),
			httpmock.NewJsonResponderOrPanic(200, validReadDeviceRecordsOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcAllDeviceRecords(baseUrl, fmcDomainUid),
			httpmock.NewJsonResponderOrPanic(500, "intentional error"),
		)
	}
}

func readFtdIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadDevice(baseUrl, ftdUid),
			httpmock.NewJsonResponderOrPanic(200, validReadFtdOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadDevice(baseUrl, ftdUid),
			httpmock.NewJsonResponderOrPanic(500, "intentional error"),
		)
	}
}

func waitForFtdDeleteStateMachineEndedSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadStateMachineInstance(baseUrl),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, []statemachine.ReadInstanceByDeviceUidOutput{
				validReadStateMachineOutput,
			}),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodPut,
			url.UpdateFmcAppliance(baseUrl, fmcApplianceUid),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func waitForFtdDeleteStateMachineTriggeredReturnedNotFound() {
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadStateMachineInstance(baseUrl),
		httpmock.NewJsonResponderOrPanic(http.StatusNotFound, statemachine.NotFoundError),
	)
}

func triggerFtdDeleteOnFmcIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodPut,
			url.UpdateFmcAppliance(baseUrl, fmcApplianceUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validUpdateFmcSpecificOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodPut,
			url.UpdateFmcAppliance(baseUrl, fmcApplianceUid),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func readFmcSpecificIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadSpecificDevice(baseUrl, fmcUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadSpecificOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadSpecificDevice(baseUrl, fmcUid),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

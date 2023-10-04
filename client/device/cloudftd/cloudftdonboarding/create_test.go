package cloudftdonboarding_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd/cloudftdonboarding"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestCloudFtdOnboardingCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := []struct {
		testName   string
		input      cloudftdonboarding.CreateInput
		setupFunc  func(t *testing.T)
		assertFunc func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T)
	}{
		{
			testName: "successful ftd onboarding, with no existing records",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				ReadFmcIsSuccessful(true)
				ReadFmcDomainInfoIsSuccessful(true)
				CheckFtdDuplicateIsSuccessful_NoExistingRecords(true)
				ReadApiTokenInfoIsSuccessful(true)
				CreateSystemApiTokenIsSuccessful(true)
				ReadFtdMetadataIsSuccessful(true)
				CreateFmcDeviceRecordIsSuccessful(true)
				ReadTaskStatusIsSuccessful(true)
				ReadFtdSpecificDeviceIsSuccessful(true)
				TriggerRegisterFmcStateMachineSuccess(true)
				TriggerRegisterFmcStateMachineEndsInDone(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, *output, validCreateFmcDeviceRecordOutput)
			},
		},
		{
			testName: "successful ftd onboarding, with existing records of different name",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				ReadFmcIsSuccessful(true)
				ReadFmcDomainInfoIsSuccessful(true)
				CheckFtdDuplicateIsSuccessful_ExistingRecordsWithDifferentName(true)
				ReadApiTokenInfoIsSuccessful(true)
				CreateSystemApiTokenIsSuccessful(true)
				ReadFtdMetadataIsSuccessful(true)
				CreateFmcDeviceRecordIsSuccessful(true)
				ReadTaskStatusIsSuccessful(true)
				ReadFtdSpecificDeviceIsSuccessful(true)
				TriggerRegisterFmcStateMachineSuccess(true)
				TriggerRegisterFmcStateMachineEndsInDone(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, *output, validCreateFmcDeviceRecordOutput)
			},
		},
		{
			testName: "successful ftd onboarding, with existing records of same name and different type",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				ReadFmcIsSuccessful(true)
				ReadFmcDomainInfoIsSuccessful(true)
				CheckFtdDuplicateIsSuccessful_ExistingRecordsWithSameNameButDifferentType(true)
				ReadApiTokenInfoIsSuccessful(true)
				CreateSystemApiTokenIsSuccessful(true)
				ReadFtdMetadataIsSuccessful(true)
				CreateFmcDeviceRecordIsSuccessful(true)
				ReadTaskStatusIsSuccessful(true)
				ReadFtdSpecificDeviceIsSuccessful(true)
				TriggerRegisterFmcStateMachineSuccess(true)
				TriggerRegisterFmcStateMachineEndsInDone(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, *output, validCreateFmcDeviceRecordOutput)
			},
		},
		{
			testName: "error when there is existing record of same name and same type",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				ReadFmcIsSuccessful(true)
				ReadFmcDomainInfoIsSuccessful(true)
				CheckFtdDuplicateIsFail_ExistingRecordsWithSameNameAndSameType()
				ReadApiTokenInfoIsSuccessful(true)
				CreateSystemApiTokenIsSuccessful(true)
				ReadFtdMetadataIsSuccessful(true)
				CreateFmcDeviceRecordIsSuccessful(true)
				ReadTaskStatusIsSuccessful(true)
				ReadFtdSpecificDeviceIsSuccessful(true)
				TriggerRegisterFmcStateMachineSuccess(true)
				TriggerRegisterFmcStateMachineEndsInDone(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.True(t, strings.Contains(err.Error(), "is already registered"))
			},
		},
		{
			testName: "error when read fmc failed",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				ReadFmcIsSuccessful(false)
				ReadFmcDomainInfoIsSuccessful(true)
				CheckFtdDuplicateIsSuccessful_NoExistingRecords(true)
				ReadApiTokenInfoIsSuccessful(true)
				CreateSystemApiTokenIsSuccessful(true)
				ReadFtdMetadataIsSuccessful(true)
				CreateFmcDeviceRecordIsSuccessful(true)
				ReadTaskStatusIsSuccessful(true)
				ReadFtdSpecificDeviceIsSuccessful(true)
				TriggerRegisterFmcStateMachineSuccess(true)
				TriggerRegisterFmcStateMachineEndsInDone(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.True(t, strings.Contains(err.Error(), url.ReadAllDevicesByType(baseUrl)))
			},
		},
		{
			testName: "error when read fmc domain info failed",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				ReadFmcIsSuccessful(true)
				ReadFmcDomainInfoIsSuccessful(false)
				CheckFtdDuplicateIsSuccessful_NoExistingRecords(true)
				ReadApiTokenInfoIsSuccessful(true)
				CreateSystemApiTokenIsSuccessful(true)
				ReadFtdMetadataIsSuccessful(true)
				CreateFmcDeviceRecordIsSuccessful(true)
				ReadTaskStatusIsSuccessful(true)
				ReadFtdSpecificDeviceIsSuccessful(true)
				TriggerRegisterFmcStateMachineSuccess(true)
				TriggerRegisterFmcStateMachineEndsInDone(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.True(t, strings.Contains(err.Error(), url.ReadFmcDomainInfo(fmcHost)))
			},
		},
		{
			testName: "error when read all fmc device records failed",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				ReadFmcIsSuccessful(true)
				ReadFmcDomainInfoIsSuccessful(true)
				CheckFtdDuplicateIsSuccessful_NoExistingRecords(false)
				ReadApiTokenInfoIsSuccessful(true)
				CreateSystemApiTokenIsSuccessful(true)
				ReadFtdMetadataIsSuccessful(true)
				CreateFmcDeviceRecordIsSuccessful(true)
				ReadTaskStatusIsSuccessful(true)
				ReadFtdSpecificDeviceIsSuccessful(true)
				TriggerRegisterFmcStateMachineSuccess(true)
				TriggerRegisterFmcStateMachineEndsInDone(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.True(t, strings.Contains(err.Error(), url.ReadFmcAllDeviceRecords(baseUrl, fmcDomainUid)))
			},
		},
		{
			testName: "error when read api token info failed",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				ReadFmcIsSuccessful(true)
				ReadFmcDomainInfoIsSuccessful(true)
				CheckFtdDuplicateIsSuccessful_NoExistingRecords(true)
				ReadApiTokenInfoIsSuccessful(false)
				CreateSystemApiTokenIsSuccessful(true)
				ReadFtdMetadataIsSuccessful(true)
				CreateFmcDeviceRecordIsSuccessful(true)
				ReadTaskStatusIsSuccessful(true)
				ReadFtdSpecificDeviceIsSuccessful(true)
				TriggerRegisterFmcStateMachineSuccess(true)
				TriggerRegisterFmcStateMachineEndsInDone(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.True(t, strings.Contains(err.Error(), url.ReadTokenInfo(baseUrl)))
			},
		},
		{
			testName: "error when create system token failed",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				ReadFmcIsSuccessful(true)
				ReadFmcDomainInfoIsSuccessful(true)
				CheckFtdDuplicateIsSuccessful_NoExistingRecords(true)
				ReadApiTokenInfoIsSuccessful(true)
				CreateSystemApiTokenIsSuccessful(false)
				ReadFtdMetadataIsSuccessful(true)
				CreateFmcDeviceRecordIsSuccessful(true)
				ReadTaskStatusIsSuccessful(true)
				ReadFtdSpecificDeviceIsSuccessful(true)
				TriggerRegisterFmcStateMachineSuccess(true)
				TriggerRegisterFmcStateMachineEndsInDone(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.True(t, strings.Contains(err.Error(), url.CreateSystemToken(baseUrl, systemTokenScope)))
			},
		},
		{
			testName: "error when read FTD metadata failed",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				ReadFmcIsSuccessful(true)
				ReadFmcDomainInfoIsSuccessful(true)
				CheckFtdDuplicateIsSuccessful_NoExistingRecords(true)
				ReadApiTokenInfoIsSuccessful(true)
				CreateSystemApiTokenIsSuccessful(true)
				ReadFtdMetadataIsSuccessful(false)
				CreateFmcDeviceRecordIsSuccessful(true)
				ReadTaskStatusIsSuccessful(true)
				ReadFtdSpecificDeviceIsSuccessful(true)
				TriggerRegisterFmcStateMachineSuccess(true)
				TriggerRegisterFmcStateMachineEndsInDone(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.True(t, strings.Contains(err.Error(), url.ReadDevice(baseUrl, ftdUid)))
			},
		},
		{
			testName: "error when create FTD device record failed",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				t.Skip("requires override inner retry config support")
				ReadFmcIsSuccessful(true)
				ReadFmcDomainInfoIsSuccessful(true)
				CheckFtdDuplicateIsSuccessful_NoExistingRecords(true)
				ReadApiTokenInfoIsSuccessful(true)
				CreateSystemApiTokenIsSuccessful(true)
				ReadFtdMetadataIsSuccessful(true)
				CreateFmcDeviceRecordIsSuccessful(false)
				ReadTaskStatusIsSuccessful(true)
				ReadFtdSpecificDeviceIsSuccessful(true)
				TriggerRegisterFmcStateMachineSuccess(true)
				TriggerRegisterFmcStateMachineEndsInDone(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.True(t, strings.Contains(err.Error(), url.CreateFmcDeviceRecord(baseUrl, fmcDomainUid)))
			},
		},
		{
			testName: "error when read task status failed",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				t.Skip("requires override inner retry config support")
				ReadFmcIsSuccessful(true)
				ReadFmcDomainInfoIsSuccessful(true)
				CheckFtdDuplicateIsSuccessful_NoExistingRecords(true)
				ReadApiTokenInfoIsSuccessful(true)
				CreateSystemApiTokenIsSuccessful(true)
				ReadFtdMetadataIsSuccessful(true)
				CreateFmcDeviceRecordIsSuccessful(true)
				ReadTaskStatusIsSuccessful(false)
				ReadFtdSpecificDeviceIsSuccessful(true)
				TriggerRegisterFmcStateMachineSuccess(true)
				TriggerRegisterFmcStateMachineEndsInDone(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.True(t, strings.Contains(err.Error(), url.ReadFmcTaskStatus(baseUrl, fmcDomainUid, fmcCreateDeviceTaskId)))

			},
		},
		{
			testName: "error when read FTD specific device failed",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				ReadFmcIsSuccessful(true)
				ReadFmcDomainInfoIsSuccessful(true)
				CheckFtdDuplicateIsSuccessful_NoExistingRecords(true)
				ReadApiTokenInfoIsSuccessful(true)
				CreateSystemApiTokenIsSuccessful(true)
				ReadFtdMetadataIsSuccessful(true)
				CreateFmcDeviceRecordIsSuccessful(true)
				ReadTaskStatusIsSuccessful(true)
				ReadFtdSpecificDeviceIsSuccessful(false)
				TriggerRegisterFmcStateMachineSuccess(true)
				TriggerRegisterFmcStateMachineEndsInDone(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.True(t, strings.Contains(err.Error(), url.ReadSpecificDevice(baseUrl, ftdUid)))
			},
		},
		{
			testName: "error when update FTD specific device failed",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				ReadFmcIsSuccessful(true)
				ReadFmcDomainInfoIsSuccessful(true)
				CheckFtdDuplicateIsSuccessful_NoExistingRecords(true)
				ReadApiTokenInfoIsSuccessful(true)
				CreateSystemApiTokenIsSuccessful(true)
				ReadFtdMetadataIsSuccessful(true)
				CreateFmcDeviceRecordIsSuccessful(true)
				ReadTaskStatusIsSuccessful(true)
				ReadFtdSpecificDeviceIsSuccessful(true)
				TriggerRegisterFmcStateMachineSuccess(false)
				TriggerRegisterFmcStateMachineEndsInDone(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.True(t, strings.Contains(err.Error(), url.UpdateSpecificCloudFtd(baseUrl, ftdSpecificUid)))
			},
		},
		{
			testName: "error when read FTD specific device failed",
			input:    cloudftdonboarding.NewCreateInput(ftdUid),
			setupFunc: func(t *testing.T) {
				t.Skip("requires override inner retry config support")
				ReadFmcIsSuccessful(true)
				ReadFmcDomainInfoIsSuccessful(true)
				CheckFtdDuplicateIsSuccessful_NoExistingRecords(true)
				ReadApiTokenInfoIsSuccessful(true)
				CreateSystemApiTokenIsSuccessful(true)
				ReadFtdMetadataIsSuccessful(true)
				CreateFmcDeviceRecordIsSuccessful(true)
				ReadTaskStatusIsSuccessful(true)
				ReadFtdSpecificDeviceIsSuccessful(true)
				TriggerRegisterFmcStateMachineSuccess(true)
				TriggerRegisterFmcStateMachineEndsInDone(false)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.True(t, strings.Contains(err.Error(), url.ReadSpecificDevice(baseUrl, ftdSpecificUid)))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(t)

			output, err := cloudftdonboarding.Create(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

func ReadFmcIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadAllDevicesByType(baseUrl),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadFmcOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadAllDevicesByType(baseUrl),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func ReadFmcDomainInfoIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcDomainInfo(fmcHost),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadDomainInfo),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcDomainInfo(fmcHost),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func ReadApiTokenInfoIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadTokenInfo(baseUrl),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadApiTokenInfo),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadTokenInfo(baseUrl),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func CheckFtdDuplicateIsSuccessful_NoExistingRecords(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcAllDeviceRecords(baseUrl, fmcDomainUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, allDeviceRecords_NoExistingRecords),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcAllDeviceRecords(baseUrl, fmcDomainUid),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func CheckFtdDuplicateIsSuccessful_ExistingRecordsWithDifferentName(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcAllDeviceRecords(baseUrl, fmcDomainUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, allDeviceRecords_ExistingRecords_DifferentName),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcAllDeviceRecords(baseUrl, fmcDomainUid),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func CheckFtdDuplicateIsSuccessful_ExistingRecordsWithSameNameButDifferentType(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcAllDeviceRecords(baseUrl, fmcDomainUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, allDeviceRecords_ExistingRecords_SameName),
		)
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcDeviceRecord(baseUrl, fmcDomainUid, allDeviceRecords_ExistingRecords_SameName.Items[0].Id),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, deviceRecord_NotFtd),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcAllDeviceRecords(baseUrl, fmcDomainUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, allDeviceRecords_ExistingRecords_SameName),
		)
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcDeviceRecord(baseUrl, fmcDomainUid, allDeviceRecords_ExistingRecords_SameName.Items[0].Id),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func CheckFtdDuplicateIsFail_ExistingRecordsWithSameNameAndSameType() {
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadFmcAllDeviceRecords(baseUrl, fmcDomainUid),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, allDeviceRecords_ExistingRecords_SameName),
	)
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadFmcDeviceRecord(baseUrl, fmcDomainUid, allDeviceRecords_ExistingRecords_SameName.Items[0].Id),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, deviceRecord_IsFtd),
	)
}

func CreateSystemApiTokenIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodPost,
			url.CreateSystemToken(baseUrl, systemTokenScope),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validCreateSystemApiTokenOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodPost,
			url.CreateSystemToken(baseUrl, systemTokenScope),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func ReadFtdMetadataIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadDevice(baseUrl, ftdUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadFtdOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadDevice(baseUrl, ftdUid),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func CreateFmcDeviceRecordIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodPost,
			url.CreateFmcDeviceRecord(baseUrl, fmcDomainUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validCreateFmcDeviceRecordOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodPost,
			url.CreateFmcDeviceRecord(baseUrl, fmcDomainUid),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func ReadTaskStatusIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcTaskStatus(baseUrl, fmcDomainUid, fmcCreateDeviceTaskId),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadTaskOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcTaskStatus(baseUrl, fmcDomainUid, fmcCreateDeviceTaskId),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func ReadFtdSpecificDeviceIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadSpecificDevice(baseUrl, ftdUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadSpecificOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadSpecificDevice(baseUrl, ftdUid),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func TriggerRegisterFmcStateMachineSuccess(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodPut,
			url.UpdateSpecificCloudFtd(baseUrl, ftdSpecificUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validUpdateSpecificUidOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodPut,
			url.UpdateSpecificCloudFtd(baseUrl, ftdSpecificUid),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func TriggerRegisterFmcStateMachineEndsInDone(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadSpecificDevice(baseUrl, ftdSpecificUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadFtdSpecificOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadSpecificDevice(baseUrl, ftdSpecificUid),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

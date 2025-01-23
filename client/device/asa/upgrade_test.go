package asa_test

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactionstatus"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactiontype"
	internalTesting "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestValidateVersionCompatibility(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := []struct {
		testName        string
		deviceUid       string
		softwareVersion string
		asdmVersion     string
		setupFunc       func(deviceUid string, softwareVersion string, asdmVersion string)
		assertFunc      func(err error, t *testing.T)
	}{
		{
			testName:        "should not fail if versions are compatible with the ASA",
			deviceUid:       uuid.New().String(),
			softwareVersion: "9.8.4",
			asdmVersion:     "7.8.2",
			setupFunc: func(deviceUid string, softwareVersion string, asdmVersion string) {
				configureCompatibleVersionsToRespondSuccessfully(deviceUid, model.CdoListResponse[asa.CompatibleVersion]{
					Items: []asa.CompatibleVersion{
						{SoftwareVersion: "9.16(6)100", AsdmVersion: "7.12(2)"},
						{SoftwareVersion: softwareVersion, AsdmVersion: asdmVersion},
					},
					Count: 2,
				})
			},
			assertFunc: func(err error, t *testing.T) {
				assert.Nil(t, err)
			},
		},
		{
			testName:        "should not fail if software version nil and ASDM version compatible with the ASA",
			deviceUid:       uuid.New().String(),
			softwareVersion: "",
			asdmVersion:     "7.16(3.100)",
			setupFunc: func(deviceUid string, softwareVersion string, asdmVersion string) {
				configureCompatibleVersionsToRespondSuccessfully(deviceUid, model.CdoListResponse[asa.CompatibleVersion]{
					Items: []asa.CompatibleVersion{
						{SoftwareVersion: "9.8.4", AsdmVersion: asdmVersion},
						{SoftwareVersion: "9.16(6)100", AsdmVersion: "7.12(2)"},
					},
					Count: 2,
				})
			},
			assertFunc: func(err error, t *testing.T) {
				assert.Nil(t, err)
			},
		},
		{
			testName:        "should not fail if ASDM version nil and software version compatible with the ASA",
			deviceUid:       uuid.New().String(),
			softwareVersion: "9.18(2)",
			asdmVersion:     "",
			setupFunc: func(deviceUid string, softwareVersion string, asdmVersion string) {
				configureCompatibleVersionsToRespondSuccessfully(deviceUid, model.CdoListResponse[asa.CompatibleVersion]{
					Items: []asa.CompatibleVersion{
						{SoftwareVersion: "9.8.4", AsdmVersion: asdmVersion},
						{SoftwareVersion: softwareVersion, AsdmVersion: "7.12(2)"},
					},
					Count: 2,
				})
			},
			assertFunc: func(err error, t *testing.T) {
				assert.Nil(t, err)
			},
		},
		{
			testName:        "should fail if software versions and ASDM are not compatible with the ASA",
			deviceUid:       uuid.New().String(),
			softwareVersion: "9.18(2)",
			asdmVersion:     "7.16(3.100)",
			setupFunc: func(deviceUid string, softwareVersion string, asdmVersion string) {
				configureCompatibleVersionsToRespondSuccessfully(deviceUid, model.CdoListResponse[asa.CompatibleVersion]{
					Items: []asa.CompatibleVersion{
						{SoftwareVersion: "9.16(6)100", AsdmVersion: "7.12(2)"},
					},
					Count: 1,
				})
			},
			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), "Device cannot be upgraded to the specified software and ASDM versions.")
			},
		},
		{
			testName:        "should fail if ASDM version nil and software version not compatible with the ASA",
			deviceUid:       uuid.New().String(),
			softwareVersion: "9.18(2)",
			asdmVersion:     "",
			setupFunc: func(deviceUid string, softwareVersion string, asdmVersion string) {
				configureCompatibleVersionsToRespondSuccessfully(deviceUid, model.CdoListResponse[asa.CompatibleVersion]{
					Items: []asa.CompatibleVersion{
						{SoftwareVersion: "9.16(6)100", AsdmVersion: "7.12(2)"},
					},
					Count: 1,
				})
			},
			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), "Device cannot be upgraded to the specified software and ASDM versions.")
			},
		},
		{
			testName:        "should fail if software version nil and ASDM version not compatible with the ASA",
			deviceUid:       uuid.New().String(),
			softwareVersion: "",
			asdmVersion:     "7.16(3.100)",
			setupFunc: func(deviceUid string, softwareVersion string, asdmVersion string) {
				configureCompatibleVersionsToRespondSuccessfully(deviceUid, model.CdoListResponse[asa.CompatibleVersion]{
					Items: []asa.CompatibleVersion{
						{SoftwareVersion: "9.16(6)100", AsdmVersion: "7.12(2)"},
					},
					Count: 1,
				})
			},
			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), "Device cannot be upgraded to the specified software and ASDM versions.")
			},
		},
		{
			testName:        "should fail if http response is not successful",
			deviceUid:       uuid.New().String(),
			softwareVersion: "9.18(2)",
			asdmVersion:     "7.16(3.100)",
			setupFunc: func(deviceUid string, softwareVersion string, asdmVersion string) {
				configureCompatibleVersionsToFailToRespond(deviceUid)
			},
			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.deviceUid, testCase.softwareVersion, testCase.asdmVersion)

			err := asa.ValidateVersionCompatibility(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.deviceUid,
				testCase.softwareVersion,
				testCase.asdmVersion,
			)

			testCase.assertFunc(err, t)
		})
	}
}

func TestUpgradeAsa(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := []struct {
		testName        string
		deviceUid       string
		softwareVersion string
		asdmVersion     string
		setupFunc       func(deviceUid string)
		assertFunc      func(err error, t *testing.T)
	}{
		{
			testName:        "should not fail if upgrade is successful",
			deviceUid:       uuid.New().String(),
			softwareVersion: "9.18(2)",
			asdmVersion:     "7.16(3.100)",
			setupFunc: func(deviceUid string) {
				transactionUid := uuid.New().String()
				inProgressTransaction := transaction.Type{
					TransactionUid:  uuid.New().String(),
					TenantUid:       uuid.New().String(),
					EntityUid:       uuid.New().String(),
					EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/inventory/devices/" + deviceUid,
					PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
					SubmissionTime:  "2024-09-10T20:10:00Z",
					LastUpdatedTime: "2024-10-10T20:10:00Z",
					Type:            transactiontype.UPGRADE_ASA,
					Status:          transactionstatus.IN_PROGRESS,
				}
				doneTransaction := transaction.Type{
					TransactionUid:  inProgressTransaction.TransactionUid,
					TenantUid:       inProgressTransaction.TenantUid,
					EntityUid:       inProgressTransaction.EntityUid,
					EntityUrl:       inProgressTransaction.EntityUrl,
					PollingUrl:      inProgressTransaction.PollingUrl,
					SubmissionTime:  inProgressTransaction.SubmissionTime,
					LastUpdatedTime: "2024-10-10T20:11:00Z",
					Type:            inProgressTransaction.Type,
					Status:          transactionstatus.DONE,
				}
				configureUpgradeAsaToRespondSuccessfully(deviceUid, inProgressTransaction)
				internalTesting.MockGetOk(fmt.Sprintf("%s/api/rest/v1/transactions/%s", baseUrl, transactionUid), doneTransaction)
			},
			assertFunc: func(err error, t *testing.T) {
				assert.Nil(t, err)
			},
		},
		{
			testName:        "should fail if the transaction is in error state",
			deviceUid:       uuid.New().String(),
			softwareVersion: "9.18(2)",
			asdmVersion:     "7.16(3.100)",
			setupFunc: func(deviceUid string) {
				transactionUid := uuid.New().String()
				inProgressTransaction := transaction.Type{
					TransactionUid:  uuid.New().String(),
					TenantUid:       uuid.New().String(),
					EntityUid:       uuid.New().String(),
					EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/inventory/devices/" + deviceUid,
					PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
					SubmissionTime:  "2024-09-10T20:10:00Z",
					LastUpdatedTime: "2024-10-10T20:10:00Z",
					Type:            transactiontype.UPGRADE_ASA,
					Status:          transactionstatus.IN_PROGRESS,
				}
				errorTransaction := transaction.Type{
					TransactionUid:  inProgressTransaction.TransactionUid,
					TenantUid:       inProgressTransaction.TenantUid,
					EntityUid:       inProgressTransaction.EntityUid,
					EntityUrl:       inProgressTransaction.EntityUrl,
					PollingUrl:      inProgressTransaction.PollingUrl,
					SubmissionTime:  inProgressTransaction.SubmissionTime,
					LastUpdatedTime: "2024-10-10T20:11:00Z",
					Type:            inProgressTransaction.Type,
					Status:          transactionstatus.ERROR,
				}
				configureUpgradeAsaToRespondSuccessfully(deviceUid, inProgressTransaction)
				internalTesting.MockGetOk(fmt.Sprintf("%s/api/rest/v1/transactions/%s", baseUrl, transactionUid), errorTransaction)
			},
			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)
			},
		},
		{
			testName:        "should fail if the upgrade ASA request returns an error transaction",
			deviceUid:       uuid.New().String(),
			softwareVersion: "9.18(2)",
			asdmVersion:     "7.16(3.100)",
			setupFunc: func(deviceUid string) {
				transactionUid := uuid.New().String()
				errorTransaction := transaction.Type{
					TransactionUid:  uuid.New().String(),
					TenantUid:       uuid.New().String(),
					EntityUid:       uuid.New().String(),
					EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/inventory/devices/" + deviceUid,
					PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
					SubmissionTime:  "2024-09-10T20:10:00Z",
					LastUpdatedTime: "2024-10-10T20:10:00Z",
					Type:            transactiontype.UPGRADE_ASA,
					Status:          transactionstatus.ERROR,
				}
				configureUpgradeAsaToRespondSuccessfully(deviceUid, errorTransaction)
			},
			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)
			},
		},
		{
			testName:        "should fail if the upgrade ASA http request is not successful",
			deviceUid:       uuid.New().String(),
			softwareVersion: "9.18(2)",
			asdmVersion:     "7.16(3.100)",
			setupFunc: func(deviceUid string) {
				configureUpgradeAsaToFailToRespond(deviceUid)
			},
			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)
			},
		},
		{
			testName:        "should fail if the upgrade ASA transaction polling http request is not successful",
			deviceUid:       uuid.New().String(),
			softwareVersion: "9.18(2)",
			asdmVersion:     "7.16(3.100)",
			setupFunc: func(deviceUid string) {
				transactionUid := uuid.New().String()
				inProgressTransaction := transaction.Type{
					TransactionUid:  uuid.New().String(),
					TenantUid:       uuid.New().String(),
					EntityUid:       uuid.New().String(),
					EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/inventory/devices/" + deviceUid,
					PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
					SubmissionTime:  "2024-09-10T20:10:00Z",
					LastUpdatedTime: "2024-10-10T20:10:00Z",
					Type:            transactiontype.UPGRADE_ASA,
					Status:          transactionstatus.IN_PROGRESS,
				}
				configureUpgradeAsaToRespondSuccessfully(deviceUid, inProgressTransaction)
				internalTesting.MockGetError(fmt.Sprintf("%s/api/rest/v1/transactions/%s", baseUrl, transactionUid), "a failed transaction error string")
			},
			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), "a failed transaction error string")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.deviceUid)

			err := asa.UpgradeAsa(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.deviceUid,
				testCase.softwareVersion,
				testCase.asdmVersion,
			)

			testCase.assertFunc(err, t)
		})
	}
}

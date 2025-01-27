package cloudftd_test

import (
	"errors"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	mockhttp "net/http"
	"testing"
	"time"
)

var upgradePackages []cloudftd.UpgradePackage = []cloudftd.UpgradePackage{
	{
		UpgradePackageUid: uuid.New().String(),
		SoftwareVersion:   "7.2.5.1-29",
	},
	{
		UpgradePackageUid: uuid.New().String(),
		SoftwareVersion:   "7.2.6-293",
	},
}

func TestUpgrade(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := []struct {
		testName          string
		uid               string
		softwareVersion   string
		expectedFtdDevice *cloudftd.FtdDevice
		expectedError     error
		setupFunc         func(deviceUid string, softwareVersion string, ftdDevice *cloudftd.FtdDevice)
		assertFunc        func(ftdDevice *cloudftd.FtdDevice, err error, expectedFtdDevice *cloudftd.FtdDevice, expectedError error, t *testing.T)
	}{
		{
			testName:          "Fail to upgrade if FTD device not found",
			uid:               uuid.New().String(),
			softwareVersion:   "7.2.5",
			expectedFtdDevice: nil,
			expectedError:     nil,
			setupFunc: func(deviceUid string, softwareVersion string, ftdDevice *cloudftd.FtdDevice) {
				httpmock.RegisterResponder(
					mockhttp.MethodGet,
					baseUrl+"/aegis/rest/v1/services/targets/devices/"+deviceUid,
					httpmock.NewStringResponder(404, ""),
				)
			},
			assertFunc: func(ftdDevice *cloudftd.FtdDevice, err error, expectedFtdDevice *cloudftd.FtdDevice, expectedError error, t *testing.T) {
				assert.Nil(t, ftdDevice)
				assert.NotNil(t, err)
			},
		},
		{
			testName:          "Fail to upgrade if FTD device connectivity state is not ONLINE",
			uid:               uuid.New().String(),
			softwareVersion:   "7.2.5",
			expectedFtdDevice: nil,
			expectedError:     errors.New("FTD device connectivity state is not ONLINE. Only ONLINE devices can be upgraded"),
			setupFunc: func(deviceUid string, softwareVersion string, ftdDevice *cloudftd.FtdDevice) {
				httpmock.RegisterResponder(
					mockhttp.MethodGet,
					baseUrl+"/aegis/rest/v1/services/targets/devices/"+deviceUid,
					httpmock.NewJsonResponderOrPanic(200, &cloudftd.FtdDevice{
						Uid:               deviceUid,
						Name:              "FTD Device",
						Metadata:          cloudftd.Metadata{},
						State:             "ACTIVE",
						DeviceType:        "FTDC",
						ConnectivityState: -3,
						Tags:              nil,
						SoftwareVersion:   "7.2.4",
					}),
				)
			},
			assertFunc: func(ftdDevice *cloudftd.FtdDevice, err error, expectedFtdDevice *cloudftd.FtdDevice, expectedError error, t *testing.T) {
				assert.Nil(t, ftdDevice)
				assert.NotNil(t, err)
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
		{
			testName:          "Fail to upgrade if device is not cdFMC-managed FTD",
			uid:               uuid.New().String(),
			softwareVersion:   "7.2.5",
			expectedFtdDevice: nil,
			expectedError:     errors.New("this resource only supports cdFMC managed FTDs"),
			setupFunc: func(deviceUid string, softwareVersion string, ftdDevice *cloudftd.FtdDevice) {
				httpmock.RegisterResponder(
					mockhttp.MethodGet,
					baseUrl+"/aegis/rest/v1/services/targets/devices/"+deviceUid,
					httpmock.NewJsonResponderOrPanic(200, &cloudftd.FtdDevice{
						Uid:               deviceUid,
						Name:              "FTD Device",
						Metadata:          cloudftd.Metadata{},
						State:             "ACTIVE",
						DeviceType:        "ASA",
						ConnectivityState: 1,
						Tags:              nil,
						SoftwareVersion:   "7.2.4",
					}),
				)
			},
			assertFunc: func(ftdDevice *cloudftd.FtdDevice, err error, expectedFtdDevice *cloudftd.FtdDevice, expectedError error, t *testing.T) {
				assert.Nil(t, ftdDevice)
				assert.NotNil(t, err)
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
		{
			testName:          "Fail to upgrade if FTD device software version is less than version to upgrade to",
			uid:               uuid.New().String(),
			softwareVersion:   "7.2.5",
			expectedFtdDevice: nil,
			expectedError:     errors.New("FTD device is on version 7.3.0, which is newer than the version to upgrade to: 7.2.5"),
			setupFunc: func(deviceUid string, softwareVersion string, ftdDevice *cloudftd.FtdDevice) {
				httpmock.RegisterResponder(
					mockhttp.MethodGet,
					baseUrl+"/aegis/rest/v1/services/targets/devices/"+deviceUid,
					httpmock.NewJsonResponderOrPanic(200, &cloudftd.FtdDevice{
						Uid:               deviceUid,
						Name:              "FTD Device",
						Metadata:          cloudftd.Metadata{},
						State:             "ACTIVE",
						DeviceType:        "FTDC",
						ConnectivityState: 1,
						Tags:              nil,
						SoftwareVersion:   "7.3.0",
					}),
				)
			},
			assertFunc: func(ftdDevice *cloudftd.FtdDevice, err error, expectedFtdDevice *cloudftd.FtdDevice, expectedError error, t *testing.T) {
				assert.Nil(t, ftdDevice)
				assert.NotNil(t, err)
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
		{
			testName:        "Do not fail if device software version is equal to version to upgrade to",
			uid:             uuid.New().String(),
			softwareVersion: "7.2.5",
			expectedFtdDevice: &cloudftd.FtdDevice{
				Uid:               uuid.New().String(),
				Name:              "FTD Device",
				Metadata:          cloudftd.Metadata{},
				State:             "ACTIVE",
				ConnectivityState: 1,
				DeviceType:        "FTDC",
				Tags:              nil,
				SoftwareVersion:   "7.2.5",
			},
			expectedError: nil,
			setupFunc: func(deviceUid string, softwareVersion string, ftdDevice *cloudftd.FtdDevice) {
				httpmock.RegisterResponder(mockhttp.MethodGet,
					baseUrl+"/aegis/rest/v1/services/targets/devices/"+deviceUid,
					httpmock.NewJsonResponderOrPanic(200, ftdDevice))
			},
			assertFunc: func(ftdDevice *cloudftd.FtdDevice, err error, expectedFtdDevice *cloudftd.FtdDevice, expectedError error, t *testing.T) {
				assert.NotNil(t, ftdDevice)
				assert.Nil(t, err)
				assert.Equal(t, expectedFtdDevice, ftdDevice)
			},
		},
		{
			testName:        "Upgrade FTD device - fail because the specified version is incompatible",
			uid:             uuid.New().String(),
			softwareVersion: "7.2.5",
			expectedFtdDevice: &cloudftd.FtdDevice{
				Uid:               uuid.New().String(),
				Name:              "FTD Device",
				Metadata:          cloudftd.Metadata{},
				State:             "ACTIVE",
				DeviceType:        "FTDC",
				ConnectivityState: 1,
				Tags:              nil,
				SoftwareVersion:   "7.2.3",
			},
			expectedError: errors.New("7.2.5 is not a valid version to upgrade FTD device FTD Device to"),
			setupFunc: func(deviceUid string, softwareVersion string, ftdDevice *cloudftd.FtdDevice) {
				httpmock.RegisterResponder(mockhttp.MethodGet,
					baseUrl+"/aegis/rest/v1/services/targets/devices/"+deviceUid,
					httpmock.NewJsonResponderOrPanic(200, ftdDevice))
				httpmock.RegisterResponder(mockhttp.MethodGet, baseUrl+"/api/rest/v1/inventory/devices/ftds/"+ftdDevice.Uid+"/upgrades/versions", httpmock.NewJsonResponderOrPanic(200, model.CdoListResponse[cloudftd.UpgradePackage]{
					Items: upgradePackages,
					Count: len(upgradePackages),
				}))
			},
			assertFunc: func(ftdDevice *cloudftd.FtdDevice, err error, expectedFtdDevice *cloudftd.FtdDevice, expectedError error, t *testing.T) {
				assert.Nil(t, ftdDevice)
				assert.NotNil(t, err)
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
		// TODO this test should be changed to a success test once the upgrade implementation is done
		{
			testName:        "Upgrade FTD device - fail because the code has not been implemented yet",
			uid:             uuid.New().String(),
			softwareVersion: "7.2.5.1-29",
			expectedFtdDevice: &cloudftd.FtdDevice{
				Uid:               uuid.New().String(),
				Name:              "FTD Device",
				Metadata:          cloudftd.Metadata{},
				State:             "ACTIVE",
				DeviceType:        "FTDC",
				ConnectivityState: 1,
				Tags:              nil,
				SoftwareVersion:   "7.2.3",
			},
			expectedError: errors.New("upgrade implementation coming soon"),
			setupFunc: func(deviceUid string, softwareVersion string, ftdDevice *cloudftd.FtdDevice) {
				httpmock.RegisterResponder(mockhttp.MethodGet,
					baseUrl+"/aegis/rest/v1/services/targets/devices/"+deviceUid,
					httpmock.NewJsonResponderOrPanic(200, ftdDevice))
				httpmock.RegisterResponder(mockhttp.MethodGet, baseUrl+"/api/rest/v1/inventory/devices/ftds/"+ftdDevice.Uid+"/upgrades/versions", httpmock.NewJsonResponderOrPanic(200, model.CdoListResponse[cloudftd.UpgradePackage]{
					Items: upgradePackages,
					Count: len(upgradePackages),
				}))
			},
			assertFunc: func(ftdDevice *cloudftd.FtdDevice, err error, expectedFtdDevice *cloudftd.FtdDevice, expectedError error, t *testing.T) {
				assert.Nil(t, ftdDevice)
				assert.NotNil(t, err)
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.uid, testCase.softwareVersion, testCase.expectedFtdDevice)

			ftdDevice, err := cloudftd.NewFtdUpgradeService(
				context.Background(),
				http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
			).Upgrade(
				testCase.uid,
				testCase.softwareVersion,
			)

			testCase.assertFunc(ftdDevice, err, testCase.expectedFtdDevice, testCase.expectedError, t)
		})
	}
}

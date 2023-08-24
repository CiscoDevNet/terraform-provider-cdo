package cdfmc_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cdfmc"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cdfmc/smartlicense"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestSmartLicenseRead(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validSmartLicenseMetadata := smartlicense.NewMetadata(
		smartLicenseAuthStatus,
		smartLicenseEvalExpiresInDays,
		smartLicenseEvalUsed,
		smartLicenseExportControl,
		smartLicenseVirtualAccount,
	)
	validSmartLicenseItems := smartlicense.NewItems(
		smartlicense.NewItem(
			validSmartLicenseMetadata,
			smartLicenseRegStatus,
			smartLicenseType,
		),
	)
	validSmartLicenseLinks := smartlicense.NewLinks(smartLicenseSelfLink)
	validSmartLicensePaging := smartlicense.NewPaging(
		smartLicenseCount,
		smartLicenseOffset,
		smartLicenseLimit,
		smartLicensePages,
	)
	validSmartLicense := smartlicense.NewSmartLicense(
		validSmartLicenseItems,
		validSmartLicenseLinks,
		validSmartLicensePaging,
	)

	testCases := []struct {
		testName   string
		setupFunc  func()
		assertFunc func(output *cdfmc.ReadSmartLicenseOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully read Smart License",
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadSmartLicense(baseUrl),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, validSmartLicense),
				)
			},
			assertFunc: func(output *cdfmc.ReadSmartLicenseOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validSmartLicense, *output)
			},
		},
		{
			testName: "return error when read Smart License error",
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.CreateDevice(baseUrl),
					httpmock.NewStringResponder(http.StatusInternalServerError, "internal server error"),
				)
			},
			assertFunc: func(output *cdfmc.ReadSmartLicenseOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := cdfmc.ReadSmartLicense(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				cdfmc.NewReadSmartLicenseInput(),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

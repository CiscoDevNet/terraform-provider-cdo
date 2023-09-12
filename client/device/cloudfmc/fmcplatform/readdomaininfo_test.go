package fmcplatform_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcplatform"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmcdomain"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"

	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestReadDomainInfo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validOutput := &fmcplatform.ReadDomainInfoOutput{
		Links:  fmcdomain.NewLinks(links),
		Paging: fmcdomain.NewPaging(count, offset, limit, pages),
		Items: []fmcdomain.Item{
			fmcdomain.NewItem(uuid, name, type_),
		},
	}

	testCases := []struct {
		testName    string
		fmcHostname string
		setupFunc   func()
		assertFunc  func(output *fmcplatform.ReadDomainInfoOutput, err error, t *testing.T)
	}{
		{
			testName:    "successfully read FMC domain info",
			fmcHostname: fmcHostname,

			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadFmcDomainInfo(fmcHostname),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, validOutput),
				)
			},

			assertFunc: func(output *fmcplatform.ReadDomainInfoOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validOutput, output)
			},
		},

		{
			testName:    "error when read FMC domain info error",
			fmcHostname: fmcHostname,

			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadFmcDomainInfo(fmcHostname),
					httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
				)
			},

			assertFunc: func(output *fmcplatform.ReadDomainInfoOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := fmcplatform.ReadFmcDomainInfo(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				fmcplatform.NewReadDomainInfoInput(fmcHostname),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

package fmcappliance_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc/fmcappliance"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"

	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validUpdateOutput := fmcappliance.NewUpdateOutputBuilder().
		Uid(uid).
		State(state).
		DomainUid(domainUid).
		Build()

	testCases := []struct {
		testName   string
		input      fmcappliance.UpdateInput
		setupFunc  func()
		assertFunc func(input fmcappliance.UpdateInput, output *fmcappliance.UpdateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully updates FMC Appliance name",
			input:    fmcappliance.NewUpdateInput(fmcApplianceUid, queueTriggerState, stateMachineContext),

			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodPut,
					url.UpdateFmcAppliance(baseUrl, fmcApplianceUid),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, validUpdateOutput),
				)
			},

			assertFunc: func(input fmcappliance.UpdateInput, output *fmcappliance.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validUpdateOutput, *output)
			},
		},

		{
			testName: "error when update FMC Appliance name error",
			input:    fmcappliance.NewUpdateInput(fmcApplianceUid, queueTriggerState, stateMachineContext),

			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodPut,
					url.UpdateFmcAppliance(baseUrl, fmcApplianceUid),
					httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
				)
			},

			assertFunc: func(input fmcappliance.UpdateInput, output *fmcappliance.UpdateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := fmcappliance.Update(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(testCase.input, output, err, t)
		})
	}
}

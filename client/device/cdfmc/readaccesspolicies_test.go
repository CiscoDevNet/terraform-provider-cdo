package cdfmc_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cdfmc"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cdfmc/accesspolicies"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestAccessPoliciesRead(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validAccessPolicesItems := []accesspolicies.Item{
		accesspolicies.NewItem(
			accessPolicyId,
			accessPolicyName,
			accessPolicyType,
			accesspolicies.NewLinks(accessPolicySelfLink),
		),
	}
	validAccessPoliciesPaging := accesspolicies.NewPaging(
		accessPolicyCount,
		accessPolicyOffset,
		accessPolicyLimit,
		accessPolicyPages,
	)
	validAccessPoliciesLink := accesspolicies.NewLinks(accessPolicySelfLink)

	validAccessPolicies := accesspolicies.New(
		validAccessPolicesItems,
		validAccessPoliciesLink,
		validAccessPoliciesPaging,
	)

	testCases := []struct {
		testName   string
		domainUid  string
		limit      int
		setupFunc  func()
		assertFunc func(output *cdfmc.ReadAccessPoliciesOutput, err error, t *testing.T)
	}{
		{
			testName:  "successfully read Access Policies",
			domainUid: domainUid,
			limit:     limit,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadAccessPolicies(baseUrl, domainUid, limit),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, validAccessPolicies),
				)
			},
			assertFunc: func(output *cdfmc.ReadAccessPoliciesOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validAccessPolicies, *output)
			},
		},
		{
			testName:  "return error when read access policy error",
			domainUid: domainUid,
			limit:     limit,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadAccessPolicies(baseUrl, domainUid, limit),
					httpmock.NewStringResponder(http.StatusInternalServerError, "internal server error"),
				)
			},
			assertFunc: func(output *cdfmc.ReadAccessPoliciesOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := cdfmc.ReadAccessPolicies(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				cdfmc.NewReadAccessPoliciesInput(fmcHostname, domainUid, limit),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

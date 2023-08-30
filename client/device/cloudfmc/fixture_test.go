package cloudfmc_test

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"time"
)

const (
	smartLicenseEvalExpiresInDays = 123456
	smartLicenseEvalUsed          = false
	smartLicenseExportControl     = false
	smartLicenseVirtualAccount    = "unit-test-virtual-account"
	smartLicenseAuthStatus        = "unit-test-auth-status"
	smartLicenseRegStatus         = "unit-test-reg-status"
	smartLicenseType              = "unit-test-type"
	smartLicenseSelfLink          = "unit-test-self-link"
	smartLicenseCount             = 123456
	smartLicenseOffset            = 123456
	smartLicenseLimit             = 123456
	smartLicensePages             = 123456

	baseUrl = "https://unit-test.net"

	fmcHostname = "https://fmc-hostname.unit-test.net"
	fmcUid      = "unit-test-fmc-uid"
	domainUid   = "unit-test-domain-uid"
	limit       = 123456
	status      = "unit-test-status"

	deviceName              = "unit-test-device-name"
	deviceUid               = "unit-test-uid"
	deviceHost              = "https://unit-test.com"
	devicePort              = 1234
	deviceCloudConnectorUId = "unit-test-uid"
	deviceState             = state.DONE

	specificDeviceUid = "unit-test-specific-device-uid"

	accessPolicySelfLink = "https://unit-test.cdo.cisco.com/api/fmc_config/v1/domain/unit-test-domain-uid/policy/accesspolicies/unit-test-uid"
	accessPolicyName     = "Unit Test Access Control Policy"
	accessPolicyType     = "UnitTestAccessPolicy"
	accessPolicyId       = "unit-test-id"
	accessPolicyCount    = 123456
	accessPolicyOffset   = 123456
	accessPolicyLimit    = 123456
	accessPolicyPages    = 123456
)

var (
	deviceCreatedDate     = time.Date(1999, 1, 1, 0, 0, 0, 0, time.Local)
	deviceLastUpdatedDate = time.Date(1999, 1, 1, 0, 0, 0, 0, time.Local)
)

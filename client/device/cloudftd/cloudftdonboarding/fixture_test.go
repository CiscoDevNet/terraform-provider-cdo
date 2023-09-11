package cloudftdonboarding_test

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
)

const (
	baseUrl = "https://unit-test.cdo.cisco.com"

	fmcName      = "unit-test-device-name"
	fmcUid       = "unit-test-uid"
	fmcDomainUid = "unit-test-domain-uid"
	fmcHost      = "unit-test-fmc-host.com"
	fmcPort      = 1234
	fmcLink      = "unit-test-fmc-link"
)

var (
	validReadFmcOutput = []device.ReadOutput{
		device.NewReadOutputBuilder().
			AsCloudFmc().
			WithName(fmcName).
			WithUid(fmcUid).
			WithLocation(fmcHost, fmcPort).
			Build(),
	}

	//validReadDomainInfo := &fmcplatform.ReadDomainInfoOutput{
	//	Links:  fmcdomain.NewLinks(links),
	//	Paging: fmcdomain.NewPaging(count, offset, limit, pages),
	//	Items: []fmcdomain.Item{
	//		fmcdomain.NewItem(uuid, name, type_),
	//	},
	//}
)

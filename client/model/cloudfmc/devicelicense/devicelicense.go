package devicelicense

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/internal"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
)

type DeviceLicense struct {
	Links  Links  `json:"links"`
	Paging Paging `json:"paging"`
	Items  []Item `json:"items"`
}

func NewDeviceLicense(links Links, paging Paging, items []Item) DeviceLicense {
	return DeviceLicense{
		Links:  links,
		Paging: paging,
		Items:  items,
	}
}

type Item struct {
	Id              string         `json:"id"`
	Type            string         `json:"type"`
	LicenseTypes    []license.Type `json:"licenseTypes"`
	PerformanceTier tier.Type      `json:"performanceTier"`
	Links           Links          `json:"links"`
}

func NewItem(id, type_ string, licenseTypes []license.Type, performanceTier tier.Type, links Links) Item {
	return Item{
		Type:            type_,
		Id:              id,
		LicenseTypes:    licenseTypes,
		PerformanceTier: performanceTier,
		Links:           links,
	}
}

type Links = internal.Links
type Paging = internal.Paging

var NewLinks = internal.NewLinks
var NewPaging = internal.NewPaging

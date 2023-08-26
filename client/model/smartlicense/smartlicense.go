package smartlicense

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmccommon"
)

type SmartLicense struct {
	Items  Items            `json:"items"`
	Links  fmccommon.Links  `json:"links"`
	Paging fmccommon.Paging `json:"paging"`
}

func NewSmartLicense(items Items, links fmccommon.Links, paging fmccommon.Paging) SmartLicense {
	return SmartLicense{
		Items:  items,
		Links:  links,
		Paging: paging,
	}
}

type Items struct {
	Items []Item `json:"items"`
}

func NewItems(items ...Item) Items {
	return Items{
		Items: items,
	}
}

type Item struct {
	Metadata  Metadata `json:"metadata"`
	RegStatus string   `json:"regStatus"`
	Type      string   `json:"type"`
}

func NewItem(metadata Metadata, regStatus string, type_ string) Item {
	return Item{
		Metadata:  metadata,
		RegStatus: regStatus,
		Type:      type_,
	}
}

type Metadata struct {
	AuthStatus        string `json:"authStatus"`
	EvalExpiresInDays int    `json:"evalExpiresInDays"`
	EvalUsed          bool   `json:"evalUsed"`
	ExportControl     bool   `json:"exportControl"`
	VirtualAccount    string `json:"virtualAccount"`
}

func NewMetadata(
	authStatus string,
	evalExpiresInDays int,
	evalUsed bool,
	exportControl bool,
	virtualAccount string,
) Metadata {
	return Metadata{
		AuthStatus:        authStatus,
		EvalExpiresInDays: evalExpiresInDays,
		EvalUsed:          evalUsed,
		ExportControl:     exportControl,
		VirtualAccount:    virtualAccount,
	}
}

package smartlicense

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/internal"
)

type SmartLicense struct {
	Items  []Item          `json:"items"`
	Links  internal.Links  `json:"links"`
	Paging internal.Paging `json:"paging"`
}

func NewSmartLicense(items []Item, links internal.Links, paging internal.Paging) SmartLicense {
	return SmartLicense{
		Items:  items,
		Links:  links,
		Paging: paging,
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

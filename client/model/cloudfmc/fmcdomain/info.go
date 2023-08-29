package fmcdomain

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/internal"
)

type Info struct {
	Links  Links  `json:"links"`
	Paging Paging `json:"paging"`
	Items  []Item `json:"items"`
}

func NewInfo(links Links, paging Paging, items []Item) Info {
	return Info{
		Links:  links,
		Paging: paging,
		Items:  items,
	}
}

type Item struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func NewItem(uuid, name, type_ string) Item {
	return Item{
		Uuid: uuid,
		Name: name,
		Type: type_,
	}
}

type Links = internal.Links
type Paging = internal.Paging

var NewLinks = internal.NewLinks
var NewPaging = internal.NewPaging

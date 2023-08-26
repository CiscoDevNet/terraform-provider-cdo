package fmcdomain

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmccommon"
)

type Info struct {
	Links  fmccommon.Links  `json:"links"`
	Paging fmccommon.Paging `json:"paging"`
	Items  []Item           `json:"items"`
}

func NewInfo(links fmccommon.Links, paging fmccommon.Paging, items []Item) Info {
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

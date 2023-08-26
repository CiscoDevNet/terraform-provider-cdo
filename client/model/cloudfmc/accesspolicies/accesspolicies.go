package accesspolicies

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmccommon"
)

type AccessPolicies struct {
	Items  []Item           `json:"items"`
	Links  fmccommon.Links  `json:"links"`
	Paging fmccommon.Paging `json:"paging"`
}

func New(items []Item, links fmccommon.Links, paging fmccommon.Paging) AccessPolicies {
	return AccessPolicies{
		Items:  items,
		Links:  links,
		Paging: paging,
	}
}

// Find return the access policy item with the given name, second return value ok indicate whether the item is found.
func (policies *AccessPolicies) Find(name string) (item Item, ok bool) {
	for _, policy := range policies.Items {
		if policy.Name == name {
			return policy, true
		}
	}
	return Item{}, false
}

type Item struct {
	Links fmccommon.Links `json:"links"`
	Id    string          `json:"id"`
	Name  string          `json:"name"`
	Type  string          `json:"type"`
}

func NewItem(id, name, type_ string, links fmccommon.Links) Item {
	return Item{
		Id:    id,
		Name:  name,
		Type:  type_,
		Links: links,
	}
}

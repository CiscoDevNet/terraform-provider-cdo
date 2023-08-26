package accesspolicies

import "github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc"

type AccessPolicies struct {
	Items  []Item          `json:"items"`
	Links  cloudfmc.Links  `json:"links"`
	Paging cloudfmc.Paging `json:"paging"`
}

func New(items []Item, links cloudfmc.Links, paging cloudfmc.Paging) AccessPolicies {
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
	Links cloudfmc.Links `json:"links"`
	Id    string         `json:"id"`
	Name  string         `json:"name"`
	Type  string         `json:"type"`
}

func NewItem(id, name, type_ string, links cloudfmc.Links) Item {
	return Item{
		Id:    id,
		Name:  name,
		Type:  type_,
		Links: links,
	}
}

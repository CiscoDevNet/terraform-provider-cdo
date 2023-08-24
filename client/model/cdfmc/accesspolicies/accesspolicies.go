package accesspolicies

type AccessPolicies struct {
	Items  []Item `json:"items"`
	Links  Links  `json:"links"`
	Paging Paging `json:"paging"`
}

func New(items []Item, links Links, paging Paging) AccessPolicies {
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
	Links Links  `json:"links"`
	Id    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

func NewItem(id, name, type_ string, links Links) Item {
	return Item{
		Id:    id,
		Name:  name,
		Type:  type_,
		Links: links,
	}
}

type Paging struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Pages  int `json:"pages"`
}

func NewPaging(count, offset, limit, pages int) Paging {
	return Paging{
		Count:  count,
		Offset: offset,
		Limit:  limit,
		Pages:  pages,
	}
}

type Links struct {
	Self string `json:"self"`
}

func NewLinks(self string) Links {
	return Links{
		Self: self,
	}
}

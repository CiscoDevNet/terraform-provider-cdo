package accesspolicies

type AccessPolicies struct {
	Items  Items  `json:"items"`
	Links  Links  `json:"links"`
	Paging Paging `json:"paging"`
}

func New(items Items, links Links, paging Paging) AccessPolicies {
	return AccessPolicies{
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

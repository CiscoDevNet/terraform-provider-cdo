package internal

type Paging struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Pages  int `json:"pages"`
}

func NewPaging(
	count int,
	offset int,
	limit int,
	pages int,
) Paging {
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

type NoLinkItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func NewNoLinkItem(id, name, type_ string) NoLinkItem {
	return NoLinkItem{
		Id:   id,
		Name: name,
		Type: type_,
	}
}

type Response struct {
	Items  []Item `json:"items"`
	Links  Links  `json:"links"`
	Paging Paging `json:"paging"`
}

func NewResponse(items []Item, links Links, paging Paging) Response {
	return Response{
		Items:  items,
		Links:  links,
		Paging: paging,
	}
}

// Find return the item with the given name, second return value ok indicate whether the item is found.
func (response *Response) Find(name string) (Item, bool) {
	for _, item := range response.Items {
		if item.Name == name {
			return item, true
		}
	}
	return Item{}, false
}

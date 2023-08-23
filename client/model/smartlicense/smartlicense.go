package smartlicense

type SmartLicense struct {
	Items  Items  `json:"items"`
	Links  Links  `json:"links"`
	Paging Paging `json:"paging"`
}

func NewSmartLicense(items Items, links Links, paging Paging) SmartLicense {
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

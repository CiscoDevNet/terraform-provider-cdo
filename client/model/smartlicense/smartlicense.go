package smartlicense

type SmartLicense struct {
	Items  Items  `json:"items"`
	Links  Links  `json:"links"`
	Paging Paging `json:"paging"`
}

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	Metadata  Metadata `json:"metadata"`
	RegStatus string   `json:"regStatus"`
	Type      string   `json:"type"`
}

type Metadata struct {
	AuthStatus        string `json:"authStatus"`
	EvalExpiresInDays int    `json:"evalExpiresInDays"`
	EvalUsed          bool   `json:"evalUsed"`
	ExportControl     bool   `json:"exportControl"`
	VirtualAccount    string `json:"virtualAccount"`
}

type Paging struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Pages  int `json:"pages"`
}

type Links struct {
	Self string `json:"self"`
}

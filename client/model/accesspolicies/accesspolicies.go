package accesspolicies

type AccessPolicies struct {
	Items  Items  `json:"items"`
	Links  Links  `json:"links"`
	Paging Paging `json:"paging"`
}

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	Links Links  `json:"links"`
	Id    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
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

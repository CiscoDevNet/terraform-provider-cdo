package model

type UserDetails struct {
	Uid                 string   `json:"uid"`
	Name                string   `json:"name"`
	UserRoles           []string `json:"roles"`
	ApiOnlyUser         bool     `json:"isApiOnlyUser"`
	LastSuccessfulLogin int64    `json:"lastSuccessfulLogin"`
	ApiTokenId          string   `json:"apiTokenId"`
}

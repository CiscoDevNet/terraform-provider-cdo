package auth

type Token struct {
	TenantUid    string `json:"tenantUid"`
	TenantName   string `json:"tenantName"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

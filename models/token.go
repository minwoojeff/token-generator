package models

type Token struct {
	AccessToken  string   `json:"access_token"`
	TokenType    string   `json:"token_type"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	Scope        []string `json:"scope"`
	JTI          string   `json:"jti"`
}

package models

type AuthRequest struct {
	ClientID		string `json:"client_id"`
	ClientSecret	string `json:"client_secret"`
	GrantType    	string `json:"grant_type"`
	Code         	string `json:"code"`
	RedirectURI  	string `json:"redirect_uri"`
}

type AuthResponse struct {
	AccessToken  	string `json:"access_token"`
	RefreshToken 	string `json:"refresh_token"`
	TokenType    	string `json:"token_type"`
	ExpiresIn    	int    `json:"expires_in"`
}

type AccountResponse struct {
	ID int `json:"id"`
}
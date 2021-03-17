package apple

type AppValidationTokenRequest struct {
	ClientID     string
	ClientSecret string
	Code         string
}

type ValidationResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	Error        string `json:"error"`
}

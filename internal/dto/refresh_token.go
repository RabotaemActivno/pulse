package dto

type RefreshTokenInput struct {
	Token string `json:"token"`
}

type RefreshTokenOutput struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

package dto

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserOutput struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

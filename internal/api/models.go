package api

type GenerateTokenRequest struct {
	UserId string `json:"user_id"`
}

type RefreshTokenRequest struct {
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

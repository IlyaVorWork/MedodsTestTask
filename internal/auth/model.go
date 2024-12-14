package auth

type RefreshDTO struct {
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	GUID             string
	Email            string
	RefreshTokenHash *string
}

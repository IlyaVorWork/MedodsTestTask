package token

import "github.com/dgrijalva/jwt-go"

type TokenClaims struct {
	UserGUID string `json:"user_guid"`
	IPv4     string `json:"ipv4"`
	jwt.StandardClaims
}

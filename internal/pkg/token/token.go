package token

import (
	"os"
	"time"

	"MedodsTestTask/internal/pkg"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func GeneratePair(guid, ipv4 string) (map[string]string, error) {
	JWT_SECRET_KEY, _ := os.LookupEnv("JWT_SECRET_KEY")

	claims := jwt.MapClaims{
		"exp":       time.Now().Add(time.Minute * 30).Unix(),
		"jti":       uuid.New().String(),
		"user_guid": guid,
		"ipv4":      ipv4,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	at, err := token.SignedString([]byte(JWT_SECRET_KEY))
	if err != nil {
		return nil, err
	}

	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	rt, err := refresh.SignedString([]byte(JWT_SECRET_KEY))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  at,
		"refresh_token": rt,
	}, err
}

func GetClaims(token string) (*TokenClaims, error) {
	JWT_SECRET_KEY, _ := os.LookupEnv("JWT_SECRET_KEY")

	claims := TokenClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET_KEY), nil
	})
	if err != nil && err.Error() != "" {
		return nil, err
	}

	return &claims, nil
}

func ValidateRefreshToken(rtClaims *TokenClaims, ipv4 string) error {
	expired := rtClaims.VerifyExpiresAt(time.Now().Unix(), true)
	if !expired {
		return pkg.ExpiredTokenError
	}

	if rtClaims.IPv4 != ipv4 {
		return pkg.UnmatchedIPsError
	}

	return nil
}

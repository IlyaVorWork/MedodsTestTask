package pkg

import "errors"

var (
	InvalidTokenError   = errors.New("provided token is invalid")
	ExpiredTokenError   = errors.New("provided token is expired")
	UnmatchedIPsError   = errors.New("IPs does not match")
	UnpairedTokensError = errors.New("provided tokens are not from the same pair")
)
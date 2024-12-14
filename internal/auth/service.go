package auth

import (
	"errors"

	"MedodsTestTask/internal/pkg"
	"MedodsTestTask/internal/pkg/token"
)

type IProvider interface {
	GetUserByGUID(guid string) (*User, error)
	SetUserRefresh(guid, refreshToken string) error
}

type Service struct {
	provider IProvider
}

func NewUserService(provider IProvider) *Service {
	return &Service{
		provider: provider,
	}
}

func (service *Service) Login(guid, ipv4 string) (map[string]string, error) {
	_, err := service.provider.GetUserByGUID(guid)
	if err != nil {
		return nil, err
	}

	pair, err := token.GeneratePair(guid, ipv4)
	if err != nil {
		return nil, err
	}

	err = service.provider.SetUserRefresh(guid, pair["refresh_token"])
	if err != nil {
		return nil, err
	}

	return pair, nil
}

func (service *Service) Refresh(accessToken, refreshToken, ipv4 string) (map[string]string, error) {
	accessClaims, err := token.GetClaims(accessToken)
	if err != nil {
		return nil, err
	}

	refreshClaims, err := token.GetClaims(refreshToken)
	if err != nil {
		return nil, err
	}

	err = token.ValidateRefreshToken(refreshClaims, ipv4)
	if err != nil {
		if errors.Is(err, pkg.UnmatchedIPsError) {
			// Код отправки warning на email пользователя
		}
		suberr := service.provider.SetUserRefresh(refreshClaims.UserGUID, "")

		if suberr != nil {
			return nil, suberr
		}
		return nil, err
	}

	if accessClaims.Id != refreshClaims.Id {
		suberr := service.provider.SetUserRefresh(refreshClaims.UserGUID, "")

		if suberr != nil {
			return nil, suberr
		}
		return nil, pkg.UnpairedTokensError
	}

	user, err := service.provider.GetUserByGUID(refreshClaims.UserGUID)
	if err != nil {
		return nil, err
	}

	if user.RefreshTokenHash == nil || *user.RefreshTokenHash != refreshToken {
		suberr := service.provider.SetUserRefresh(refreshClaims.UserGUID, "")

		if suberr != nil {
			return nil, suberr
		}
		return nil, pkg.InvalidTokenError
	}

	pair, err := token.GeneratePair(refreshClaims.UserGUID, ipv4)
	if err != nil {
		return nil, err
	}

	err = service.provider.SetUserRefresh(refreshClaims.UserGUID, pair["refresh_token"])
	if err != nil {
		return nil, err
	}

	return pair, nil
}

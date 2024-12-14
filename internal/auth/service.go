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
	// Поиск пользователя по ID
	_, err := service.provider.GetUserByGUID(guid)
	if err != nil {
		return nil, pkg.UnexistingUserError
	}

	// Генерация пары access и refresh токенов
	pair, err := token.GeneratePair(guid, ipv4)
	if err != nil {
		return nil, err
	}

	// Установка refresh token'а в БД соответствующему пользователю
	err = service.provider.SetUserRefresh(guid, pair["refresh_token"])
	if err != nil {
		return nil, err
	}

	return pair, nil
}

func (service *Service) Refresh(accessToken, refreshToken, ipv4 string) (map[string]string, error) {

	// Получение claims access token'а
	accessClaims, err := token.GetClaims(accessToken)
	if err != nil {
		return nil, err
	}

	// Получение claims refresh token'а
	refreshClaims, err := token.GetClaims(refreshToken)
	if err != nil {
		return nil, err
	}

	// Валидация refresh токена
	err = token.ValidateRefreshToken(refreshClaims, ipv4)
	if err != nil {
		// Если IP не совпадают
		if errors.Is(err, pkg.UnmatchedIPsError) {
			// TODO:Код отправки warning на email пользователя
		}
		// Обнуление refresh token'а в БД
		suberr := service.provider.SetUserRefresh(refreshClaims.UserGUID, "")

		if suberr != nil {
			return nil, suberr
		}
		return nil, err
	}

	// Проверка токенов на принадлежность к одной паре
	if accessClaims.Id != refreshClaims.Id {
		// Обнуление refresh token'а в БД
		suberr := service.provider.SetUserRefresh(refreshClaims.UserGUID, "")

		if suberr != nil {
			return nil, suberr
		}
		return nil, pkg.UnpairedTokensError
	}

	// Получение пользователя по ID из refresh token'а
	user, err := service.provider.GetUserByGUID(refreshClaims.UserGUID)
	if err != nil {
		return nil, err
	}

	// Проверка на совпадение предоставленного refresh token'а с находящимся в БД
	if user.RefreshTokenHash == nil || *user.RefreshTokenHash != refreshToken {
		// Обнуление refresh token'а в БД
		suberr := service.provider.SetUserRefresh(refreshClaims.UserGUID, "")

		if suberr != nil {
			return nil, suberr
		}
		return nil, pkg.InvalidTokenError
	}

	// Генерация пары access и refresh токенов
	pair, err := token.GeneratePair(refreshClaims.UserGUID, ipv4)
	if err != nil {
		return nil, err
	}

	// Установка refresh token'а в БД соответствующему пользователю
	err = service.provider.SetUserRefresh(refreshClaims.UserGUID, pair["refresh_token"])
	if err != nil {
		return nil, err
	}

	return pair, nil
}

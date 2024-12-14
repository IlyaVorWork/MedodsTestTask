package auth

import "database/sql"

type Provider struct {
	db *sql.DB
}

func NewProvider(db *sql.DB) *Provider {
	return &Provider{db: db}
}

func (provider *Provider) GetUserByGUID(guid string) (*User, error) {
	var user User

	err := provider.db.QueryRow("SELECT id, email, refresh_token FROM public.\"User\" WHERE id = $1", guid).Scan(&user.GUID, &user.Email, &user.RefreshTokenHash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (provider *Provider) SetUserRefresh(guid, refreshToken string) error {
	_, err := provider.db.Exec("UPDATE public.\"User\" SET refresh_token = $1 WHERE id = $2", refreshToken, guid)
	if err != nil {
		return err
	}

	return nil
}

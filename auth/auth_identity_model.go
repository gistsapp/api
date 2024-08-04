package auth

import (
	"database/sql"
	"errors"

	"github.com/gistapp/api/storage"
	"github.com/gistapp/api/user"
	"github.com/gofiber/fiber/v2/log"
)

type AuthIdentitySQL struct {
	ID         sql.NullInt32
	Data       sql.NullString
	Type       sql.NullString
	ProviderID sql.NullString
	OwnerID    sql.NullString
}

type AuthIdentity struct {
	ID         string `json:"id"`
	Data       string `json:"data"`
	Owner      string `json:"owner"`
	Type       string `json:"type"`
	ProviderID string `json:"provider_id"`
}

type AuthIdentityAndUser struct {
	AuthIdentity AuthIdentity
	User         user.User
}

type AuthIdentityModel interface {
	Save() (*AuthIdentity, error)
}

func (ai *AuthIdentitySQL) Save() (*AuthIdentity, error) {
	row, err := storage.Database.Query("INSERT INTO auth_identity(data, type, owner_id, provider_id) VALUES ($1, $2, $3, $4) RETURNING auth_id, data, owner_id, type, provider_id", ai.Data.String, ai.Type.String, ai.OwnerID.String, ai.ProviderID.String)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't create auth identity")
	}
	var authIdentity AuthIdentity
	row.Next()
	err = row.Scan(&authIdentity.ID, &authIdentity.Data, &authIdentity.Owner, &authIdentity.Type, &authIdentity.ProviderID)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't find auth identity")
	}
	return &authIdentity, nil
}

func (ai *AuthIdentitySQL) GetWithUser(provider_id string) (*AuthIdentityAndUser, error) {
	query := "SELECT a.auth_id, a.data, a.owner_id, a.type, a.provider_id, u.user_id, u.email, u.name, u.picture FROM auth_identity a JOIN users u ON a.owner_id = u.user_id WHERE a.provider_id = $1"

	row, err := storage.Database.Query(query, provider_id)

	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't get auth identity")
	}

	row.Next()
	var authIdentity AuthIdentityAndUser
	err = row.Scan(&authIdentity.AuthIdentity.ID, &authIdentity.AuthIdentity.Data, &authIdentity.AuthIdentity.Owner, &authIdentity.AuthIdentity.Type, &authIdentity.AuthIdentity.ProviderID, &authIdentity.User.ID, &authIdentity.User.Email, &authIdentity.User.Name, &authIdentity.User.Picture)

	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't find auth identity")
	}

	return &authIdentity, nil

}

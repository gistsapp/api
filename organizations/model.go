package organizations

import (
	"database/sql"
	"errors"

	"github.com/gistapp/api/storage"
	"github.com/gofiber/fiber/v2/log"
)

type Role string

const (
	Owner Role = "owner"
)

type OrganizationSQL struct {
	ID   sql.NullInt32
	Name sql.NullString
}

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type OrganizationModel interface {
	Save() error
}

func (o *OrganizationSQL) Save(owner_id string) (*Organization, error) {
	row, err := storage.Database.Query("INSERT INTO organization(name) VALUES ($1) RETURNING org_id, name", o.Name.String)

	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't create organization")
	}

	var organization Organization

	row.Next()
	err = row.Scan(&organization.ID, &organization.Name)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't find organization")
	}

	_, err = storage.Database.Exec("INSERT INTO member(org_id, user_id, role) VALUES ($1, $2, $3)", organization.ID, owner_id, Owner)

	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't create member")
	}
	return &organization, nil
}

func (o *OrganizationSQL) Delete() error {
	_, err := storage.Database.Exec("DELETE FROM organization WHERE org_id = $1", o.ID.Int32)
	if err != nil {
		log.Error(err)
		return errors.New("couldn't delete organization")
	}
	return nil
}

package organizations

import (
	"database/sql"
	"errors"

	"github.com/gistapp/api/storage"
	"github.com/gofiber/fiber/v2/log"
)

type OrganizationSQL struct {
	ID sql.NullInt32
	Name sql.NullString
}

type Organization struct {
	ID string `json:"id"`
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

	_, err = storage.Database.Exec("INSERT INTO member(org_id, member_id) VALUES ($1, $2)", organization.ID, owner_id)

	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't create member")
	}
	return &organization, nil
}
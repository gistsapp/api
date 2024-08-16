package organizations

import (
	"database/sql"
	"errors"
)

type OrganizationServiceImpl struct{}

func (g *OrganizationServiceImpl) Save(name string, owner_id string) (*Organization, error) {
	m := OrganizationSQL{
		ID: sql.NullInt32{
			Valid: false,
			Int32: 0,
		}, // useless ID
		Name: sql.NullString{
			String: name,
			Valid:  true,
		},
	}

	organization, err := m.Save(owner_id)
	if err != nil {
		return nil, errors.New("couldn't insert into database organization")
	}
	return organization, nil
}

var OrganizationService *OrganizationServiceImpl = &OrganizationServiceImpl{}

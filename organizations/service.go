package organizations

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2/log"
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

	log.Info("saving organization")

	organization, err := m.Save(owner_id)
	if err != nil {
		return nil, errors.New("couldn't insert into database organization")
	}
	return organization, nil
}

// returns a list of organizations that the user is a member of
func (g *OrganizationServiceImpl) GetAsMember(user_id string) ([]Organization, error) {
	m := OrganizationSQL{
		ID: sql.NullInt32{
			Valid: false,
			Int32: 0,
		},
		Name: sql.NullString{
			String: "",
			Valid:  false,
		},
	}
	organizations, err := m.GetByMember(user_id)
	if err != nil {
		return nil, errors.New("couldn't find organizations")
	}
	return organizations, nil
}

func (g *OrganizationServiceImpl) GetByID(org_id string, user_id string) (*Organization, error) {
	m := OrganizationSQL{
		ID: sql.NullInt32{
			Valid: true,
			Int32: 0,
		},
		Name: sql.NullString{
			String: "",
			Valid:  false,
		},
	}
	organization, err := m.GetByID(user_id, org_id)
	if err != nil {
		return nil, errors.New("couldn't find organization")
	}
	return organization, nil
}

var OrganizationService *OrganizationServiceImpl = &OrganizationServiceImpl{}

package organizations

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2/log"
)

type OrganizationServiceImpl struct{}

func (g *OrganizationServiceImpl) Save(name string, owner_id string) (*Organization, error) {
	m := OrganizationSQL{
		ID: sql.NullString{
			Valid:  false,
			String: "",
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

func (g *OrganizationServiceImpl) IsOwner(org_id string, user_id string) bool {
	member_sql := MemberSQL{
		OrgID: sql.NullString{
			String: org_id,
			Valid:  true,
		},
		UserID: sql.NullString{
			String: user_id,
			Valid:  true,
		},
	}

	member, err := member_sql.Get()
	if err != nil {
		return false
	}

	return member.Role == Owner
}

// returns a list of organizations that the user is a member of
func (g *OrganizationServiceImpl) GetAsMember(user_id string) ([]Organization, error) {
	m := OrganizationSQL{
		ID: sql.NullString{
			Valid:  false,
			String: "",
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
		ID: sql.NullString{
			Valid:  true,
			String: "",
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

func (g *OrganizationServiceImpl) Delete(org_id string, user_id string) error {

	if !g.IsOwner(org_id, user_id) {
		return errors.New("user is not the owner of the organization")
	}

	m := OrganizationSQL{
		ID: sql.NullString{
			Valid:  true,
			String: org_id,
		},
		Name: sql.NullString{
			String: "",
			Valid:  false,
		},
	}
	log.Info("org_id: %s\n", org_id)
	err := m.Delete()
	if err != nil {
		log.Error(err)
		return errors.New("couldn't delete organization")
	}
	return nil
}

var OrganizationService *OrganizationServiceImpl = &OrganizationServiceImpl{}
var ErrUserNotOwner = errors.New("user is not the owner of the organization")

package organizations

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

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

func (g *OrganizationServiceImpl) IsOwner(org_id int, user_id int) bool {
	member_sql := MemberSQL{
		OrgID: sql.NullInt32{
			Int32: int32(org_id),
			Valid: true,
		},
		UserID: sql.NullInt32{
			Int32: int32(user_id),
			Valid: true,
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

func (g *OrganizationServiceImpl) Delete(org_id string, user_id string) error {
	i_org_id, err := strconv.Atoi(org_id)

	if err != nil {
		return errors.New("organization ID must be an integer")
	}

	i_user_id, err := strconv.Atoi(user_id)

	if err != nil {
		return errors.New("user ID must be an integer")
	}

	if !g.IsOwner(i_org_id, i_user_id) {
		return errors.New("user is not the owner of the organization")
	}

	m := OrganizationSQL{
		ID: sql.NullInt32{
			Valid: true,
			Int32: int32(i_org_id),
		},
		Name: sql.NullString{
			String: "",
			Valid:  false,
		},
	}
	fmt.Printf("org_id: %d\n", i_org_id)
	err = m.Delete()
	if err != nil {
		fmt.Println(err)
		return errors.New("couldn't delete organization")
	}
	return nil
}

var OrganizationService *OrganizationServiceImpl = &OrganizationServiceImpl{}
var ErrUserNotOwner = errors.New("user is not the owner of the organization")

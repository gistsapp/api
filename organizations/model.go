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
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Gists []string `json:"gists,omitempty"`
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

func (o *OrganizationSQL) GetByMember(user_id string) ([]Organization, error) {
	query := "SELECT o.org_id, o.name FROM organization o JOIN member ON o.org_id = member.org_id WHERE user_id = $1"
	rows, err := storage.Database.Query(query, user_id)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't find organizations")
	}
	var organizations []Organization
	for rows.Next() {
		var organization Organization
		err = rows.Scan(&organization.ID, &organization.Name)
		if err != nil {
			log.Error(err)
			return nil, errors.New("couldn't find organization")
		}
		organizations = append(organizations, organization)
	}
	return organizations, nil
}

func (o *OrganizationSQL) GetByID(user_id string, org_id string) (*Organization, error) {
	query := "SELECT o.org_id, o.name, gists.gist_id FROM organization o JOIN gists ON o.org_id = gists.org_id WHERE o.org_id=$1 AND gists.owner=$2"
	log.Info("user is ", user_id)

	rows, err := storage.Database.Query(query, org_id, user_id)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't find organization")
	}

	type organization_line struct {
		ID   string
		Name string
		Gist string
	}

	orgs := []organization_line{}

	for rows.Next() {
		organization := new(organization_line)
		err = rows.Scan(&organization.ID, &organization.Name, &organization.Gist)
		if err != nil {
			log.Error(err)
			return nil, errors.New("couldn't find organization")
		}

		orgs = append(orgs, *organization)
	}

	if len(orgs) == 0 {
		return nil, errors.New("couldn't find organization")
	}

	gists_ids := []string{}

	for _, org := range orgs {
		gists_ids = append(gists_ids, org.Gist)
	}

	return &Organization{
		Name:  orgs[0].Name,
		ID:    orgs[0].ID,
		Gists: gists_ids,
	}, nil
}

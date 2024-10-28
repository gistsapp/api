package gists

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
)

type GistServiceImpl struct{}

func (g *GistServiceImpl) Save(name string, content string, ownerID string, orgID string, language string, description string) (*Gist, error) {
	// Helper function to set NullString type based on value
	toNullString := func(s string) sql.NullString {
		return sql.NullString{
			String: s,
			Valid:  true,
		}
	}

	// Helper function to set NullInt32 type based on value
	toNullInt32 := func(s string) (sql.NullInt32, error) {
		if s == "" {
			return sql.NullInt32{Valid: false}, nil
		}
		intValue, err := strconv.Atoi(s)
		if err != nil {
			return sql.NullInt32{}, errors.New("org_id must be an integer")
		}
		return sql.NullInt32{
			Int32: int32(intValue),
			Valid: true,
		}, nil
	}

	// Set values for GistSQL fields
	orgIDNullInt32, err := toNullInt32(orgID)
	if err != nil {
		return nil, err
	}

	if language == "" {
		language = "text"
	}

	m := GistSQL{
		ID: sql.NullInt32{
			Valid: false, // Assuming ID is auto-generated and not required here
			Int32: 0,
		},
		Name:        toNullString(name),
		Content:     toNullString(content),
		OwnerID:     toNullString(ownerID),
		Language:    toNullString(language),
		Description: toNullString(description),
		OrgID:       orgIDNullInt32,
	}

	// Save and handle errors
	gist, err := m.Save()
	if err != nil {
		return nil, errors.New("couldn't insert into database gists")
	}
	return gist, nil
}

func (g *GistServiceImpl) UpdateName(id string, name string, owner_id string) (*Gist, error) {

	f := gistExists(id, owner_id)

	if f != nil {
		return nil, f
	}

	m := GistSQL{
		ID: sql.NullInt32{
			Valid: true,
			Int32: 0,
		},
		Name: sql.NullString{
			String: name,
			Valid:  true,
		},
		Content: sql.NullString{
			String: "",
			Valid:  false,
		},

		OwnerID: sql.NullString{
			String: owner_id,
			Valid:  true,
		},
	}
	gist, err := m.UpdateName(id)
	if err != nil {
		return nil, errors.New("couldn't update name in database gists")
	}
	return gist, nil
}

func (g *GistServiceImpl) UpdateContent(id string, content string, owner_id string) (*Gist, error) {
	err := gistExists(id, owner_id)

	if err != nil {
		return nil, err
	}
	m := GistSQL{
		ID: sql.NullInt32{
			Valid: true,
			Int32: 0,
		},
		Name: sql.NullString{
			String: "",
			Valid:  false,
		},
		Content: sql.NullString{
			String: content,
			Valid:  true,
		},

		OwnerID: sql.NullString{
			String: owner_id,
			Valid:  true,
		},
	}

	gist, err := m.UpdateContent(id)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't update content in database gists")
	}
	return gist, nil
}

func (g *GistServiceImpl) UpdateDescription(id string, description string, owner_id string) (*Gist, error) {
	err := gistExists(id, owner_id)
	if err != nil {
		return nil, err
	}
	m := GistSQL{
		ID: sql.NullInt32{
			Valid: true,
			Int32: 0,
		},
		Name: sql.NullString{
			String: "",
			Valid:  false,
		},
		Content: sql.NullString{
			String: "",
			Valid:  false,
		},
		Description: sql.NullString{
			String: description,
			Valid:  true,
		},
		OwnerID: sql.NullString{
			String: owner_id,
			Valid:  true,
		},
	}
	gist, err := m.UpdateField(id, "description", description)
	if err != nil {
		return nil, errors.New("couldn't update description in database gists")
	}
	return gist, nil
}

func (g *GistServiceImpl) UpdateLanguage(id string, language string, owner_id string) (*Gist, error) {
	err := gistExists(id, owner_id)
	if err != nil {
		return nil, err
	}
	m := GistSQL{
		ID: sql.NullInt32{
			Valid: true,
			Int32: 0,
		},
		Name: sql.NullString{
			String: "",
			Valid:  false,
		},
		Content: sql.NullString{
			String: "",
			Valid:  false,
		},
		Language: sql.NullString{
			String: language,
			Valid:  true,
		},
		OwnerID: sql.NullString{
			String: owner_id,
			Valid:  true,
		},
	}
	return m.UpdateField(id, "language", language)
}

func (g *GistServiceImpl) Delete(id string, owner_id string) error {
	err := gistExists(id, owner_id)

	if err != nil {
		return err
	}
	m := GistSQL{
		ID: sql.NullInt32{
			Valid: true,
			Int32: 0,
		},
		Name: sql.NullString{
			String: "",
			Valid:  false,
		},
		Content: sql.NullString{
			String: "",
			Valid:  false,
		},

		OwnerID: sql.NullString{
			String: owner_id,
			Valid:  true,
		},
	}
	err = m.Delete(id)
	if err != nil {
		return errors.New("couldn't delete from database gists")
	}
	return nil
}

func (g *GistServiceImpl) FindAll(owner_id string) ([]Gist, error) {
	m := GistSQL{

		OwnerID: sql.NullString{
			String: owner_id,
			Valid:  true,
		},
	}
	gists, err := m.FindAll()
	if err != nil {
		return nil, errors.New("couldn't get gists")
	}
	return gists, nil
}

func (g *GistServiceImpl) FindByID(id string, owner_id string) (*Gist, error) {
	m := GistSQL{

		OwnerID: sql.NullString{
			String: owner_id,
			Valid:  true,
		},
	}
	gist, err := m.FindByID(id)
	if err != nil {
		return nil, errors.New("couldn't get gist")
	}
	return gist, nil
}

func gistExists(id string, owner_id string) error {
	m := GistSQL{

		OwnerID: sql.NullString{
			String: owner_id,
			Valid:  true,
		},
	}

	gists, err := m.FindByID(id)

	if err != nil {
		return ErrGistNotFound
	}

	if gists == nil {
		return ErrGistNotFound
	}
	return nil
}

var ErrGistNotFound = errors.New("gist not found")

var GistService GistServiceImpl = GistServiceImpl{}

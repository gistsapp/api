package gists

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2/log"
)

type GistServiceImpl struct{}

func (g *GistServiceImpl) Save(name string, content string, ownerID string, orgID string, language string, description string, visibility string) (*Gist, error) {
	// Helper function to set NullString type based on value
	toNullString := func(s string) sql.NullString {
		return sql.NullString{
			String: s,
			Valid:  true,
		}
	}

	toReallyNullString := func(s string) sql.NullString {
		if s == "" {
			return sql.NullString{
				String: "",
				Valid:  false,
			}
		}
		return sql.NullString{
			String: s,
			Valid:  true,
		}
	}

	if language == "" {
		language = "text"
	}

	m := GistSQL{
		ID: sql.NullString{
			Valid:  false, // Assuming ID is auto-generated and not required here
			String: "",
		},
		Name:        toNullString(name),
		Content:     toNullString(content),
		OwnerID:     toNullString(ownerID),
		Language:    toNullString(language),
		Description: toNullString(description),
		OrgID:       toReallyNullString(orgID),
		Visibility:  toNullString(visibility),
	}

	// Save and handle errors
	gist, err := m.Save()
	if err != nil {
		return nil, errors.New("couldn't insert into database gists")
	}

	rights := GistRights{
		UserID: ownerID,
		GistID: gist.ID,
		Right:  string(Write),
	}

	_, err = rights.Save()

	return gist, err
}

func (g *GistServiceImpl) UpdateName(id string, name string, owner_id string) (*Gist, error) {

	f := gistExists(id, owner_id)

	if f != nil {
		return nil, f
	}

	m := GistSQL{
		ID: sql.NullString{
			Valid:  true,
			String: "",
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
		ID: sql.NullString{
			Valid:  true,
			String: "",
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
		ID: sql.NullString{
			Valid:  true,
			String: "",
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
		ID: sql.NullString{
			Valid:  true,
			String: "",
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
		ID: sql.NullString{
			Valid:  true,
			String: "",
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

func (g *GistServiceImpl) FindAll(owner_id string, limit int, offset int) ([]Gist, error) {
	m := GistSQL{

		OwnerID: sql.NullString{
			String: owner_id,
			Valid:  true,
		},
	}
	gists, err := m.FindAll(limit, offset)
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

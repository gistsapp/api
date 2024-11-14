package gists

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2/log"
)

type GistServiceImpl struct{}

func (g *GistServiceImpl) Save(name string, content string, ownerID string, orgID sql.NullString, language string, description string, visibility string) (*Gist, error) {
	// Helper function to set NullString type based on value

	m := GistSQL{
		ID:          "",
		Name:        name,
		Content:     content,
		OwnerID:     ownerID,
		Language:    language,
		Description: description,
		OrgID:       orgID,
		Visibility:  visibility,
	}

	// Save and handle errors
	gist, err := m.Save()
	if err != nil {
		log.Error(err)
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

	m := NewGistSQL(id, name, "", owner_id, sql.NullString{}, "", "", "")

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
	m := NewGistSQL(id, "", content, owner_id, sql.NullString{}, "", "", "")
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
	m := NewGistSQL(id, "", "", owner_id, sql.NullString{}, "", "", "")
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
	m := NewGistSQL(id, "", "", owner_id, sql.NullString{}, "", "", "")
	return m.UpdateField(id, "language", language)
}

func (g *GistServiceImpl) Delete(id string, owner_id string) error {
	err := gistExists(id, owner_id)

	if err != nil {
		return err
	}
	m := NewGistSQL(id, "", "", owner_id, sql.NullString{}, "", "", "")
	err = m.Delete(id)
	if err != nil {
		return errors.New("couldn't delete from database gists")
	}
	return nil
}

func (g *GistServiceImpl) FindAll(owner_id string, limit int, offset int, short bool) ([]Gist, error) {
	m := NewGistSQL("", "", "", owner_id, sql.NullString{}, "", "", "")
	gists, err := m.FindAll(limit, offset)
	if err != nil {
		return nil, errors.New("couldn't get gists")
	}
	THRESHOLD := 700
	if short {
		for gist := range gists {
			if len(gists[gist].Content) > THRESHOLD {
				gists[gist].Content = strings.Join(strings.Split(gists[gist].Content, "")[:THRESHOLD], "")
			}
		}
	}
	return gists, nil
}

func (g *GistServiceImpl) FindByID(id string, owner_id string) (*Gist, error) {
	m := NewGistSQL(id, "", "", owner_id, sql.NullString{}, "", "", "")
	gist, err := m.FindByID(id)
	if err != nil {
		return nil, errors.New("couldn't get gist")
	}
	return gist, nil
}

func (g *GistServiceImpl) GetPageCount(owner_id string, limit int) (int, error) {
	m := GistSQL{
		OwnerID: owner_id,
	}
	nb_gists, err := m.Count()
	if err != nil {
		return 0, errors.New("couldn't get gists count")
	}
	nb_pages := int(nb_gists / limit)
	return nb_pages, nil
}

func (g *GistServiceImpl) Update(id string, name string, org_id utils.ZeroString, content string, language string, description string, visibility string, owner_id string) (*Gist, error) {
	err := gistExists(id, owner_id)
	if err != nil {
		return nil, err
	}
	m := NewGistSQL(id, name, content, owner_id, org_id.SqlString(), description, language, visibility)
	gist, err := m.Update()
	if err != nil {
		return nil, errors.New("couldn't update gist")
	}
	return gist, nil
}

func gistExists(id string, owner_id string) error {
	m := NewGistSQL(id, "", "", owner_id, sql.NullString{}, "", "", "")

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

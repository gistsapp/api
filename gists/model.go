package gists

import (
	"database/sql"
	"errors"

	"github.com/gistapp/api/storage"
	"github.com/gistsapp/pogo/pogo"
	"github.com/gofiber/fiber/v2/log"
)

type GistVisibility string

const (
	Private GistVisibility = "private"
	Public  GistVisibility = "public"
)

type GistSQL struct {
	ID          sql.NullString
	Name        sql.NullString
	Content     sql.NullString
	OwnerID     sql.NullString
	OrgID       sql.NullString
	Description sql.NullString
	Language    sql.NullString
	Visibility  sql.NullString
}

type Gist struct {
	ID          string  `json:"id" pogo:"gist_id"`
	Name        string  `json:"name" pogo:"name"`
	Content     string  `json:"content" pogo:"content"`
	OwnerID     string  `json:"owner_id" pogo:"owner"`
	OrgID       *string `json:"org_id,omitempty" pogo:"org_id"`
	Description string  `json:"description" pogo:"description"`
	Language    string  `json:"language" pogo:"language"`
	Visibility  string  `json:"visibility" pogo:"visibility"`
}

type GistModel interface {
	Save() error
}

func (g *GistSQL) Save() (*Gist, error) {
	db := storage.PogoDatabase
	gists := make([]Gist, 0)
	err := pogo.SuperQuery(db, "INSERT INTO gists(name, content, owner, org_id, language, description, visibility) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING :fields", &gists, g.Name.String, g.Content.String, g.OwnerID.String, g.OrgID, g.Language, g.Description, g.Visibility.String)

	if len(gists) <= 0 {
		return nil, errors.New("couldn't create gist")
	}
	return &gists[0], err
}

func (g *GistSQL) UpdateName(id string) (*Gist, error) {
	return g.UpdateField(id, "name", g.Name.String)
}

func (g *GistSQL) UpdateContent(id string) (*Gist, error) {
	db := storage.PogoDatabase

	gists := make([]Gist, 0)
	err := pogo.SuperQuery(db, "UPDATE gists SET content = $1 WHERE gist_id = $2 AND owner = $3 RETURNING :fields", &gists, g.Content.String, id, g.OwnerID.String)

	if len(gists) <= 0 {
		return nil, errors.New("gist not found")
	}

	return &gists[0], err
}

func (g *GistSQL) UpdateField(id string, field string, val string) (*Gist, error) {
	db := storage.PogoDatabase
	gists := make([]Gist, 0)
	err := pogo.SuperQuery(db, "UPDATE gists SET "+field+" = $1 WHERE gist_id = $2 AND owner = $3 RETURNING :fields", &gists, val, id, g.OwnerID.String)
	if len(gists) <= 0 {
		return nil, errors.New("gist not found")
	}
	return &gists[0], err
}

func (g *GistSQL) UpdateVisibility(id string, visibility string) (*Gist, error) {
	return g.UpdateField(id, "visibility", visibility)
}

func (g *GistSQL) Delete(id string) error {
	_, err := storage.Database.Exec("DELETE FROM gists WHERE gist_id = $1 AND owner = $2", id, g.OwnerID.String)
	if err != nil {
		log.Error(err)
		return errors.New("couldn't delete gist")
	}
	return nil
}

func (g *GistSQL) FindByID(id string) (*Gist, error) {
	db := storage.PogoDatabase

	gists := make([]Gist, 0)
	err := pogo.SuperQuery(db, "SELECT :fields FROM gists WHERE gist_id = $1", &gists, id)
	if len(gists) <= 0 {
		log.Error(err)
		return nil, errors.New("gist not found")
	}
	return &gists[0], err
}

func (g *GistSQL) FindAll() ([]Gist, error) {
	db := storage.PogoDatabase

	gists := make([]Gist, 0)
	err := pogo.SuperQuery(db, "SELECT :fields FROM gists WHERE owner = $1", &gists, g.OwnerID.String)
	if len(gists) <= 0 {
		log.Error(err)
		return nil, errors.New("gist not found")
	}
	return gists, err
}

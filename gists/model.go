package gists

import (
	"database/sql"
	"errors"

	"github.com/gistapp/api/storage"
	"github.com/gistapp/api/utils"
	"github.com/gistsapp/pogo/pogo"
	"github.com/gofiber/fiber/v2/log"
)

type GistVisibility string

const (
	Private GistVisibility = "private"
	Public  GistVisibility = "public"
)

type GistSQL struct {
	ID          string
	Name        string
	Content     string
	OwnerID     string
	OrgID       sql.NullString
	Description string
	Language    string
	Visibility  string
}

type Gist struct {
	ID          string         `json:"id" pogo:"gist_id"`
	Name        string         `json:"name" pogo:"name"`
	Content     string         `json:"content" pogo:"content"`
	OwnerID     string         `json:"owner_id" pogo:"owner"`
	OrgID       sql.NullString `json:"org_id,omitempty" pogo:"org_id"`
	Description string         `json:"description" pogo:"description"`
	Language    string         `json:"language" pogo:"language"`
	Visibility  string         `json:"visibility" pogo:"visibility"`
}

type GistModel interface {
	Save() error
}

func (g *Gist) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":          g.ID,
		"name":        g.Name,
		"content":     g.Content,
		"owner_id":    g.OwnerID,
		"org_id":      utils.ToNullString(g.OrgID),
		"description": g.Description,
		"language":    g.Language,
		"visibility":  g.Visibility,
	}
}

func NewGistSQL(ID string, Name string, Content string, OwnerID string, OrgID sql.NullString, Description string, Language string, Visibility string) *GistSQL {
	return &GistSQL{
		ID:          ID,
		Name:        Name,
		Content:     Content,
		OwnerID:     OwnerID,
		OrgID:       OrgID,
		Description: Description,
		Language:    Language,
		Visibility:  Visibility,
	}
}

func (g *GistSQL) Save() (*Gist, error) {
	db := storage.PogoDatabase
	gists := make([]Gist, 0)
	err := pogo.SuperQuery(db, "INSERT INTO gists(name, content, owner, org_id, language, description, visibility) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING :fields", &gists, g.Name, g.Content, g.OwnerID, g.OrgID, g.Language, g.Description, g.Visibility)

	if len(gists) <= 0 {
		log.Error(err)
		return nil, errors.New("couldn't create gist")
	}
	return &gists[0], err
}

func (g *GistSQL) UpdateName(id string) (*Gist, error) {
	return g.UpdateField(id, "name", g.Name)
}

func (g *GistSQL) UpdateContent(id string) (*Gist, error) {
	db := storage.PogoDatabase

	gists := make([]Gist, 0)
	err := pogo.SuperQuery(db, "UPDATE gists SET content = $1 WHERE gist_id = $2 AND owner = $3 RETURNING :fields", &gists, g.Content, id, g.OwnerID)

	if len(gists) <= 0 {
		return nil, errors.New("gist not found")
	}

	return &gists[0], err
}

func (g *GistSQL) UpdateField(id string, field string, val string) (*Gist, error) {
	db := storage.PogoDatabase
	gists := make([]Gist, 0)
	err := pogo.SuperQuery(db, "UPDATE gists SET "+field+" = $1 WHERE gist_id = $2 AND owner = $3 RETURNING :fields", &gists, val, id, g.OwnerID)
	if len(gists) <= 0 {
		return nil, errors.New("gist not found")
	}
	return &gists[0], err
}

func (g *GistSQL) UpdateVisibility(id string, visibility string) (*Gist, error) {
	return g.UpdateField(id, "visibility", visibility)
}

func (g *GistSQL) UpdateGist() (*Gist, error) {
	db := storage.PogoDatabase
	gists := make([]Gist, 0)
	err := pogo.SuperQuery(db, "UPDATE gists SET name = $1, content = $2, language = $3, description = $4, visibility = $5 WHERE gist_id = $6 AND owner = $7 RETURNING :fields", &gists, g.Name, g.Content, g.Language, g.Description, g.Visibility, g.ID, g.OwnerID)
	if len(gists) <= 0 {
		return nil, errors.New("gist not found")
	}
	return &gists[0], err
}

func (g *GistSQL) Delete(id string) error {
	_, err := storage.Database.Exec("DELETE FROM gists WHERE gist_id = $1 AND owner = $2", id, g.OwnerID)
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

func (g *GistSQL) FindAll(limit int, offset int) ([]Gist, error) {
	db := storage.PogoDatabase

	gists := make([]Gist, 0)
	err := pogo.SuperQuery(db, "SELECT :fields FROM gists WHERE owner = $1 LIMIT $2 OFFSET $3", &gists, g.OwnerID, limit, offset)
	if len(gists) <= 0 {
		log.Error(err)
		return nil, errors.New("gist not found")
	}
	return gists, err
}

func (g *GistSQL) Count() (int, error) {
	db := storage.PogoDatabase
	var count int
	rows, err := db.Query("SELECT COUNT(*) FROM gists WHERE owner = $1", g.OwnerID.String)

	rows.Next()

	rows.Scan(&count)

	if err != nil {
		log.Error(err)
		return 0, errors.New("couldn't get gists")
	}
	return count, nil
}

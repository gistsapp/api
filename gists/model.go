package gists

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gistapp/api/storage"
	"github.com/gofiber/fiber/v2/log"
)

type GistSQL struct {
	ID      sql.NullInt32
	Name    sql.NullString
	Content sql.NullString
	OwnerID sql.NullString
	OrgID   sql.NullInt32
}

type Gist struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Content string  `json:"content"`
	OwnerID string  `json:"owner_id"`
	OrgID   *string `json:"org_id,omitempty"`
}

type GistModel interface {
	Save() error
}

func (g *GistSQL) Save() (*Gist, error) {
	var row *sql.Rows
	var err error

	if g.OrgID.Valid {
		row, err = storage.Database.Query("INSERT INTO gists(name, content, owner, org_id) VALUES ($1, $2, $3, $4) RETURNING gist_id, name, content, owner, org_id", g.Name.String, g.Content.String, g.OwnerID.String, g.OrgID.Int32)
	} else {
		row, err = storage.Database.Query("INSERT INTO gists(name, content, owner) VALUES ($1, $2, $3) RETURNING gist_id, name, content, owner", g.Name.String, g.Content.String, g.OwnerID.String)
	}

	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't create gist")
	}

	var gist Gist

	row.Next()
	if g.OrgID.Valid {
		err = row.Scan(&gist.ID, &gist.Name, &gist.Content, &gist.OwnerID, &gist.OrgID)
	} else {
		err = row.Scan(&gist.ID, &gist.Name, &gist.Content, &gist.OwnerID)
		gist.OrgID = nil
	}
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't find gist")
	}
	return &gist, nil
}

func (g *GistSQL) UpdateName(id string) (*Gist, error) {
	row, err := storage.Database.Query("UPDATE gists SET name = $1 WHERE gist_id = $2 AND owner = $3 RETURNING gist_id, name, content, owner", g.Name.String, id, g.OwnerID.String)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't update name")
	}

	var gist Gist

	row.Next()

	err = row.Scan(&gist.ID, &gist.Name, &gist.Content, &gist.OwnerID)

	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't scan gist")
	}

	return &gist, nil
}

func (g *GistSQL) UpdateContent(id string) (*Gist, error) {
	row, err := storage.Database.Query("UPDATE gists SET content = $1 WHERE gist_id = $2 AND owner = $3 RETURNING gist_id, name, content, owner", g.Content.String, id, g.OwnerID.String)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't update content")
	}

	var gist Gist

	row.Next()

	err = row.Scan(&gist.ID, &gist.Name, &gist.Content, &gist.OwnerID)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't scan gist")
	}

	return &gist, nil
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
	row, err := storage.Database.Query("SELECT gist_id, name, content, owner, org_id FROM gists WHERE gist_id = $1 AND owner = $2", id, g.OwnerID.String)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't find gist")
	}
	row.Next()
	var gist Gist
	err = row.Scan(&gist.ID, &gist.Name, &gist.Content, &gist.OwnerID, &gist.OrgID)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't find gist")
	}
	return &gist, nil
}

func (g *GistSQL) FindAll() ([]Gist, error) {
	rows, err := storage.Database.Query("SELECT gist_id, name, content, owner, org_id FROM gists WHERE owner = $1", g.OwnerID.String)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't find gists")
	}
	var gists []Gist
	for rows.Next() {
		var gist GistSQL
		err = rows.Scan(&gist.ID, &gist.Name, &gist.Content, &gist.OwnerID, &gist.OrgID)
		if err != nil {
			log.Error(err)
			return nil, errors.New("couldn't find gists")
		}
		gists = append(gists, Gist{
			ID:      strconv.Itoa(int(gist.ID.Int32)),
			Name:    gist.Name.String,
			Content: gist.Content.String,
			OwnerID: gist.OwnerID.String,
			OrgID: func() *string {
				if gist.OrgID.Valid {
					orgID := strconv.Itoa(int(gist.OrgID.Int32))
					return &orgID
				}
				return nil
			}(),
		})
	}
	return gists, nil
}

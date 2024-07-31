package gists

import (
	"database/sql"
	"errors"

	"github.com/gistapp/api/storage"
	"github.com/gofiber/fiber/v2/log"
)

type GistSQL struct {
	ID      sql.NullInt32
	Name    sql.NullString
	Content sql.NullString
}

type GistModel interface {
	Save() error
}

func (g *GistSQL) Save() error {
	_, err := storage.Database.Exec("INSERT INTO gists(name, content) VALUES ($1, $2)", g.Name.String, g.Content.String)

	if err != nil {
		log.Error(err)
		return errors.New("couldn't create gist")
	}
	return nil
}

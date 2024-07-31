package gists

import (
	"database/sql"
	"errors"
)

type GistServiceImpl struct{}

func (g *GistServiceImpl) Save(name string, content string) error {
	m := GistSQL{
		ID: sql.NullInt32{
			Valid: false,
			Int32: 0,
		},
		Name: sql.NullString{
			String: name,
			Valid:  true,
		},
		Content: sql.NullString{
			String: content,
			Valid:  true,
		},
	}

	err := m.Save()
	if err != nil {
		return errors.New("couldn't insert into database gists")
	}
	return nil
}

var GistService GistServiceImpl = GistServiceImpl{}

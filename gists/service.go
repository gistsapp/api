package gists

import (
	"database/sql"
	"errors"
)

type GistServiceImpl struct{}

func (g *GistServiceImpl) Save(name string, content string) (*Gist, error) {
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

	gist, err := m.Save()
	if err != nil {
		return nil, errors.New("couldn't insert into database gists")
	}
	return gist, nil
}

func (g *GistServiceImpl) UpdateName(id string, name string) error {
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
	}
	err := m.UpdateName(id)
	if err != nil {
		return errors.New("couldn't update name in database gists")
	}
	return nil
}

func (g *GistServiceImpl) FindAll() ([]Gist, error) {
	m := GistSQL{}
	gists, err := m.FindAll()
	if err != nil {
		return nil, errors.New("couldn't get gists")
	}
	return gists, nil
}

var GistService GistServiceImpl = GistServiceImpl{}

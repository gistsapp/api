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

func (g *GistServiceImpl) UpdateContent(id string, content string) error {
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
	}
	err := m.UpdateContent(id)
	if err != nil {
		return errors.New("couldn't update content in database gists")
	}
	return nil
}

func (g *GistServiceImpl) Delete(id string) error {
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
	}
	err := m.Delete(id)
	if err != nil {
		return errors.New("couldn't delete from database gists")
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

func (g *GistServiceImpl) FindByID(id string) (*Gist, error) {
	m := GistSQL{}
	gist, err := m.FindByID(id)
	if err != nil {
		return nil, errors.New("couldn't get gist")
	}
	return gist, nil
}

var GistService GistServiceImpl = GistServiceImpl{}

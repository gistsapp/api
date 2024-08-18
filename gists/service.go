package gists

import (
	"database/sql"
	"errors"
	"strconv"
)

type GistServiceImpl struct{}

func (g *GistServiceImpl) Save(name string, content string, owner_id string, org_id string) (*Gist, error) {
	var m GistSQL

	if org_id != "" {
		org_id_int, err := strconv.Atoi(org_id)

		if err != nil {
			return nil, errors.New("org_id must be an integer")
		}
		m = GistSQL{
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
			OwnerID: sql.NullString{
				String: owner_id,
				Valid:  true,
			},
			OrgID: sql.NullInt32{
				Int32: int32(org_id_int),
				Valid: true,
			},
		}
	} else {
		m = GistSQL{
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
			OwnerID: sql.NullString{
				String: owner_id,
				Valid:  true,
			},
		}
	}

	gist, err := m.Save()
	if err != nil {
		return nil, errors.New("couldn't insert into database gists")
	}
	return gist, nil
}

func (g *GistServiceImpl) UpdateName(id string, name string, owner_id string) error {

	f := gistExists(id, owner_id)

	if f != nil {
		return f
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
	err := m.UpdateName(id)
	if err != nil {
		return errors.New("couldn't update name in database gists")
	}
	return nil
}

func (g *GistServiceImpl) UpdateContent(id string, content string, owner_id string) error {
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
			String: content,
			Valid:  true,
		},

		OwnerID: sql.NullString{
			String: owner_id,
			Valid:  true,
		},
	}

	err = m.UpdateContent(id)
	if err != nil {
		return errors.New("couldn't update content in database gists")
	}
	return nil
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

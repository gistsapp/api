package gists

import (
	"github.com/gistapp/api/storage"
	"github.com/gistsapp/pogo/pogo"
	"github.com/gofiber/fiber/v2/log"
)

type GistRights struct {
	UserID string `json:"user_id" pogo:"user_id"`
	GistID string `json:"gist_id" pogo:"gist_id"`
	Right  string `json:"right" pogo:"rights"`
}

type GistRightsModel interface {
	Save() (*GistRights, error)
	Delete() error
	Update() (*GistRights, error)
	GetByGistID() ([]GistRights, error)
	GetByUserID() ([]GistRights, error)
	GetByGistIDAndUserID() (*GistRights, error)
}

func (gr *GistRights) Save() (*GistRights, error) {
	db := storage.PogoDatabase
	gist := make([]GistRights, 0)
	err := pogo.SuperQuery(db, "INSERT INTO user_gists(user_id, gist_id, rights) VALUES ($1, $2, $3) RETURNING :fields", &gist, gr.UserID, gr.GistID, gr.Right)

	return &gist[0], err
}

func (gr *GistRights) Delete() error {
	db := storage.PogoDatabase
	_, err := db.Exec("DELETE FROM user_gists WHERE user_id = $1 AND gist_id = $2", gr.UserID, gr.GistID)
	return err
}

func (gr *GistRights) Update() (*GistRights, error) {
	db := storage.PogoDatabase
	gist := make([]GistRights, 0)
	err := pogo.SuperQuery(db, "UPDATE user_gists SET rights = $1 WHERE user_id = $2 AND gist_id = $3 RETURNING :fields", &gist, gr.Right, gr.UserID, gr.GistID)
	return &gist[0], err
}

func (gr *GistRights) GetByGistID() ([]GistRights, error) {
	db := storage.PogoDatabase
	gists := make([]GistRights, 0)
	err := pogo.SuperQuery(db, "SELECT :fields FROM user_gists WHERE gist_id = $1", &gists, gr.GistID)
	return gists, err
}

func (gr *GistRights) GetByUserID() ([]GistRights, error) {
	db := storage.PogoDatabase
	gists := make([]GistRights, 0)
	err := pogo.SuperQuery(db, "SELECT :fields FROM user_gists WHERE user_id = $1", &gists, gr.UserID)
	return gists, err
}

func (gr *GistRights) GetByGistIDAndUserID() (*GistRights, error) {
	log.Info(gr.UserID)
	log.Info(gr.GistID)
	db := storage.PogoDatabase
	gists := make([]GistRights, 0)
	err := pogo.SuperQuery(db, "SELECT :fields FROM user_gists WHERE user_id = $1 AND gist_id = $2", &gists, gr.UserID, gr.GistID)
	log.Info(gists)
	if len(gists) <= 0 {
		return nil, nil
	}
	return &gists[0], err
}

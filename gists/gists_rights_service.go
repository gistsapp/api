package gists

import "github.com/gofiber/fiber/v2/log"

type RightGistService interface {
	HasRight(gistID string, userID string, action GistAction) (bool, error)
	AddRight(gistID string, userID string, action GistAction) error
	CanEdit(gistID string, userID string) (bool, error)
	CanRead(gistID string, userID string) (bool, error)
}

type GistGuard struct{}

type GistAction string

const (
	Read  GistAction = "read"
	Write GistAction = "write"
)

func NewRightsGistService() *GistGuard {
	return &GistGuard{}
}

func (r *GistGuard) HasRight(gistID string, userID string, action GistAction) (bool, error) {
	gist_sql := GistSQL{}
	gist, err := gist_sql.FindByID(gistID)

	if err != nil {
		return false, err
	}

	if action == Read {
		if gist.Visibility == string(Public) {
			return true, nil
		}
	}

	gr := GistRights{
		UserID: userID,
		GistID: gistID,
		Right:  string(action),
	}
	gists_rights, err := gr.GetByGistIDAndUserID()
	log.Info(gists_rights)
	if err != nil {
		return false, nil
	}
	if gists_rights == nil {
		return false, nil
	}

	if action == Read && gists_rights.Right == string(Read) {
		return true, nil
	}

	if action == Write && (gists_rights.Right == string(Write) || gists_rights.Right == string(Read)) {
		return true, nil
	}
	return false, nil
}

func (r *GistGuard) CanEdit(gistID string, userID string) (bool, error) {
	return r.HasRight(gistID, userID, Write)
}

func (r *GistGuard) CanRead(gistID string, userID string) (bool, error) {
	return r.HasRight(gistID, userID, Read)
}

func (r *GistGuard) AddRight(gistID string, userID string, action GistAction) error {
	gr := GistRights{
		UserID: userID,
		GistID: gistID,
		Right:  string(action),
	}
	_, err := gr.Save()
	return err
}

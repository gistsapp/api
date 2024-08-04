package auth

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/gistapp/api/user"
	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
)

type AuthServiceImpl struct{}

func (a *AuthServiceImpl) Authenticate(c *fiber.Ctx) error {
	if user, err := goth_fiber.CompleteUserAuth(c); err == nil {
		log.Info(user)
		return nil
	} else {
		return goth_fiber.BeginAuthHandler(c)
	}
}

func (a *AuthServiceImpl) LocalAuth(email string) error {
	return errors.New("not implemented")
}

func (a *AuthServiceImpl) Callback(c *fiber.Ctx) error {
	auth_user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		log.Error(err)
		return ErrCantCompleteAuth
	}

	user_md, _, err := a.GetUser(auth_user)

	if err == nil {
		token, err := utils.CreateToken(user_md.Email, user_md.ID)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"token": token,
		})
	}

	user_md, err = a.Register(auth_user)

	if err != nil {
		return err
	}

	jwt, err := utils.CreateToken(user_md.Email, user_md.ID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"token": jwt,
	})
}

func (a *AuthServiceImpl) GetUser(auth_user goth.User) (*user.User, *AuthIdentity, error) {
	auth_and_user, err := new(AuthIdentitySQL).GetWithUser(auth_user.UserID)
	if err != nil {
		return nil, nil, err
	}

	return &auth_and_user.User, &auth_and_user.AuthIdentity, nil
}

func (a *AuthServiceImpl) Register(auth_user goth.User) (*user.User, error) {
	data, err := json.Marshal(auth_user)
	if err != nil {
		return nil, errors.New("couldn't marshal user")
	}

	user_model := user.UserSQL{
		ID:      sql.NullString{String: auth_user.UserID, Valid: true},
		Email:   sql.NullString{String: auth_user.Email, Valid: true},
		Name:    sql.NullString{String: auth_user.Name, Valid: true},
		Picture: sql.NullString{String: auth_user.AvatarURL, Valid: true},
	}

	user_data, err := user_model.Save()

	if err != nil {
		return nil, err
	}

	auth_identity_model := AuthIdentitySQL{
		Data:       sql.NullString{String: string(data), Valid: true},
		Type:       sql.NullString{String: auth_user.Provider, Valid: true},
		OwnerID:    sql.NullString{String: user_data.ID, Valid: true},
		ProviderID: sql.NullString{String: auth_user.UserID, Valid: true},
	}

	auth_identity, err := auth_identity_model.Save()
	log.Info(auth_identity)
	return user_data, err
}

func (a *AuthServiceImpl) RegisterProviders() {
	goth.UseProviders(
		google.New(utils.Get("GOOGLE_KEY"), utils.Get("GOOGLE_SECRET"), utils.Get("PUBLIC_URL")+"/auth/callback/google"),
		github.New(utils.Get("GITHUB_KEY"), utils.Get("GITHUB_SECRET"), utils.Get("PUBLIC_URL")+"/auth/callback/github"),
	)
}

var AuthService AuthServiceImpl = AuthServiceImpl{}
var ErrCantCompleteAuth = errors.New("can't complete auth")

package auth

import (
	"errors"

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

func (a *AuthServiceImpl) Callback(c *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		log.Error(err)
		return ErrCantCompleteAuth
	}
	log.Info(user)
	return errors.New("not implemented")
}

func (a *AuthServiceImpl) RegisterProviders() {
	goth.UseProviders(
		google.New(utils.Get("GOOGLE_KEY"), utils.Get("GOOGLE_SECRET"), utils.Get("PUBLIC_URL")+"/auth/callback/google"),
		github.New(utils.Get("GITHUB_KEY"), utils.Get("GITHUB_SECRET"), utils.Get("PUBLIC_URL")+"/auth/callback/github"),
	)
}

var AuthService AuthServiceImpl = AuthServiceImpl{}
var ErrCantCompleteAuth = errors.New("can't complete auth")

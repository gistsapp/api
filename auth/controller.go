package auth

import (
	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthControllerImpl struct{}

type AuthLocalValidator struct {
	Email string `json:"email"`
}

type AuthLocalVerificationValidator struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

func (a *AuthControllerImpl) Callback() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := AuthService.Callback(c)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		token_cookie := new(fiber.Cookie)
		token_cookie.Name = "gists.access_token"
		token_cookie.HTTPOnly = true
		token_cookie.Value = token
		c.Cookie(token_cookie)
		return c.Redirect(utils.Get("FRONTEND_URL"))
	}
}

func (a *AuthControllerImpl) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return AuthService.Authenticate(c)
	}
}

func (a *AuthControllerImpl) LocalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		e := new(AuthLocalValidator)
		if err := c.BodyParser(e); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with field email as text")
		}

		if err := AuthService.LocalAuth(e.Email); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		return c.JSON(fiber.Map{"message": "Check your email for the token"})
	}
}

func (a *AuthControllerImpl) VerifyAuthToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		e := new(AuthLocalVerificationValidator)

		if err := c.BodyParser(e); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields token and email as text")
		}

		token := e.Token
		email := e.Email

		jwt_token, err := AuthService.VerifyLocalAuthToken(token, email)

		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		token_cookie := new(fiber.Cookie)
		token_cookie.Name = "gists.access_token"
		token_cookie.HTTPOnly = true
		token_cookie.Value = jwt_token
		c.Cookie(token_cookie)
		return c.Redirect(utils.Get("FRONTEND_URL"))
	}
}

var AuthController AuthControllerImpl = AuthControllerImpl{}

package user

import (
	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2"
)

type IAuthController interface {
	Callback() fiber.Handler
	Authenticate() fiber.Handler
	LocalAuth() fiber.Handler
	VerifyAuthToken() fiber.Handler
}

type AuthControllerImpl struct {
	AuthService IAuthService
}

type AuthLocalValidator struct {
	Email string `json:"email"`
}

type AuthLocalVerificationValidator struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

func (a *AuthControllerImpl) Callback() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := a.AuthService.Callback(c)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		token_cookie := new(fiber.Cookie)
		token_cookie.Name = "gists.access_token"
		token_cookie.HTTPOnly = true
		if utils.Get("ENV") == "development" {
			token_cookie.Secure = false
		} else {
			token_cookie.Domain = ".gists.app" // hardcoded
			token_cookie.Secure = true
		}
		token_cookie.Value = token
		c.Cookie(token_cookie)

		return c.Redirect(utils.Get("FRONTEND_URL") + "/mygist")
	}
}

func (a *AuthControllerImpl) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return a.AuthService.Authenticate(c)
	}
}

func (a *AuthControllerImpl) LocalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		e := new(AuthLocalValidator)
		if err := c.BodyParser(e); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with field email as text")
		}

		if _, err := a.AuthService.LocalAuth(e.Email); err != nil {
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

		jwt_token, err := a.AuthService.VerifyLocalAuthToken(token, email)

		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		token_cookie := new(fiber.Cookie)
		token_cookie.Name = "gists.access_token"
		token_cookie.HTTPOnly = true
		token_cookie.Value = jwt_token

		if utils.Get("ENV") == "development" {
			token_cookie.Secure = false
		} else {
			token_cookie.Domain = ".gists.app" // hardcoded
			token_cookie.Secure = true
		}
		c.Cookie(token_cookie)
		return c.Status(200).JSON(fiber.Map{"message": "You are now logged in"})
	}
}

var AuthController AuthControllerImpl = AuthControllerImpl{}

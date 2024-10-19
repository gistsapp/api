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
	Renew() fiber.Handler
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
		c.Cookie(utils.Cookie("gists.access_token", token.AccessToken))
		c.Cookie(utils.Cookie("gists.refresh_token", token.RefreshToken))

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

		c.Cookie(utils.Cookie("gists.access_token", jwt_token.AccessToken))   //set access token
		c.Cookie(utils.Cookie("gists.refresh_token", jwt_token.RefreshToken)) //set refresh token

		return c.Status(200).JSON(fiber.Map{"message": "You are now logged in"})
	}
}

func (a *AuthControllerImpl) Renew() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user_id := c.Locals("pub").(string)

		tokens, err := a.AuthService.Renew(user_id)

		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		c.Cookie(utils.Cookie("gists.access_token", tokens.AccessToken))   //set access token
		c.Cookie(utils.Cookie("gists.refresh_token", tokens.RefreshToken)) //set refresh token

		return c.Status(200).JSON(fiber.Map{"message": "Welcome back"})
	}
}

var AuthController AuthControllerImpl = AuthControllerImpl{}

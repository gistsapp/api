package auth

import "github.com/gofiber/fiber/v2"

type AuthControllerImpl struct{}

type AuthLocalValidator struct {
	Email string `json:"email"`
}

func (a *AuthControllerImpl) Callback() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return AuthService.Callback(c)
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
		return AuthService.LocalAuth(e.Email)
	}
}

var AuthController AuthControllerImpl = AuthControllerImpl{}

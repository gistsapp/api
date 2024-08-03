package auth

import "github.com/gofiber/fiber/v2"

type AuthControllerImpl struct{}

func (a *AuthControllerImpl) Callback() fiber.Handler {
	return func(c *fiber.Ctx) error {
		AuthService.Callback(c)
		return c.SendString("Register")
	}
}

func (a *AuthControllerImpl) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return AuthService.Authenticate(c)
	}
}

var AuthController AuthControllerImpl = AuthControllerImpl{}

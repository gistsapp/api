package server

import (
	"github.com/gistapp/api/user"
	"github.com/gofiber/fiber/v2"
)

func AuthorizationCookieMiddleware(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("auth")
	if cookie == "" {
		return ctx.Next()
	}
	ctx.Request().Header.Set("Authorization", "Bearer "+cookie)
	return ctx.Next()
}

func AuthNeededMiddleware(ctx *fiber.Ctx) error {
	if ctx.Get("Authorization") == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	raw_token := string(ctx.Request().Header.Peek("Authorization")[7:])
	claims, err := user.AuthService.IsAuthenticated(raw_token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	ctx.Locals("pub", claims.Pub)
	ctx.Locals("email", claims.Email)
	return ctx.Next()
}

package server

import (
	"github.com/gofiber/fiber/v2"
)

func AuthorizationCookieMiddleware(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("gists.access_token")
	if cookie == "" {
		return ctx.Next()
	}
	ctx.Request().Header.Set("Authorization", "Bearer "+cookie)
	return ctx.Next()
}

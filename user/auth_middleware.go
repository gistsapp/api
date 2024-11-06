package user

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func AuthNeededMiddleware(ctx *fiber.Ctx) error {
	if ctx.Get("Authorization") == "" {
		// return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		// 	"error": "Unauthorized",
		// })
		ctx.Locals("pub", "Guest")
		ctx.Locals("email", "Guest")
		return ctx.Next()
	}
	raw_token := string(ctx.Request().Header.Peek("Authorization")[7:])
	claims, err := AuthService.IsAuthenticated(raw_token)
	if err != nil {
		log.Info(err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	ctx.Locals("pub", claims.Pub)
	ctx.Locals("email", claims.Email)
	return ctx.Next()
}

func RenewTokenMiddleware(ctx *fiber.Ctx) error {
	if ctx.Get("Authorization") == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	raw_token := string(ctx.Request().Header.Peek("Authorization")[7:])
	fmt.Println(raw_token)
	_, err := AuthService.IsAuthenticated(raw_token)
	fmt.Println(err)
	if err == nil {
		log.Error(err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token not expired yet",
		})
	}

	refresh_token := ctx.Cookies("gists.refresh_token", "")

	if refresh_token == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token not found",
		})
	}

	renew_claims, err := AuthService.CanRefresh(refresh_token)

	ctx.Locals("pub", renew_claims.Pub)

	return ctx.Next()
}

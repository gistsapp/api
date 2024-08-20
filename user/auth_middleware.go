package user

import "github.com/gofiber/fiber/v2"

func AuthNeededMiddleware(ctx *fiber.Ctx) error {
	if ctx.Get("Authorization") == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	raw_token := string(ctx.Request().Header.Peek("Authorization")[7:])
	claims, err := AuthService.IsAuthenticated(raw_token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	ctx.Locals("pub", claims.Pub)
	ctx.Locals("email", claims.Email)
	return ctx.Next()
}

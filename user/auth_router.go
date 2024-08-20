package user

import "github.com/gofiber/fiber/v2"

type AuthRouter struct {
	Controller IAuthController
}

func (r *AuthRouter) SubscribeRoutes(app *fiber.Router) {
	(*app).Get("/auth/callback/:provider", r.Controller.Callback())
	(*app).Get("/auth/:provider", r.Controller.Authenticate())
	(*app).Post("/auth/local/begin", r.Controller.LocalAuth())
	(*app).Post("/auth/local/verify", r.Controller.VerifyAuthToken())
}

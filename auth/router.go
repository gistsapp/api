package auth

import "github.com/gofiber/fiber/v2"

type AuthRouter struct {
	Controller AuthControllerImpl
}

func (r *AuthRouter) SubscribeRoutes(app *fiber.Router) {
	(*app).Get("/auth/callback/:provider", r.Controller.Callback())
	(*app).Get("/auth/:provider", r.Controller.Authenticate())
}

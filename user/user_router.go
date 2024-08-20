package user

import (
	"github.com/gofiber/fiber/v2"
)

type UserRouter struct {
	Controller *UserControllerImpl
}

func (r *UserRouter) SubscribeRoutes(app *fiber.Router) {
	user_router := (*app).Group("/user", AuthNeededMiddleware)

	user_router.Get("/me", r.Controller.Get())
}

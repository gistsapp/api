package gists

import "github.com/gofiber/fiber/v2"

type GistRouter struct {
	Controller GistControllerImpl
}

func (r *GistRouter) SubscribeRoutes(app *fiber.Router) {
	(*app).Post("/gists", r.Controller.Save())
	(*app).Patch("/gists/:id/name", r.Controller.UpdateName())
	(*app).Get("/gists", r.Controller.FindAll())
}

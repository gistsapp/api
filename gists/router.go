package gists

import "github.com/gofiber/fiber/v2"

type GistRouter struct {
	Controller GistControllerImpl
}

func (r *GistRouter) SubscribeRoutes(app *fiber.Router) {
	(*app).Post("/gists", r.Controller.Save())
	(*app).Patch("/gists/:id/name", r.Controller.UpdateName())
	(*app).Patch("/gists/:id/content", r.Controller.UpdateContent())
	(*app).Get("/gists", r.Controller.FindAll())
	(*app).Get("/gists/:id", r.Controller.FindByID())
	(*app).Delete("/gists/:id", r.Controller.Delete())
}

package gists

import (
	"github.com/gistapp/api/user"
	"github.com/gofiber/fiber/v2"
)

type GistRouter struct {
	Controller GistControllerImpl
}

func (r *GistRouter) SubscribeRoutes(app *fiber.Router) {
	gists_router := (*app).Group("/gists", user.AuthNeededMiddleware)

	gists_router.Post("/", r.Controller.Save())
	gists_router.Patch("/:id/name", r.Controller.UpdateName())
	gists_router.Patch("/:id/content", r.Controller.UpdateContent())
	gists_router.Patch("/:id/language", r.Controller.UpdateLanguage())
	gists_router.Patch("/:id/description", r.Controller.UpdateDescription())
	gists_router.Get("/", r.Controller.FindAll())
	gists_router.Get("/:id", r.Controller.FindByID())
	gists_router.Delete("/:id", r.Controller.Delete())
}

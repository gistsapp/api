package organizations

import (
	"github.com/gistapp/api/user"
	"github.com/gofiber/fiber/v2"
)

type OrganizationRouter struct {
	Controller OrganizationControllerImpl
}

func (r *OrganizationRouter) SubscribeRoutes(app *fiber.Router) {
	organizations_router := (*app).Group("/orgs", user.AuthNeededMiddleware)

	organizations_router.Post("/", r.Controller.Save())
	organizations_router.Get("/", r.Controller.GetAsMember())
	organizations_router.Get("/:id", r.Controller.GetByID())
}

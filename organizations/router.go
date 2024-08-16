package organizations

import "github.com/gofiber/fiber/v2"

type OrganizationRouter struct {
	Controller OrganizationControllerImpl
}

func (r *OrganizationRouter) SubscribeRoutes(app *fiber.Router) {
	organizations_router := (*app).Group("/orgs")

	organizations_router.Post("/", r.Controller.Save())
}

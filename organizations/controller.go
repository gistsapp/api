package organizations

import "github.com/gofiber/fiber/v2"

type OrganizationControllerImpl struct{}

type OrganizationValidator struct {
	Name string `json:"name"`
}

func (c *OrganizationControllerImpl) Save() fiber.Handler {
	return func(c *fiber.Ctx) error {
		org_payload := new(OrganizationValidator)
		owner_id := c.Locals("pub").(string)

		if err := c.BodyParser(org_payload); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields name as text")
		}
		org, err := OrganizationService.Save(org_payload.Name, owner_id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.Status(201).JSON(org)
	}
}

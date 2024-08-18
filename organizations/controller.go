package organizations

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type OrganizationControllerImpl struct{}

type OrganizationValidator struct {
	Name string `json:"name"`
}

func (c *OrganizationControllerImpl) Save() fiber.Handler {
	return func(c *fiber.Ctx) error {
		org_payload := new(OrganizationValidator)
		owner_id := c.Locals("pub").(string)
		//
		if err := c.BodyParser(org_payload); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields name as text")
		}
		org, err := OrganizationService.Save(org_payload.Name, owner_id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		//
		return c.Status(201).JSON(org)
	}
}

func (c *OrganizationControllerImpl) GetAsMember() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user_id := c.Locals("pub").(string)
		log.Info("user_id: ", user_id)

		organizations, err := OrganizationService.GetAsMember(user_id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(organizations)
	}
}

func (c *OrganizationControllerImpl) GetByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		org_id := c.Params("id")
		user_id := c.Locals("pub").(string)
		org, err := OrganizationService.GetByID(org_id, user_id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(org)
	}
}

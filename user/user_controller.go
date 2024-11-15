package user

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type UserControllerImpl struct{}

func (u *UserControllerImpl) Get() fiber.Handler {
	return func(c *fiber.Ctx) error {
		owner_id := c.Locals("pub").(string)

		user, err := UserService.GetUserByID(owner_id)

		if err != nil {
			if errors.Is(err, ErrUserNotFound) {
				return c.Status(404).SendString(err.Error())
			}
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(user)
	}
}

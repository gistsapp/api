package gists

import "github.com/gofiber/fiber/v2"

type GistControllerImpl struct{}

type GistSaveValidator struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (g *GistControllerImpl) Save() fiber.Handler {
	return func(c *fiber.Ctx) error {
		g := new(GistSaveValidator)

		if err := c.BodyParser(g); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields name and content as text")
		}
		gist, err := GistService.Save(g.Name, g.Content)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.Status(201).JSON(gist)
	}
}

func (g *GistControllerImpl) UpdateName() fiber.Handler {
	return func(c *fiber.Ctx) error {
		g := new(GistSaveValidator)
		if err := c.BodyParser(g); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields name and content as text")
		}
		if err := GistService.UpdateName(c.Params("id"), g.Name); err != nil {
			return c.Status(400).SendString(err.Error()) //could be because gist not found
		}
		return c.Status(200).SendString("Gist updated successfully")
	}
}

func (g *GistControllerImpl) FindAll() fiber.Handler {
	return func(c *fiber.Ctx) error {
		gists, err := GistService.FindAll()
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(gists)
	}
}

var GistController GistControllerImpl = GistControllerImpl{}

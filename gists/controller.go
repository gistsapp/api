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

		if err := GistService.Save(g.Name, g.Content); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.Status(200).SendString("Gist created successfully")
	}
}

var GistController GistControllerImpl = GistControllerImpl{}

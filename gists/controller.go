package gists

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type GistControllerImpl struct{}

type GistSaveValidator struct {
	Name        string `json:"name"`
	Content     string `json:"content"`
	OrgID       string `json:"org_id,omitempty"`
	Language    string `json:"language,omitempty"`
	Description string `json:"description,omitempty"`
}

func (g *GistControllerImpl) Save() fiber.Handler {
	return func(c *fiber.Ctx) error {
		g := new(GistSaveValidator)
		owner_id := c.Locals("pub").(string)

		if err := c.BodyParser(g); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields name and content as text")
		}
		gist, err := GistService.Save(g.Name, g.Content, owner_id, g.OrgID, g.Language, g.Description)
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

		owner_id := c.Locals("pub").(string)
		gist, err := GistService.UpdateName(c.Params("id"), g.Name, owner_id)
		if err != nil {
			if err == ErrGistNotFound {
				return c.Status(404).SendString(err.Error())
			}
			return c.Status(400).SendString(err.Error()) //could be because gist not found
		}
		return c.Status(200).JSON(gist)
	}
}

func (g *GistControllerImpl) FindAll() fiber.Handler {
	return func(c *fiber.Ctx) error {
		owner_id := c.Locals("pub").(string)
		gists, err := GistService.FindAll(owner_id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(gists)
	}
}

func (g *GistControllerImpl) FindByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		owner_id := c.Locals("pub").(string)

		gist, err := GistService.FindByID(c.Params("id"), owner_id)
		if err != nil {
			return c.Status(404).SendString(err.Error())
		}

		raw_params := c.Queries()["raw"]
		log.Info(raw_params)
		if raw_params == "true" {
			return c.Status(200).SendString(gist.Content)
		}
		return c.JSON(gist)
	}
}

func (g *GistControllerImpl) RawFindByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		owner_id := c.Locals("pub").(string)

		gist, err := GistService.FindByID(c.Params("id"), owner_id)
		if err != nil {
			return c.Status(404).SendString(err.Error())
		}

		return c.Status(200).SendString(gist.Content)
	}
}

func (g *GistControllerImpl) UpdateContent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		g := new(GistSaveValidator)
		if err := c.BodyParser(g); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields name and content as text")
		}

		owner_id := c.Locals("pub").(string)
		gist, err := GistService.UpdateContent(c.Params("id"), g.Content, owner_id)
		if err != nil {
			if err == ErrGistNotFound {
				return c.Status(404).SendString(err.Error())
			}
			return c.Status(400).SendString(err.Error()) //could be because gist not found
		}
		return c.Status(200).JSON(gist)
	}
}

func (g *GistControllerImpl) UpdateLanguage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		g := new(GistSaveValidator)
		if err := c.BodyParser(g); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields name and content as text")
		}
		owner_id := c.Locals("pub").(string)
		gist, err := GistService.UpdateLanguage(c.Params("id"), g.Language, owner_id)
		if err != nil {
			if err == ErrGistNotFound {
				return c.Status(404).SendString(err.Error())
			}
			return c.Status(400).SendString(err.Error()) //could be because gist not found
		}
		return c.Status(200).JSON(gist)
	}
}

func (g *GistControllerImpl) UpdateDescription() fiber.Handler {
	return func(c *fiber.Ctx) error {
		g := new(GistSaveValidator)
		if err := c.BodyParser(g); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields name and content as text")
		}
		owner_id := c.Locals("pub").(string)
		gist, err := GistService.UpdateDescription(c.Params("id"), g.Description, owner_id)
		if err != nil {
			if err == ErrGistNotFound {
				return c.Status(404).SendString(err.Error())
			}
			return c.Status(400).SendString(err.Error()) //could be because gist not found
		}
		return c.Status(200).JSON(gist)
	}
}

func (g *GistControllerImpl) Delete() fiber.Handler {
	return func(c *fiber.Ctx) error {
		owner_id := c.Locals("pub").(string)
		if err := GistService.Delete(c.Params("id"), owner_id); err != nil {
			if err == ErrGistNotFound {
				return c.Status(404).SendString(err.Error())
			}
			return c.Status(400).SendString(err.Error()) //could be because gist not found
		}
		return c.Status(200).SendString("Gist deleted successfully")
	}
}

var GistController GistControllerImpl = GistControllerImpl{}

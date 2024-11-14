package gists

import (
	"errors"
	"strconv"

	"github.com/gistapp/api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type GistControllerImpl struct {
	gist_guard GistGuard
}

type GistSaveValidator struct {
	Name        string           `json:"name" validate:"required"`
	Content     string           `json:"content" validate:"required"`
	OrgID       utils.ZeroString `json:"org_id,omitempty"`
	Language    string           `json:"language,omitempty"`
	Description string           `json:"description,omitempty"`
	Visibility  string           `json:"visibility,omitempty"`
}

type GistUpdateValidator struct {
	Name        string           `json:"name" validate:"required"`
	Content     string           `json:"content" validate:"required"`
	OrgID       utils.ZeroString `json:"org_id,omitempty"`
	Language    string           `json:"language,omitempty"`
	Description string           `json:"description,omitempty"`
	Visibility  string           `json:"visibility,omitempty"`
}

func defaultGist() *GistSaveValidator {
	return &GistSaveValidator{
		Language:    "plaintext",
		Visibility:  "public",
		Description: "",
	}

}

func (g *GistControllerImpl) Save() fiber.Handler {
	return func(c *fiber.Ctx) error {
		g := defaultGist()
		owner_id := c.Locals("pub").(string)
		log.Info(owner_id)

		if err := c.BodyParser(g); err != nil {
			log.Info(err)
			return c.Status(400).SendString("Request must be valid JSON with fields name and content as text")
		}
		validate := validator.New(validator.WithRequiredStructEnabled())

		err := validate.Struct(g)

		if err != nil {
			log.Error(err)
			return c.Status(400).SendString("Request must be valid JSON with fields name and content as text")
		}

		if g.Visibility != "public" && g.Visibility != "private" {
			return c.Status(400).SendString("Visibility must be either public or private")
		}
		log.Info(g)

		gist, err := GistService.Save(g.Name, g.Content, owner_id, g.OrgID.SqlString(), g.Language, g.Description, g.Visibility)
		if err != nil {
			log.Error(err)
			return c.Status(500).SendString(err.Error())
		}

		return c.Status(201).JSON(gist.ToJSON())
	}
}

func (g *GistControllerImpl) UpdateName() fiber.Handler {
	return func(c *fiber.Ctx) error {
		g_body := new(GistSaveValidator)
		if err := c.BodyParser(g_body); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields name and content as text")
		}

		owner_id := c.Locals("pub").(string)

		can_edit, err := g.gist_guard.CanEdit(c.Params("id"), owner_id)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		if !can_edit {
			return c.Status(403).SendString("You do not have permission to edit this gist")
		}

		gist, err := GistService.UpdateName(c.Params("id"), g_body.Name, owner_id)
		if err != nil {
			if err == ErrGistNotFound {
				return c.Status(404).SendString(err.Error())
			}
			return c.Status(400).SendString(err.Error()) //could be because gist not found
		}
		return c.Status(200).JSON(gist.ToJSON())
	}
}

func (g *GistControllerImpl) FindAll() fiber.Handler {
	return func(c *fiber.Ctx) error {
		owner_id := c.Locals("pub").(string)
		limit_param := c.Query("limit")
		offset_param := c.Query("offset")
		short_param := c.Query("short")

		if limit_param == "" {
			limit_param = "50"
		}

		if offset_param == "" {
			offset_param = "0"
		}

		short := short_param == "true"

		limit, err := strconv.Atoi(limit_param)
		if err != nil {
			return c.Status(400).SendString("limit must be a number")
		}
		offset, err := strconv.Atoi(offset_param)
		if err != nil {
			return c.Status(400).SendString("offset must be a number")
		}

		gists, err := GistService.FindAll(owner_id, limit, offset, short)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		nb_pages, err := GistService.GetPageCount(owner_id, limit)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		gists_json := make([]map[string]interface{}, 0)

		for _, gist := range gists {
			gists_json = append(gists_json, gist.ToJSON())

		}

		return c.JSON(map[string]interface{}{
			"gists":    gists_json,
			"nb_pages": nb_pages,
		})
	}
}

func (g *GistControllerImpl) FindByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		owner_id := c.Locals("pub").(string)

		can_read, err := g.gist_guard.CanRead(c.Params("id"), owner_id)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		if !can_read {
			return c.Status(403).SendString("You do not have permission to read this gist")
		}

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

		can_read, err := g.gist_guard.CanRead(c.Params("id"), owner_id)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		if !can_read {
			return c.Status(403).SendString("You do not have permission to read this gist")
		}

		gist, err := GistService.FindByID(c.Params("id"), owner_id)
		if err != nil {
			return c.Status(404).SendString(err.Error())
		}

		return c.Status(200).SendString(gist.Content)
	}
}

func (g *GistControllerImpl) UpdateContent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		g_body := new(GistSaveValidator)
		if err := c.BodyParser(g_body); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields name and content as text")
		}

		owner_id := c.Locals("pub").(string)

		can_edit, err := g.gist_guard.CanEdit(c.Params("id"), owner_id)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		if !can_edit {
			return c.Status(403).SendString("You do not have permission to edit this gist")
		}

		gist, err := GistService.UpdateContent(c.Params("id"), g_body.Content, owner_id)
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
		g_body := defaultGist()
		if err := c.BodyParser(g_body); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields name and content as text")
		}
		owner_id := c.Locals("pub").(string)

		can_edit, err := g.gist_guard.CanEdit(c.Params("id"), owner_id)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		if !can_edit {
			return c.Status(403).SendString("You do not have permission to edit this gist")
		}

		gist, err := GistService.UpdateLanguage(c.Params("id"), g_body.Language, owner_id)
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
		g_body := new(GistSaveValidator)
		if err := c.BodyParser(g_body); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields name and content as text")
		}
		owner_id := c.Locals("pub").(string)

		can_edit, err := g.gist_guard.CanEdit(c.Params("id"), owner_id)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		if !can_edit {
			return c.Status(403).SendString("You do not have permission to edit this gist")
		}

		gist, err := GistService.UpdateDescription(c.Params("id"), g_body.Description, owner_id)
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

		can_edit, err := g.gist_guard.CanEdit(c.Params("id"), owner_id)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		if !can_edit {
			return c.Status(403).SendString("You do not have permission to edit this gist")
		}

		if err := GistService.Delete(c.Params("id"), owner_id); err != nil {
			if err == ErrGistNotFound {
				return c.Status(404).SendString(err.Error())
			}
			return c.Status(400).SendString(err.Error()) //could be because gist not found
		}
		return c.Status(200).SendString("Gist deleted successfully")
	}
}

func (g *GistControllerImpl) Update() fiber.Handler {
	return func(c *fiber.Ctx) error {
		gist_validate, err := validateGist(c)
		id := c.Params("id")
		owner_id := c.Locals("pub").(string)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		can_edit, err := g.gist_guard.CanEdit(id, owner_id)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		if !can_edit {
			return c.Status(403).SendString("You do not have permission to edit this gist")
		}

		gist, err := GistService.Update(id, gist_validate.Name, gist_validate.OrgID, gist_validate.Content, gist_validate.Language, gist_validate.Description, gist_validate.Visibility, owner_id)
		return c.Status(200).JSON(gist.ToJSON())
	}
}

func validateGist(c *fiber.Ctx) (*GistSaveValidator, error) {
	g := defaultGist()
	owner_id := c.Locals("pub").(string)
	log.Info(owner_id)

	if err := c.BodyParser(g); err != nil {
		log.Info(err)
		return nil, err
	}
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(g)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if g.Visibility != "public" && g.Visibility != "private" {
		return nil, errors.New("Visibility must be either public or private")
	}

	return g, nil
}

var GistController GistControllerImpl = GistControllerImpl{}

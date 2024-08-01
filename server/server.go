package server

import (
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Server struct {
	listenAddr string
	app        *fiber.App
}

type DomainRouter interface {
	SubscribeRoutes(app *fiber.Router)
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		app:        fiber.New(),
	}
}

func (s *Server) Ignite(routers ...DomainRouter) {

	s.app.Get("/", func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/openapi.json",
			Theme:   scalar.ThemeKepler,
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Gists App API Reference",
			},
			DarkMode: true,
		})

		if err != nil {
			log.Error(err)
			return c.Status(500).SendString("Error generating API reference")
		}
		return c.Format(htmlContent)
	})

	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	s.app.Use(logger.New())

	custom_router := s.app.Group("/")

	for _, router := range routers {
		router.SubscribeRoutes(&custom_router)
	}

	log.Fatal(s.app.Listen(s.listenAddr))
}

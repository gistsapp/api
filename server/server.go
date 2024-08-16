package server

import (
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Server struct {
	listenAddr string
	App        *fiber.App
}

type DomainRouter interface {
	SubscribeRoutes(app *fiber.Router)
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		App:        fiber.New(),
	}
}

func (s *Server) Setup(routers ...DomainRouter) {

	s.App.Get("/", func(c *fiber.Ctx) error {
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

	s.App.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     utils.Get("FRONTEND_URL"),
	}))

	s.App.Use(logger.New())

	custom_router := s.App.Group("/")

	for _, router := range routers {
		router.SubscribeRoutes(&custom_router)
	}

}

func (s *Server) Ignite() {
	log.Fatal(s.App.Listen(s.listenAddr))
}

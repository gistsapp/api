package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	listenAddr string
	app        *fiber.App
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		app:        fiber.New(),
	}
}

func (s *Server) Start() {
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! Don't fool me twice!")
	})

	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	log.Fatal(s.app.Listen(s.listenAddr))
}

package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

//go:embed frontend/dist
var distFS embed.FS

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Zen",
	})

	app.Use(logger.New())

	api := app.Group("/api")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	httpFS := http.FS(distFS)
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       httpFS,
		PathPrefix: "frontend/dist",
		Browse:     false,
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

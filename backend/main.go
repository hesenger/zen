package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

//go:embed frontend/dist
var distFS embed.FS

func main() {
	params, err := loadOrCreateParams()
	if err != nil {
		log.Fatal("Failed to load params:", err)
	}
	jwtSecret = []byte(params.JWTSecret)

	appUpdater := NewAppUpdater("/opt/zen/data/setup.json")
	go appUpdater.Start()

	app := fiber.New(fiber.Config{
		AppName: "Zen",
	})

	app.Use(logger.New())

	api := app.Group("/api")
	api.Get("/health", handleHealth)
	api.Get("/check", handleCheck)
	api.Post("/setup", handleSetup)
	api.Post("/login", handleLogin)
	api.Post("/logout", handleLogout)

	httpFS := http.FS(distFS)
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       httpFS,
		PathPrefix: "frontend/dist",
		Browse:     false,
	}))

	log.Printf("Server starting on port 8888")
	if err := app.Listen(":8888"); err != nil {
		log.Fatal(err)
	}
}

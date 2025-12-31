package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	appUpdater := NewDefaultAppUpdater("/opt/zen/data/setup.json")
	go appUpdater.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down, stopping all managed apps...")
		appUpdater.ProcessManager.StopAll()
		os.Exit(0)
	}()

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

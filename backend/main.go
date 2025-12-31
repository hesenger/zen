package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"golang.org/x/crypto/bcrypt"
)

//go:embed frontend/dist
var distFS embed.FS

type App struct {
	Provider string `json:"provider"`
	Key      string `json:"key"`
	Command  string `json:"command"`
}

type SetupData struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	GithubToken string `json:"githubToken"`
	Apps        []App  `json:"apps"`
}

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

	api.Get("/check", func(c *fiber.Ctx) error {
		setupFilePath := "/opt/zen/data/setup.json"
		if _, err := os.Stat(setupFilePath); err == nil {
			return c.JSON(fiber.Map{
				"status": "ready",
			})
		}
		return c.JSON(fiber.Map{
			"status": "pending-setup",
		})
	})

	api.Post("/setup", func(c *fiber.Ctx) error {
		var setupData SetupData
		if err := c.BodyParser(&setupData); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(setupData.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}
		setupData.Password = string(hashedPassword)

		setupFilePath := "/opt/zen/data/setup.json"
		setupDir := filepath.Dir(setupFilePath)
		if err := os.MkdirAll(setupDir, 0755); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to create data directory",
			})
		}

		jsonData, err := json.MarshalIndent(setupData, "", "  ")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to marshal setup data",
			})
		}

		if err := os.WriteFile(setupFilePath, jsonData, 0600); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to write setup file",
			})
		}

		return c.SendStatus(201)
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

package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

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

func handleCheck(c *fiber.Ctx) error {
	setupFilePath := "/opt/zen/data/setup.json"
	if _, err := os.Stat(setupFilePath); err != nil {
		return c.JSON(fiber.Map{
			"status": "pending-setup",
		})
	}

	tokenString := c.Cookies("auth_token")
	if tokenString != "" {
		if _, err := validateJWT(tokenString); err == nil {
			return c.JSON(fiber.Map{
				"status": "authenticated",
			})
		}
	}

	return c.JSON(fiber.Map{
		"status": "ready",
	})
}

func handleSetup(c *fiber.Ctx) error {
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
}

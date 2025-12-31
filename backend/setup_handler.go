package main

import (
	"github.com/gofiber/fiber/v2"
)

var setupService = NewDefaultSetupService()

func handleCheck(c *fiber.Ctx) error {
	tokenString := c.Cookies("auth_token")
	status, err := setupService.CheckSetupStatus(tokenString)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to check setup status",
		})
	}

	return c.JSON(fiber.Map{
		"status": status,
	})
}

func handleSetup(c *fiber.Ctx) error {
	var setupData SetupData
	if err := c.BodyParser(&setupData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := setupService.PerformSetup(setupData); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to complete setup",
		})
	}

	return c.SendStatus(201)
}

package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(getEnvOrDefault("JWT_SECRET", "your-secret-key-change-in-production"))

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func generateJWT(username string) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func validateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
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

	api.Post("/login", func(c *fiber.Ctx) error {
		var loginReq LoginRequest
		if err := c.BodyParser(&loginReq); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid payload",
			})
		}

		setupFilePath := "/opt/zen/data/setup.json"
		fileData, err := os.ReadFile(setupFilePath)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"message": "Invalid credentials",
			})
		}

		var setupData SetupData
		if err := json.Unmarshal(fileData, &setupData); err != nil {
			return c.Status(401).JSON(fiber.Map{
				"message": "Invalid credentials",
			})
		}

		if loginReq.Username != setupData.Username {
			return c.Status(401).JSON(fiber.Map{
				"message": "Invalid credentials",
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(setupData.Password), []byte(loginReq.Password)); err != nil {
			return c.Status(401).JSON(fiber.Map{
				"message": "Invalid credentials",
			})
		}

		token, err := generateJWT(loginReq.Username)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed to generate token",
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "auth_token",
			Value:    token,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			MaxAge:   86400,
		})

		return c.JSON(fiber.Map{
			"message": "Login successful",
		})
	})

	api.Post("/logout", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     "auth_token",
			Value:    "",
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			MaxAge:   -1,
		})

		return c.JSON(fiber.Map{
			"message": "Logged out",
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

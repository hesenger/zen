package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(getEnvOrDefault("JWT_SECRET", "your-secret-key-change-in-production"))

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

func handleLogin(c *fiber.Ctx) error {
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
}

func handleLogout(c *fiber.Ctx) error {
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
}

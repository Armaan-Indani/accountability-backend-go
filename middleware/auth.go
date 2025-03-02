package middleware

import (
	"app/config"
	"log"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(config.Config("SECRET"))},
		ErrorHandler:   jwtError,
		SuccessHandler: jwtSuccessHandler, // Custom success handler
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	log.Println("Error in auth.go middleware: ", err.Error())
	if err.Error() == "missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

// jwtSuccessHandler extracts claims from the JWT and stores them in fiber.Ctx
func jwtSuccessHandler(c *fiber.Ctx) error {

	// Type assert user to *jwt.Token
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid token structure"})
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid token claims"})
	}

	// Extracting userID from JWT payload
	if userID, exists := claims["user_id"].(float64); exists {
		c.Locals("userID", int(userID)) // Store in Fiber context
	} else {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user ID from token",
			"data":    nil,
		})
	}

	return c.Next() // Proceed to the next handler
}

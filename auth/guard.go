package auth

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWT Middleware
var JwtMiddleware = jwtware.New(jwtware.Config{
	SigningKey: jwtware.SigningKey{Key: []byte("secret")},
})

// RoleGuard
func RoleGuard(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Missing or invalid token")
		}

		token, ok := user.(*jwt.Token)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token claims")
		}

		userRoles, ok := claims["roles"].([]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).SendString("Roles not found in token")
		}

		// Convert user roles to a map for efficient lookup
		userRolesMap := make(map[string]bool)
		for _, role := range userRoles {
			if roleStr, ok := role.(string); ok {
				userRolesMap[roleStr] = true
			}
		}

		// Check if the user has at least one of the required roles
		for _, requiredRole := range requiredRoles {
			if userRolesMap[requiredRole] {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).SendString("Access denied")
	}
}

// Auth Middleware
func Accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

// Restricted
func Restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}

func AdminOnly(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome Admin " + name)
}

func UserOnly(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome User " + name)
}

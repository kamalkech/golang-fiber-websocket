package user

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App) {
	var UserRoutes = app.Group("/users")

	UserRoutes.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("list users")
	})
}

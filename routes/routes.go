package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kriengsak.ko/backend-lab/controllers"
)

// Setup registers all routes
func Setup(app *fiber.App) {
	// root
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "hello world"})
	})

	api := app.Group("/api")
	auth := api.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)

	// serve swagger/openapi json and UI
	app.Static("/swagger", "./docs")
	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/index.html")
	})
}

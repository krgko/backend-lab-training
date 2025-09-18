package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "hello world",
		})
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

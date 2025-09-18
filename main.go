package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/kriengsak.ko/backend-lab/database"
	"github.com/kriengsak.ko/backend-lab/models"
	"github.com/kriengsak.ko/backend-lab/routes"
)

func main() {
	// initialize sqlite database
	database.Init("app.db")

	// auto migrate models
	if err := database.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	app := fiber.New()

	routes.Setup(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

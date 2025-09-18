package routes

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kriengsak.ko/backend-lab/controllers"
	"github.com/kriengsak.ko/backend-lab/database"
	"github.com/kriengsak.ko/backend-lab/models"
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

	// profile routes require auth
	profile := api.Group("/profile")
	profile.Use(AuthRequired)
	profile.Get("/", controllers.GetProfile)
	profile.Put("/", controllers.UpdateProfile)

	// serve swagger/openapi json and UI
	app.Static("/swagger", "./docs")
	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/index.html")
	})
}

// AuthRequired is a simple JWT middleware that sets the user in locals
func AuthRequired(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header"})
	}

	tokenStr := parts[1]
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
	}

	sub := claims["sub"]
	if sub == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
	}

	// GORM expects uint for ID; JWT decoder may decode numbers as float64
	var id uint
	switch v := sub.(type) {
	case float64:
		id = uint(v)
	case int:
		id = uint(v)
	case int64:
		id = uint(v)
	case uint:
		id = v
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid subject claim"})
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
	}

	c.Locals("user", user)
	return c.Next()
}

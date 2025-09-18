package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kriengsak.ko/backend-lab/database"
	"github.com/kriengsak.ko/backend-lab/models"
)

// ProfileResponse represents the profile returned to clients
type ProfileResponse struct {
	ID              uint   `json:"id"`
	Email           string `json:"email"`
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	Phone           string `json:"phone,omitempty"`
	MemberCode      string `json:"member_code,omitempty"`
	MembershipLevel string `json:"membership_level,omitempty"`
	Points          int    `json:"points"`
	JoinedAt        string `json:"joined_at"`
}

// UpdateProfileRequest is the payload to update a user's editable profile fields
type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

// GetProfile returns the authenticated user's profile
func GetProfile(c *fiber.Ctx) error {
	u := c.Locals("user")
	if u == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthenticated"})
	}
	user := u.(models.User)

	res := ProfileResponse{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Phone:           user.Phone,
		MemberCode:      user.MemberCode,
		MembershipLevel: user.MembershipLevel,
		Points:          user.Points,
		JoinedAt:        user.CreatedAt.Format("2/1/2006"),
	}

	return c.JSON(res)
}

// UpdateProfile updates editable fields of the authenticated user
func UpdateProfile(c *fiber.Ctx) error {
	u := c.Locals("user")
	if u == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthenticated"})
	}
	user := u.(models.User)

	var body UpdateProfileRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	// only allow changing these fields from the UI
	user.FirstName = body.FirstName
	user.LastName = body.LastName
	user.Phone = body.Phone

	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update profile"})
	}

	// return updated profile
	res := ProfileResponse{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Phone:           user.Phone,
		MemberCode:      user.MemberCode,
		MembershipLevel: user.MembershipLevel,
		Points:          user.Points,
		JoinedAt:        user.CreatedAt.Format("2/1/2006"),
	}

	return c.JSON(res)
}

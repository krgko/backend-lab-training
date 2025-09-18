package models

import "gorm.io/gorm"

// User represents a registered user
type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`

	// profile fields
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	Phone           string `json:"phone,omitempty"`
	MemberCode      string `json:"member_code,omitempty" gorm:"uniqueIndex"`
	MembershipLevel string `json:"membership_level,omitempty"`
	Points          int    `json:"points"`
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type Role uint8

// User permission level
const (
	RoleAdmin Role = iota
	RoleUser
)

type User struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Name         string      `gorm:"not null"`
	Role         Role        `gorm:"not null"` // Default to RoleUser
	CreatedAt    time.Time   `gorm:"not null"`
	UpdatedAt    time.Time   `gorm:"not null"`
	Email        string      `gorm:"not null;uniqueIndex"`
	Password     string      `gorm:"not null"`
	UserSettings Preferences `gorm:"embedded;embeddedPrefix:preferences_"`
	Enabled      bool        `gorm:"not null"`
	DeletedAt    *time.Time  // soft delete
}

// NewUser Pass validated email and hashed password
func NewUser(id uuid.UUID, name string, email string, password string, role Role, preferences Preferences, enabled bool) *User {
	return &User{
		ID:           id,
		Name:         name,
		Role:         role,
		Email:        email,
		Password:     password,
		UserSettings: preferences,
		Enabled:      enabled,
	}
}

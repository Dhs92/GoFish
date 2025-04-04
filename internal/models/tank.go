package models

import (
	"time"

	"github.com/google/uuid"
)

type Tank struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_user_name,priority:1"` // Reference to Users.UserID
	User      User       `gorm:"constraint:OnDelete:CASCADE"`                             // User who owns the tank
	Name      string     `gorm:"not null;uniqueIndex:idx_user_name,priority:2"`           // Name of the tank, must be unique for each user
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	Size      float64    `gorm:"not null"`
	SizeUnit  string     `gorm:"not null"` // gal, liter, etc // Ensures that the combination of UserID and Name is unique
	DeletedAt *time.Time // soft delete
}

func NewTank(ownerID uuid.UUID, name string, size float64, sizeUnit string) *Tank {
	return &Tank{
		ID:       uuid.New(),
		UserID:   ownerID,
		Name:     name,
		Size:     size,
		SizeUnit: sizeUnit,
	}
}

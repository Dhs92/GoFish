package models

import (
	"time"

	"github.com/google/uuid"
)

type Livestock struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;"`
	TankID    uuid.UUID `gorm:"type:uuid;not null"`          // Reference to Tanks.TankID
	Tank      Tank      `gorm:"constraint:OnDelete:CASCADE"` // Tank where the livestock resides
	Name      string    `gorm:"index;not null"`
	Type      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	Size      *float64
	SizeUnit  *string
	Birthday  *time.Time
	Colors    *string
	DeletedAt *time.Time // soft delete
}

func NewLivestock(tankID uuid.UUID, name string, livestockType string, size *float64, sizeUnit *string, birthday *time.Time, colors *string) *Livestock {
	return &Livestock{
		ID:       uuid.New(),
		TankID:   tankID,
		Name:     name,
		Type:     livestockType,
		Size:     size,
		SizeUnit: sizeUnit,
		Birthday: birthday,
		Colors:   colors,
	}
}

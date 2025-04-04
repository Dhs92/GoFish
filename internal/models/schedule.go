package models

import (
	"time"

	"github.com/google/uuid"
)

type ScheduleItem struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null"`          // Reference to Users.UserID
	User         User       `gorm:"constraint:OnDelete:CASCADE"` // User who owns the schedule item
	Name         string     `gorm:"index;not null"`
	ScheduleType string     `gorm:"not null"` // Water change, param check, feeding, etc
	DateTime     time.Time  `gorm:"not null"`
	CreatedAt    time.Time  `gorm:"not null"`
	UpdatedAt    time.Time  `gorm:"not null"`
	Repeat       bool       `gorm:"not null"`
	Consumable   *uuid.UUID `gorm:"type:uuid"` // Reference to StockItems.ItemID
	DeletedAt    *time.Time // soft delete
}

func NewScheduleItem(userID uuid.UUID, name string, scheduleType string, dateTime time.Time, repeat bool, consumable *uuid.UUID) *ScheduleItem {
	return &ScheduleItem{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         name,
		ScheduleType: scheduleType,
		DateTime:     dateTime,
		Repeat:       repeat,
		Consumable:   consumable,
	}
}

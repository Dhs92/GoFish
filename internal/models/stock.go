package models

import (
	"time"

	"github.com/google/uuid"
)

type StockItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`          // Reference to models.User.UserID
	User      User      `gorm:"constraint:OnDelete:CASCADE"` // User who owns the stock item
	Name      string    `gorm:"index;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	ItemType  string    `gorm:"not null"` // Food, chemical, etc
	Quantity  float64   `gorm:"not null"`
	Brand     *string
	DeletedAt *time.Time // soft delete
}

func NewStockItem(ownerID uuid.UUID, name string, itemType string, quantity float64, brand *string) StockItem {
	return StockItem{
		ID:       uuid.New(),
		UserID:   ownerID,
		Name:     name,
		ItemType: itemType,
		Quantity: quantity,
		Brand:    brand,
	}
}

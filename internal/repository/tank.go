package repository

import (
	"fmt"
	"github.com/Dhs92/GoFish/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TankRepository interface {
	FindByID(uuid uuid.UUID) (*models.Tank, error)
	GetAllByOwner(ownerID uuid.UUID) (*[]models.Tank, error)
	GetAll() (*[]models.Tank, error)
	Update(tank models.Tank) error
	Delete(uuid uuid.UUID) error
}

type GormTankRepository struct {
	db *gorm.DB
}

var _ TankRepository = (*GormTankRepository)(nil)

func NewGormTankRepository(db *gorm.DB) *GormTankRepository {
	return &GormTankRepository{db: db}
}

func (tp *GormTankRepository) FindByID(uuid uuid.UUID) (*models.Tank, error) {
	var tank models.Tank
	err := tp.db.Find(&tank, "id = ?", uuid).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find tank: %w", err)
	}

	return &tank, nil
}

func (tp *GormTankRepository) GetAllByOwner(ownerID uuid.UUID) (*[]models.Tank, error) {
	var tanks []models.Tank
	err := tp.db.Where("owner_id = ?", ownerID).Find(&tanks).Error

	if err != nil {
		return nil, fmt.Errorf("get tanks by owner id: %w", err)
	}

	return &tanks, nil
}

func (tp *GormTankRepository) GetAll() (*[]models.Tank, error) {
	var tanks []models.Tank
	err := tp.db.Find(&tanks).Error

	if err != nil {
		return nil, fmt.Errorf("get tanks by owner id: %w", err)
	}

	return &tanks, nil
}

func (tp *GormTankRepository) Update(tank models.Tank) error {
	err := tp.db.Save(&tank).Error

	if err != nil {
		return fmt.Errorf("update tank: %w", err)
	}

	return nil
}

func (tp *GormTankRepository) Delete(uuid uuid.UUID) error {
	err := tp.db.Delete("tank_id = ?", uuid).Error

	if err != nil {
		return fmt.Errorf("delete tank: %w", err)
	}

	return nil
}

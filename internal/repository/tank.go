package repository

import (
	"github.com/Dhs92/GoFish/internal/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type TankRepository interface {
	FindByID(uuid uuid.UUID) (*models.Tank, error)
	FindByOwner(ownerID uuid.UUID) (*models.Tank, error)
	GetAll() ([]*models.Tank, error)
	Update(tank models.Tank) error
	Delete(uuid uuid.UUID) error
}

type GormTankRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

var _ TankRepository = (*GormTankRepository)(nil)

func NewGormTankRepository(db *gorm.DB, logger zerolog.Logger) *GormTankRepository {
	return &GormTankRepository{
		db:     db,
		logger: logger.With().Str("module", "tank").Logger(),
	}
}

func (g *GormTankRepository) FindByID(uuid uuid.UUID) (*models.Tank, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormTankRepository) FindByOwner(ownerID uuid.UUID) (*models.Tank, error) {
	panic("implement me")
}

func (g *GormTankRepository) GetAll() ([]*models.Tank, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormTankRepository) Update(tank models.Tank) error {
	//TODO implement me
	panic("implement me")
}

func (g *GormTankRepository) Delete(uuid uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

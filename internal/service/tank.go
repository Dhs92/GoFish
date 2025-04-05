package service

import (
	"fmt"
	"github.com/Dhs92/GoFish/internal/models"
	"github.com/Dhs92/GoFish/internal/repository"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// TODO: Document

type TankService struct {
	db     repository.TankRepository
	logger zerolog.Logger
}

func NewTankService(db repository.TankRepository, logger zerolog.Logger) *TankService {
	return &TankService{
		db:     db,
		logger: logger,
	}
}

func (ts *TankService) CreateTank(ownerID uuid.UUID, name string, size float64, unit models.SizeUnit) (*models.Tank, error) {
	tank := models.NewTank(ownerID, name, size, unit)

	// Shouldn't ever hit it, but y'know
	if tank == nil {
		ts.logger.Error().Msg("CreateTank called with nil")
		return nil, fmt.Errorf("error creating new tank")
	}

	err := ts.db.Update(*tank)

	if err != nil {
		ts.logger.Error().Err(err).Str("owner", ownerID.String()).Msg("Error updating tank")
		return nil, fmt.Errorf("error updating tank: %w", err)
	}

	return tank, nil
}

func (ts *TankService) FindTankByID(id uuid.UUID) (*models.Tank, error) {
	tank, err := ts.db.FindByID(id)

	if err != nil {
		ts.logger.Error().Err(err).Str("tank_id", id.String()).Msg("Error finding tank")
		return nil, fmt.Errorf("error finding tank: %w", err)
	} else if tank == nil {
		ts.logger.Error().Str("tank_id", id.String()).Msg("Tank ID returned nil")
		return nil, fmt.Errorf("error finding tank: %w", err)
	}

	return tank, nil
}

func (ts *TankService) GetAllTanks() (*[]models.Tank, error) {
	tanks, err := ts.db.GetAll()

	if err != nil {
		ts.logger.Error().Err(err).Msg("Error getting all tanks")
		return nil, fmt.Errorf("error getting all tanks: %w", err)
	}

	return tanks, nil
}

func (ts *TankService) UpdateTank(tank *models.Tank) error {
	if tank == nil {
		ts.logger.Error().Msg("UpdateTank called with nil")
		return fmt.Errorf("error updating tank")
	}

	err := ts.db.Update(*tank)
	if err != nil {
		ts.logger.Error().Err(err).Msg("Error updating tank")
		return fmt.Errorf("error updating tank: %w", err)
	}

	return nil
}

func (ts *TankService) DeleteTank(id uuid.UUID) error {
	err := ts.db.Delete(id)

	if err != nil {
		ts.logger.Error().Err(err).Msg("Error deleting tank")
		return fmt.Errorf("error deleting tank: %w", err)
	}

	return nil
}

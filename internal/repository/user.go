package repository

import (
	"errors"
	"fmt"
	"github.com/Dhs92/GoFish/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	GetAll() (*[]models.User, error)
	Update(user *models.User) error
	DeleteByID(id uuid.UUID) error
}

type GormUserRepository struct {
	db *gorm.DB
}

var _ UserRepository = (*GormUserRepository)(nil)

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (repo *GormUserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := repo.db.First(&user, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if repo.db.Error != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", repo.db.Error)
	}

	return &user, nil
}

func (repo *GormUserRepository) GetAll() (*[]models.User, error) {
	users := make([]models.User, 0)

	err := repo.db.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all users: %w", err)
	}

	return &users, nil
}

// FindByEmail takes a valid, parsed email
func (repo *GormUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := repo.db.First(&user, "email = ?", email).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if repo.db.Error != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", repo.db.Error)
	}

	return &user, nil
}

func (repo *GormUserRepository) Update(user *models.User) error {
	err := repo.db.Save(user).Error

	if err != nil {
		return fmt.Errorf("error while updating user: %w", err)
	}

	return nil
}

func (repo *GormUserRepository) DeleteByID(id uuid.UUID) error {
	err := repo.db.Delete(&models.User{}, "id = ?", id).Error

	if err != nil {
		return fmt.Errorf("error while deleting user: %w", repo.db.Error)
	}

	return nil
}

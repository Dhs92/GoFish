package service

import (
	"fmt"
	"github.com/Dhs92/GoFish/internal/models"
	"github.com/Dhs92/GoFish/internal/repository"
	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
	"github.com/rs/zerolog"
	"net/mail"
)

//TODO: Document

type UserService struct {
	db     repository.UserRepository
	logger zerolog.Logger
}

func NewUserService(db repository.UserRepository, logger zerolog.Logger) *UserService {
	logger.Info().Msg("Initializing UserService")
	return &UserService{
		db:     db,
		logger: logger.With().Str("module", "user_service").Logger(),
	}
}

func (us *UserService) CreateUser(name string, email string, password string, role models.Role) (*uuid.UUID, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		us.logger.Error().Err(err).Msg("Error hashing password")
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	if !validateEmail(email) {
		us.logger.Error().Str("email", email).Msg("Invalid email address")
		return nil, fmt.Errorf("invalid email")
	}

	id := uuid.New()
	err = us.db.Update(models.NewUser(
		id,
		name,
		email,
		hashedPassword,
		role,
		models.DefaultPreferences(),
		true,
	))
	if err != nil {
		us.logger.Error().Err(err).Msg("Error creating/updating user")
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	us.logger.Info().Str("user_id", id.String()).Msg("User created")

	return &id, nil
}

func (us *UserService) FindUser(id uuid.UUID) (*models.User, error) {
	user, err := us.db.FindByID(id)

	if err != nil {
		us.logger.Error().Err(err).Msg("Error finding user")
		return nil, fmt.Errorf("error finding user: %w", err)
	} else if user == nil {
		us.logger.Error().Str("user_id", id.String()).Msg("Find User by ID returned Nil")
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (us *UserService) UpdateUser(user *models.User) error {
	err := us.db.Update(user)
	if err != nil {
		us.logger.Error().Err(err).Msg("Error updating user")
		return fmt.Errorf("error updating user: %w", err)
	}

	us.logger.Info().Msg("User updated")

	return nil
}

func (us *UserService) DeleteUser(user *models.User) error {
	err := us.db.DeleteByID(user.ID)

	if err != nil {
		us.logger.Error().Err(err).Msg("Error deleting user")
		return fmt.Errorf("error deleting user: %w", err)
	}

	us.logger.Info().Str("user_id", user.ID.String()).Msg("User deleted")

	return nil
}

func (us *UserService) GetAll() (*[]models.User, error) {
	users, err := us.db.GetAll()
	if err != nil {
		us.logger.Error().Err(err).Msg("Error getting all users")
		return nil, fmt.Errorf("error getting all users: %w", err)
	}

	return users, nil
}

func (us *UserService) VerifyPassword(user *models.User, password string) error {
	passed, err := argon2.VerifyEncoded([]byte(password), []byte(user.Password))
	if err != nil {
		us.logger.Error().Err(err).Msg("Error verifying password")
		return fmt.Errorf("error verifying password: %w", err)
	}

	if !passed {
		us.logger.Debug().Msg("Invalid password")
		return fmt.Errorf("invalid password")
	}

	return nil
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Returns empty string on error
func hashPassword(password string) (string, error) {
	argon := argon2.DefaultConfig()

	hash, err := argon.HashEncoded([]byte(password))

	if err != nil {
		return "", err
	}

	hashString := string(hash)

	return hashString, nil
}

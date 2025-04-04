package service

import (
	"errors"
	"github.com/Dhs92/GoFish/internal/models"
	"github.com/Dhs92/GoFish/internal/repository"
	"github.com/Dhs92/GoFish/internal/service"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

type UserRepositoryOptional func(m *MockUserRepository)

type MockUserRepository struct {
	db  map[uuid.UUID]*models.User
	err error
}

var _ repository.UserRepository = (*MockUserRepository)(nil)

func UserRepositoryWithError(err error) UserRepositoryOptional {
	return func(m *MockUserRepository) {
		m.err = err
	}
}

func UserRepositoryWithUser(u *models.User) UserRepositoryOptional {
	return func(m *MockUserRepository) {
		m.db[u.ID] = u
	}
}

func NewMockUserRepository(opts ...UserRepositoryOptional) *MockUserRepository {
	db := &MockUserRepository{db: make(map[uuid.UUID]*models.User)}

	for _, opt := range opts {
		opt(db)
	}

	return db
}

func (m MockUserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	user, ok := m.db[id]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (m MockUserRepository) FindByEmail(email string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockUserRepository) GetAll() (*[]models.User, error) {
	panic("implement me")
}

func (m MockUserRepository) Update(user *models.User) error {
	m.db[user.ID] = user

	return nil
}

func (m MockUserRepository) DeleteByID(id uuid.UUID) error {
	delete(m.db, id)

	return nil
}

func TestFindByID_Success(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	mu := NewMockUserRepository(
		UserRepositoryWithUser(models.NewUser(id1, "Alice", "alice@example.com", "password123", models.RoleAdmin, models.DefaultPreferences(), true)),
		UserRepositoryWithUser(models.NewUser(id2, "Jake", "jake@example.com", "password123", models.RoleUser, models.DefaultPreferences(), true)),
	)

	service.NewUserService(mu, log.Logger.Level(zerolog.Disabled))

	user, err := mu.FindByID(id1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, id1, user.ID)
}

package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Dhs92/GoFish/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	assert.NoError(t, err)

	closeFunc := func() {
		_ = db.Close()
	}

	return gormDB, mock, closeFunc
}

func TestFindByID_Success(t *testing.T) {
	gormDB, mock, closeFunc := setupTestDB(t)
	defer closeFunc()

	repo := repository.NewGormUserRepository(gormDB)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
			AddRow(userID, "Alice", "alice@example.com", "hashedpassword"))

	user, err := repo.FindByID(userID)
	assert.NotNil(t, user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, "alice@example.com", user.Email)
}

func TestFindByID_Error(t *testing.T) {
	gormDB, mock, closeFunc := setupTestDB(t)
	defer closeFunc()

	repo := repository.NewGormUserRepository(gormDB)
	wrongUserID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(wrongUserID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}))

	user, err := repo.FindByID(wrongUserID)
	assert.Nil(t, user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByEmail_Success(t *testing.T) {
	gormDB, mock, closeFunc := setupTestDB(t)
	defer closeFunc()

	repo := repository.NewGormUserRepository(gormDB)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs("alex@example.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
			AddRow(userID, "Alice", "alice@example.com", "hashedpassword"))

	user, err := repo.FindByEmail("alex@example.com")
	assert.NotNil(t, user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, "Alice", user.Name)
}

func TestFindByEmail_Error(t *testing.T) {
	gormDB, mock, closeFunc := setupTestDB(t)
	defer closeFunc()

	repo := repository.NewGormUserRepository(gormDB)
	wrongUserEmail := "jake@example.com"

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(wrongUserEmail, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}))

	user, err := repo.FindByEmail(wrongUserEmail)
	assert.Nil(t, user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

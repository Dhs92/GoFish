package db_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Dhs92/GoFish/db"
	"github.com/matthewhartstonge/argon2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func setupMongoContainer(t *testing.T) (testcontainers.Container, string) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := mongoC.Host(ctx)
	require.NoError(t, err)

	port, err := mongoC.MappedPort(ctx, "27017")
	require.NoError(t, err)

	uri := "mongodb://" + host + ":" + port.Port()
	return mongoC, uri
}

func TestConnect(t *testing.T) {
	mongoC, uri := setupMongoContainer(t)
	defer mongoC.Terminate(context.Background())

	ctx := context.Background()
	database, err := db.Connect(ctx, uri, "testdb")
	require.NoError(t, err)
	assert.NotNil(t, database)
}

func TestCreateCollections(t *testing.T) {
	mongoC, uri := setupMongoContainer(t)
	defer mongoC.Terminate(context.Background())

	ctx := context.Background()
	mockDatabase, err := db.Connect(ctx, uri, "testdb")
	require.NoError(t, err)

	err = mockDatabase.CreateCollections(ctx)
	require.NoError(t, err)

	collectionNames, err := mockDatabase.Database.ListCollectionNames(ctx, bson.D{})
	require.NoError(t, err)

	assert.ElementsMatch(t, collectionNames, []string{"users", "scheduleItems", "stockItems", "tanks"})
}

func TestCreateIndexes(t *testing.T) {
	mongoC, uri := setupMongoContainer(t)
	defer mongoC.Terminate(context.Background())

	ctx := context.Background()
	database, err := db.Connect(ctx, uri, "testdb")
	require.NoError(t, err)

	// Assuming CreateUserIndexes, CreateScheduleItemIndexes, CreateStockItemIndexes, CreateTankIndexes are implemented
	err = database.CreateIndexes(ctx)
	require.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	mongoC, uri := setupMongoContainer(t)
	defer mongoC.Terminate(context.Background())

	ctx := context.Background()
	mockDatabase, err := db.Connect(ctx, uri, "testdb")
	require.NoError(t, err)

	argon := argon2.DefaultConfig()

	passwordHash, err := argon.HashEncoded([]byte("password"))

	require.NoError(t, err)

	user, _ := db.NewUser("John Doe", "test@test.com", string(passwordHash))

	result, err := mockDatabase.Create(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, result)

	newUser := mockDatabase.FindDocument(ctx, bson.M{"name": "John Doe"}, "users")

	require.NoError(t, err)
	var userResult db.User
	err = newUser.Decode(&userResult)

	if err != nil {
		fmt.Println("No user found")
		fmt.Printf("%v\n", userResult)
		require.NoError(t, err)
	} else {
		fmt.Println("User found")
		fmt.Printf("%v\r\n", userResult)
	}
}

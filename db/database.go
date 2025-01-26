package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	client   *mongo.Client
	database *mongo.Database
}

func Connect(ctx context.Context, uri, dbName string) (*Database, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return &Database{
		client:   client,
		database: db,
	}, nil
}

func (db *Database) CreateCollections(ctx context.Context) error {
	collections := []string{"users", "scheduleItems", "stockItems", "tanks"}

	for _, collectionName := range collections {
		// Ensure collection is created
		if err := db.database.CreateCollection(ctx, collectionName); err != nil && !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}

	return nil
}

func (db *Database) CreateIndexes(ctx context.Context) error {
	if err := db.CreateUserIndexes(ctx); err != nil {
		return err
	}
	if err := db.CreateScheduleItemIndexes(ctx); err != nil {
		return err
	}
	if err := db.CreateStockItemIndexes(ctx); err != nil {
		return err
	}
	if err := db.CreateTankIndexes(ctx); err != nil {
		return err
	}
	return nil
}

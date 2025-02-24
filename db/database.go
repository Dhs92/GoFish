package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	Database *mongo.Database
}

type DocumentInterface interface {
	CollectionName() string
	ObjectID() bson.ObjectID
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
		Database: db,
	}, nil
}

func (db *Database) Create(ctx context.Context, doc DocumentInterface) (*mongo.InsertOneResult, error) {
	collection := db.Database.Collection(doc.CollectionName())
	return collection.InsertOne(ctx, doc)
}

// Does not accept a document interface because it needs to be able to accept any filter
func (db *Database) FindDocument(ctx context.Context, filter interface{}, collectionName string) *mongo.SingleResult {
	collection := db.Database.Collection(collectionName)
	return collection.FindOne(ctx, filter)
}

func (db *Database) Update(ctx context.Context, filter, update DocumentInterface) (*mongo.UpdateResult, error) {
	collection := db.Database.Collection(filter.CollectionName())
	return collection.UpdateOne(ctx, filter, update)
}

func (db *Database) Delete(ctx context.Context, filter DocumentInterface) (*mongo.DeleteResult, error) {
	collection := db.Database.Collection(filter.CollectionName())
	return collection.DeleteOne(ctx, filter)
}

func (db *Database) CreateCollections(ctx context.Context) error {
	collections := []string{"users", "scheduleItems", "stockItems", "tanks"}

	for _, collectionName := range collections {
		// Ensure collection is created
		if err := db.Database.CreateCollection(ctx, collectionName); err != nil && !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}

	return nil
}

// CreateIndexes creates indexes for all collections in the database.
// It returns an error if any index creation fails.
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

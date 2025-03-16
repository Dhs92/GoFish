package db

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Collection names
const (
	UsersCollection         = "users"
	ScheduleItemsCollection = "scheduleItems"
	StockItemsCollection    = "stockItems"
	TanksCollection         = "tanks"
)

type Database struct {
	Database *mongo.Database
}

type DocumentInterface interface {
	CollectionName() string
	ObjectID() bson.ObjectID
}

func Connect(ctx context.Context, uri, dbName string) (*Database, error) {
	log.Info().Msg("Connecting to database")
	log.Debug().Str("uri", uri).Msg("Connecting to database")
	log.Debug().Str("dbName", dbName).Msg("Connecting to database")
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to database")
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error pinging database")
		return nil, err
	}

	log.Info().Msg("Connected to database")
	db := client.Database(dbName)
	return &Database{
		Database: db,
	}, nil
}

func (db *Database) Create(ctx context.Context, doc any /*TODO: Verify it can be bson serialized*/, collectionName string) (*mongo.InsertOneResult, error) {
	log.Info().Msg("Creating document")
	log.Debug().Interface("doc", doc).Msg("Creating document")
	collection := db.Database.Collection(collectionName)
	return collection.InsertOne(ctx, doc)
}

// Does not accept a document interface because it needs to be able to accept any filter
func (db *Database) FindDocument(ctx context.Context, filter any, collectionName string) *mongo.SingleResult {
	log.Info().Msg("Finding document")
	log.Debug().Interface("filter", filter).Msg("Finding document")
	log.Debug().Str("collectionName", collectionName).Msg("Finding document")
	collection := db.Database.Collection(collectionName)
	return collection.FindOne(ctx, filter)
}

func (db *Database) Update(ctx context.Context, filter any, update any, collectionName string) (*mongo.UpdateResult, error) {
	log.Info().Msg("Updating document")
	log.Debug().Interface("filter", filter).Msg("Updating document")
	log.Debug().Interface("update", update).Msg("Updating document")
	collection := db.Database.Collection(collectionName)
	return collection.UpdateOne(ctx, filter, update)
}

func (db *Database) Delete(ctx context.Context, filter any, collectionName string) (*mongo.DeleteResult, error) {
	log.Info().Msg("Deleting document")
	log.Debug().Interface("filter", filter).Msg("Deleting document")
	collection := db.Database.Collection(collectionName)
	return collection.DeleteOne(ctx, filter)
}

func (db *Database) CreateCollections(ctx context.Context) error {
	// List of collections to create
	collections := []string{"users", "scheduleItems", "stockItems", "tanks"}

	// Create collections
	log.Info().Msg("Creating collections")
	log.Debug().Strs("collections", collections).Msg("Creating collections")
	for _, collectionName := range collections {
		// Ensure collection is created
		if err := db.Database.CreateCollection(ctx, collectionName); err != nil && !mongo.IsDuplicateKeyError(err) {
			log.Error().Err(err).Msg("Error creating collection")
			return err
		}
	}

	return nil
}

// CreateIndexes creates indexes for all collections in the database.
// It returns an error if any index creation fails.
func (db *Database) CreateIndexes(ctx context.Context) error {
	// Create indexes for each collection
	log.Info().Msg("Creating indexes")
	if err := db.CreateUserIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Error creating user indexes")
		return err
	}
	if err := db.CreateScheduleItemIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Error creating schedule item indexes")
		return err
	}
	if err := db.CreateStockItemIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Error creating stock item indexes")
		return err
	}
	if err := db.CreateTankIndexes(ctx); err != nil {
		log.Error().Err(err).Msg("Error creating tank indexes")
		return err
	}
	return nil
}

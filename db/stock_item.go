package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type StockItem struct {
	ID       bson.ObjectID `bson:"_id,omitempty"` // ItemID
	Owner    bson.ObjectID `bson:"owner"`         // Reference to Users._id
	Name     string        `bson:"name"`
	Type     string        `bson:"type"` // ItemType
	Quantity float64       `bson:"quantity"`
	Brand    *string       `bson:"brand,omitempty"`
}

func (db *Database) CreateStockItemIndexes(ctx context.Context) error {
	collection := db.database.Collection("stockItems")
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"owner": 1, // Create an ascending index on the Owner field
		},
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

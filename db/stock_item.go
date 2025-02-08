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

func (s *StockItem) CollectionName() string {
	return "stockItems"
}

func (db *Database) CreateStockItemIndexes(ctx context.Context) error {
	collection := db.Database.Collection("stockItems")
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"owner": 1, // Create an ascending index on the Owner field
		},
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func NewStockItem(owner bson.ObjectID, name string, itemType string, quantity float64, brand *string) *StockItem {
	return &StockItem{
		Owner:    owner,
		Name:     name,
		Type:     itemType,
		Quantity: quantity,
		Brand:    brand,
	}
}

package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ScheduleItem struct {
	ID           bson.ObjectID  `bson:"_id,omitempty"` // ScheduleID
	UserID       bson.ObjectID  `bson:"userId"`        // Reference to Users._id
	Name         string         `bson:"name"`
	ScheduleType string         `bson:"scheduleType"`
	DateTime     time.Time      `bson:"dateTime"`
	Repeat       bool           `bson:"repeat"`
	Consumable   *ConsumableRef `bson:"consumable,omitempty"`
}

type ConsumableRef struct {
	ItemID   bson.ObjectID `bson:"itemId"` // Reference to StockItems._id
	Quantity float64       `bson:"quantity"`
}

func (db *Database) CreateScheduleItemIndexes(ctx context.Context) error {
	collection := db.database.Collection("scheduleItems")
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"userId": 1, // Create an ascending index on the UserID field
		},
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

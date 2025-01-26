package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Tank struct {
	ID        bson.ObjectID `bson:"_id,omitempty"` // TankID
	Owner     bson.ObjectID `bson:"owner"`         // Reference to Users._id
	Name      string        `bson:"name"`
	Size      float64       `bson:"size"`
	SizeUnit  string        `bson:"sizeUnit"`
	Livestock []Livestock   `bson:"livestock"`
}

type Livestock struct {
	ID       bson.ObjectID `bson:"_id,omitempty"` // LivestockID
	Name     string        `bson:"name"`
	Species  string        `bson:"species"`
	Size     *float64      `bson:"size,omitempty"`
	SizeUnit *string       `bson:"sizeUnit,omitempty"`
	Birthday *time.Time    `bson:"birthday,omitempty"`
	Colors   *string       `bson:"colors,omitempty"`
}

func (db *Database) CreateTankIndexes(ctx context.Context) error {
	collection := db.database.Collection("tanks")
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"owner": 1, // Create an ascending index on the Owner field
		},
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

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

func NewTank(owner bson.ObjectID, name string, size float64, sizeUnit string) *Tank {
	return &Tank{
		Owner:    owner,
		Name:     name,
		Size:     size,
		SizeUnit: sizeUnit,
	}
}

func NewLivestock(name, species string, size *float64, sizeUnit *string, birthday *time.Time, colors *string) *Livestock {
	return &Livestock{
		Name:     name,
		Species:  species,
		Size:     size,
		SizeUnit: sizeUnit,
		Birthday: birthday,
		Colors:   colors,
	}
}

func (db *Database) GetTank(ctx context.Context, tankID bson.ObjectID) (*Tank, error) {
	collection := db.database.Collection("tanks")
	var tank Tank
	err := collection.FindOne(ctx, bson.M{"_id": tankID}).Decode(&tank)
	return &tank, err
}

func (db *Database) GetTanks(ctx context.Context, owner bson.ObjectID) ([]Tank, error) {
	collection := db.database.Collection("tanks")
	cursor, err := collection.Find(ctx, bson.M{"owner": owner})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tanks []Tank
	if err = cursor.All(ctx, &tanks); err != nil {
		return nil, err
	}
	return tanks, nil
}

func (db *Database) UpdateTank(ctx context.Context, tankID bson.ObjectID, tank *Tank) error {
	collection := db.database.Collection("tanks")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": tankID}, bson.M{"$set": tank})

	return err
}

func (db *Database) DeleteTank(ctx context.Context, tankID bson.ObjectID) error {
	collection := db.database.Collection("tanks")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": tankID})

	return err
}

func (db *Database) AddLivestock(ctx context.Context, tankID bson.ObjectID, livestock *Livestock) error {
	collection := db.database.Collection("tanks")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": tankID}, bson.M{"$push": bson.M{"livestock": livestock}})

	return err
}

func (db *Database) RemoveLivestock(ctx context.Context, tankID bson.ObjectID, livestockID bson.ObjectID) error {
	collection := db.database.Collection("tanks")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": tankID}, bson.M{"$pull": bson.M{"livestock": bson.M{"_id": livestockID}}})

	return err
}

func (db *Database) UpdateLivestock(ctx context.Context, tankID bson.ObjectID, livestockID bson.ObjectID, livestock *Livestock) error {
	collection := db.database.Collection("tanks")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": tankID, "livestock._id": livestockID}, bson.M{"$set": bson.M{"livestock.$": livestock}})

	return err
}

func (db *Database) GetLivestock(ctx context.Context, tankID bson.ObjectID, livestockID bson.ObjectID) (*Livestock, error) {
	collection := db.database.Collection("tanks")
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"_id": tankID}}},
		bson.D{{Key: "$unwind", Value: "$livestock"}},
		bson.D{{Key: "$match", Value: bson.M{"livestock._id": livestockID}}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var livestock Livestock
	if cursor.Next(ctx) {
		if err := cursor.Decode(&livestock); err != nil {
			return nil, err
		}
		return &livestock, nil
	}
	return nil, mongo.ErrNoDocuments
}

func (db *Database) FindLivestockByName(ctx context.Context, ownerID bson.ObjectID, name string) ([]Tank, error) {
	collection := db.database.Collection("tanks")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "owner", Value: ownerID}}}},
		{{Key: "$unwind", Value: "$livestock"}},
		{{Key: "$match", Value: bson.D{{Key: "livestock.name", Value: name}}}}}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []Tank
	for cursor.Next(ctx) {
		var tank Tank
		if err := cursor.Decode(&tank); err != nil {
			return nil, err
		}
		results = append(results, tank)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}

func (db *Database) CreateTank(ctx context.Context, tank *Tank) error {
	collection := db.database.Collection("tanks")
	_, err := collection.InsertOne(ctx, tank)

	return err
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

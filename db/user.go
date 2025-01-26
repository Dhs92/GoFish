package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID          bson.ObjectID `bson:"_id,omitempty"` // MongoDB ObjectID
	Name        string        `bson:"name"`
	Email       string        `bson:"email"`
	Password    string        `bson:"password"`
	CreatedDate time.Time     `bson:"createdDate"`
	Settings    UserSettings  `bson:"settings"`
}

type UserSettings struct {
	Theme         string `bson:"theme"`
	Notifications bool   `bson:"notifications"`
}

func (db *Database) CreateUser(ctx context.Context, name, email, pwHash string, settings UserSettings) (*User, error) {
	now := time.Now()
	user := &User{
		Name:        name,
		Email:       email,
		Password:    pwHash,
		CreatedDate: now,
		Settings:    settings,
	}

	collection := db.database.Collection("users")
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	id := result.InsertedID.(bson.ObjectID)
	user.ID = id

	return user, nil
}

func (db *Database) CreateUserIndexes(ctx context.Context) error {
	collection := db.database.Collection("users")
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"email": 1, // Create an ascending index on the Email field
		},
		Options: options.Index().SetUnique(true), // Ensure the Email field is unique
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (db *Database) DeleteUser(ctx context.Context, userID bson.ObjectID) error {
	collection := db.database.Collection("users")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": userID})
	return err
}

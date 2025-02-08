package db

import (
	"context"
	"errors"
	"net/mail"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID          bson.ObjectID `bson:"_id,omitempty"` // MongoDB ObjectID
	Name        string        `bson:"name,omitempty"`
	Email       string        `bson:"email,omitempty"`
	Password    string        `bson:"password,omitempty"`
	CreatedDate bson.DateTime `bson:"createdDate,omitempty"`
	Settings    UserSettings  `bson:"settings,omitempty"`
	Enabled     bool          `bson:"enabled,omitempty"`
}

type UserSettings struct {
	Theme         string `bson:"theme"`
	Notifications bool   `bson:"notifications"`
	Timezone      string `bson:"timezone"`
	Locale        string `bson:"locale"`
	StartOfWeek   string `bson:"startOfWeek"`
}

func (u *User) CollectionName() string {
	return "users"
}

func (u *User) ObjectID() bson.ObjectID {
	return u.ID
}

func VerifyEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func InitialUserPreferences() UserSettings {
	return UserSettings{
		Theme:         "light",
		Notifications: true,
		Timezone:      "UTC",
		Locale:        "en-US",
		StartOfWeek:   "Sunday",
	}
}

// NewUser creates a new User object with the provided name, email, and password hash. It returns an error if the email address is invalid.
func NewUser(name string, email string, pwHash string) (*User, error) {
	if VerifyEmail(email) {
		return &User{
			Name:        name,
			Email:       email,
			Password:    pwHash,
			CreatedDate: bson.NewDateTimeFromTime(time.Now()),
			Settings:    InitialUserPreferences(),
			Enabled:     true,
		}, nil
	} else {
		return nil, errors.New("invalid email address")
	}
}

func (db *Database) CreateUserIndexes(ctx context.Context) error {
	collection := db.Database.Collection("users")
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"email": 1, // Create an ascending index on the Email field
		},
		Options: options.Index().SetUnique(true), // Ensure the Email field is unique
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// User ...
type User struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name               string             `json:"name,omitempty" bson:"name,omitempty"`
	DateOfBirth        string             `json:"dateOfBirth,omitempty" bson:"dateOfBirth,omitempty"`
	PhoneNo            string             `json:"phoneNo,omitempty" bson:"phoneNo,omitempty"`
	EmailAddress       string             `json:"emailAddress,omitempty" bson:"emailAddress,omitempty"`
	CreationTimeStamp  time.Time          `json:"creationTimeStamp,omitempty" bson:"creationTimeStamp,omitempty"`
}

var userCollectionName = "users"

func getUsers(db *mongo.Database, start, count int) ([]User, error) {
	col := usersCollection(db)
	ctx := dbContext(30)

	cursor, err := col.Find(ctx, bson.M{}) // find all
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (user *User) getUser(db *mongo.Database) error {
	col := usersCollection(db)
	ctx := dbContext(30)

	filter := bson.M{"_id": user.ID}
	err := col.FindOne(ctx, filter).Decode(&user)
	return err
}

func (user *User) createUser(db *mongo.Database) (*mongo.InsertOneResult, error) {
	col := usersCollection(db)
	ctx := dbContext(30)

	user.CreationTimeStamp = time.Now()
	result, err := col.InsertOne(ctx, user)
	// Convert to map[string]string
	// id := map[string]string{"_id": result.InsertedID.(primitive.ObjectID).Hex()}
	return result, err
}

// helpers
func usersCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection(userCollectionName)
}
package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// Contact ...
type Contact struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ContactIDOne       string             `json:"_idOne,omitempty" bson:"_idOne,omitempty"`
	ContactIDTwo       string             `json:"_idTwo,omitempty" bson:"_idTwo,omitempty"`
	TimeStamp          time.Time          `json:"timeStamp,omitempty" bson:"timeStamp,omitempty"`
}

var contactCollectionName = "contacts"

func getContacts(db *mongo.Database, start, count int) ([]Contact, error) {
	col := contactsCollection(db)
	ctx := dbContext(30)

	cursor, err := col.Find(ctx, bson.M{}) // find all
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var contacts []Contact
	for cursor.Next(ctx) {
		var contact Contact
		cursor.Decode(&contact)
		contacts = append(contacts, contact)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func (contact *Contact) createContact(db *mongo.Database) (*mongo.InsertOneResult, error) {
	col := contactsCollection(db)
	ctx := dbContext(30)
	contact.TimeStamp = time.Now()

	result, err := col.InsertOne(ctx, contact)
	return result, err
}

// helpers
func contactsCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection(contactCollectionName)
}
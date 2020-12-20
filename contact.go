package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"fmt"
)

// Contact ...
type Contact struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ContactIDOne       string             `json:"_idOne,omitempty" bson:"_idOne,omitempty"`
	ContactIDTwo       string             `json:"_idTwo,omitempty" bson:"_idTwo,omitempty"`
	TimeStamp          time.Time          `json:"timeStamp,omitempty" bson:"timeStamp,omitempty"`
}

var contactCollectionName = "contacts"

func getContactTracing(db *mongo.Database, start, count int, userID string, timestamp string) ([]User, error) {
	contactCol := contactsCollection(db)
	userCol := usersCollection(db)
	ctx := dbContext(30)

	cursor, err := contactCol.Find(ctx, bson.M{}) // find all
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		var contact Contact
		cursor.Decode(&contact)
		if userID == contact.ContactIDOne || userID == contact.ContactIDTwo{
			layout := "2006-01-02T15:04:05.000Z"
			given, err := time.Parse(layout, timestamp)
			
			if err != nil{
				continue
			}
			
			fourteenDaysAgo := given.AddDate(0, 0, -14)
			
			if  contact.TimeStamp.After(fourteenDaysAgo) && contact.TimeStamp.Before(given) {
				var user User
				if userID == contact.ContactIDOne {
					objectID, err := primitive.ObjectIDFromHex(contact.ContactIDTwo)
					if err != nil{
						fmt.Println(err)
					}
					
					err = userCol.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)				
					if err != nil{
						fmt.Println(err)
					}
					users = append(users, user)
				} else {
					objectID, err := primitive.ObjectIDFromHex(userID)
					if err != nil{
						fmt.Println(err)
					}
					
					err = userCol.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)				
					if err != nil{
						fmt.Println(err)
					}
					users = append(users, user)
				}
			}
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users[start : start + count], nil
}

func (contact *Contact) createContact(db *mongo.Database) (*mongo.InsertOneResult, error) {
	contactCol := contactsCollection(db)
	ctx := dbContext(30)
	contact.TimeStamp = time.Now()

	result, err := contactCol.InsertOne(ctx, contact)
	return result, err
}

// helpers
func contactsCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection(contactCollectionName)
}
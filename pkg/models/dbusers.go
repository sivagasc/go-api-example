package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Users is a package level variable acting as an in-memory user database
var DBUsers DBUserCollection

// User represents a user of the system
type DBUser struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	FirstName string             `json:"first_name, omitempty"`
	LastName  string             `json:"last_name, omitempty"`
}

// UserCollection is a collection of user records
type DBUserCollection struct {
	Users []DBUser
}

func AllUsers(collection *mongo.Collection) ([]DBUser, error) {

	fmt.Println("All Users - Get")
	fmt.Println(reflect.TypeOf(collection))

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := collection.Find(context.TODO(), bson.D{})

	// Find() method raised an error
	if err != nil {
		log.Fatal("Finding all documents ERROR:", err)
	} else {
		// iterate over docs using Next()

		for cursor.Next(ctx) {

			usr := DBUser{}
			err := cursor.Decode(&usr)
			if err != nil {
				fmt.Println("cursor.Next() error:", err)
				os.Exit(1)
			}
			DBUsers.AddUser(usr)

		}
	}
	fmt.Println(DBUsers)
	return DBUsers.Users, nil
}

// GetUserByID returns the user record matching privided ID
func (uc DBUserCollection) GetUserByID(id string, collection *mongo.Collection) (*DBUser, error) {

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ObjectID")
	}

	filter := bson.D{{"_id", userID}}
	var result DBUser

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal("Error in fetching:", err)
		return nil, fmt.Errorf("user not found")
	}
	return &result, nil

}
func (uc DBUserCollection) DeleteUserByID(id string, collection *mongo.Collection) (string, error) {

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ObjectID")
	}

	filter := bson.D{{"_id", userID}}

	result, err1 := collection.DeleteMany(context.TODO(), filter)
	if err1 != nil {
		log.Fatal(err1)
		return "Failed to delete!", err1
	}
	fmt.Printf("DeleteMany removed %v document(s)\n", result.DeletedCount)
	if result.DeletedCount == 0 {
		return "No user found", nil
	}
	return "Deleted success", nil
}

func (uc DBUserCollection) CreateUser(user DBUser, collection *mongo.Collection) (string, error) {

	// Insert document into DB
	user.ID = primitive.NewObjectID()
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
		return "error", fmt.Errorf("create user failed")
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	return fmt.Sprintf("Insert operation success!, ID:%s", user.ID.Hex()), nil
}

func (uc *DBUserCollection) AddUser(u DBUser) (*DBUser, error) {
	// nextID := len(uc.Users) + 1 // ID begins with 1
	// u.ID = nextID
	for _, y := range uc.Users {
		if y.FirstName == u.FirstName && y.LastName == y.LastName {
			// Not yet supporting multiple users of same name
			return nil, fmt.Errorf("user with that name already exists")
		}
	}
	// u.CreatedAt = &[]time.Time{time.Now().UTC()}[0]
	// u.UpdatedAt = &[]time.Time{time.Now().UTC()}[0]
	uc.Users = append(uc.Users, u)
	return &u, nil
}

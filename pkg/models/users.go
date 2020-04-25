package models

import (
	"context"
	"fmt"
	"log"
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
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
}

// UserCollection is a collection of user records
type DBUserCollection struct {
	Users []DBUser
}

func AllUsers(collection *mongo.Collection) ([]DBUser, error) {

	fmt.Println("All Users - Get")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := collection.Find(context.TODO(), bson.D{})

	// Find() method raised an error
	if err != nil {
		log.Println("Finding all documents ERROR:", err)
		return nil, err
	} else {
		clear(&DBUsers)
		// iterate over docs using Next()

		for cursor.Next(ctx) {

			usr := DBUser{}
			err := cursor.Decode(&usr)
			if err != nil {
				fmt.Println("cursor.Next() error:", err)
				return nil, err
			}
			DBUsers.AddUser(usr)

		}
	}
	return DBUsers.Users, nil
}

// GetUserByID returns the user record matching privided ID
func (uc DBUserCollection) GetUserByID(id string, collection *mongo.Collection) (*DBUser, error) {

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ObjectID")
		return nil, fmt.Errorf("Invalid Object ID")
	}

	filter := bson.D{{"_id", userID}}
	var result DBUser

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println("Error in fetching:", err)
		return nil, fmt.Errorf("user not found")
	}
	return &result, nil

}
func (uc DBUserCollection) DeleteUserByID(id string, collection *mongo.Collection) (string, error) {

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid ObjectID")
		return "", fmt.Errorf("Invalid Object ID")
	}

	filter := bson.D{{"_id", userID}}

	result, err1 := collection.DeleteMany(context.TODO(), filter)
	if err1 != nil {
		log.Println(err1)
		return "", fmt.Errorf("Failed to delete: %s", err1)
	}

	fmt.Printf("DeleteMany removed %v document(s)\n", result.DeletedCount)
	if result.DeletedCount == 0 {
		return "", fmt.Errorf("No user found")
	}
	return "User Details Deleted successfully", nil
}

func (uc DBUserCollection) CreateUser(user DBUser, collection *mongo.Collection) (string, error) {

	// Insert document into DB
	user.ID = primitive.NewObjectID()

	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("Error in user creation:%s", err)
	}

	fmt.Println("Created a new User, user ID: ", insertResult.InsertedID)

	return fmt.Sprintf("Created a new User, user ID::%s", user.ID.Hex()), nil
}

func (uc *DBUserCollection) AddUser(u DBUser) (*DBUser, error) {
	// nextID := len(uc.Users) + 1 // ID begins with 1
	// u.ID = nextID
	for _, y := range uc.Users {
		if y.FirstName == u.FirstName && y.LastName == u.LastName {
			// Not yet supporting multiple users of same name
			return nil, fmt.Errorf("user with that name already exists")
		}
	}
	// u.CreatedAt = &[]time.Time{time.Now().UTC()}[0]
	// u.UpdatedAt = &[]time.Time{time.Now().UTC()}[0]
	uc.Users = append(uc.Users, u)
	return &u, nil
}

// Clear the interface values
func clear(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

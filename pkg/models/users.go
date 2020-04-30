package models

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/sivagasc/go-api-example/pkg/common"

	"github.com/sivagasc/go-api-example/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Users is a package level variable acting as an in-memory user database
var Users UserCollection

// User represents a user of the system
type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	UserName  string             `json:"username"`
	Password  string             `json:"password"`
}

// UserCollection is a collection of user records
type UserCollection struct {
	Users []User
}

// AllUsers is a public method to return all the user details from the database
func (uc UserCollection) AllUsers() (*UserCollection, error) {

	// Get Logger
	logger := common.GetLoggerInstance()
	// Get DB Connection
	collection, err := common.GetDBConnection()
	if err != nil {
		return nil, fmt.Errorf("Error in getting DB instance")
	}

	logger.Info().Msg("All Users - Get")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := collection.Find(context.TODO(), bson.D{})

	// Find() method raised an error
	if err != nil {
		logger.Error().Msgf("Finding all documents ERROR: %s", err.Error())
		return nil, err
	}
	clear(&Users)
	// iterate over docs using Next()

	for cursor.Next(ctx) {

		usr := User{}
		err := cursor.Decode(&usr)
		if err != nil {
			logger.Error().Msgf("cursor.Next() error: %s", err.Error())
			return nil, err
		}
		Users.addUser(usr)

	}

	return &Users, nil
}

// GetUserByID returns the user record matching privided ID
func (uc UserCollection) GetUserByID(id string) (*User, error) {

	// Get Logger
	logger := common.GetLoggerInstance()
	// Get DB Connection
	collection, err := common.GetDBConnection()
	if err != nil {
		return nil, fmt.Errorf("Error in getting DB instance")
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Info().Msg("Invalid ObjectID")
		return nil, fmt.Errorf("Invalid Object ID")
	}

	filter := bson.D{{"_id", userID}}
	var result User

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		logger.Info().Msgf("Error in fetching: %s", err.Error())
		return nil, fmt.Errorf("user not found")
	}
	return &result, nil

}

//GetUserByUserName returns the user record matching privided ID
func GetUserByUserName(username string) (*User, error) {

	// Get Logger
	logger := common.GetLoggerInstance()
	// Get DB Connection
	collection, err := common.GetDBConnection()
	if err != nil {
		return nil, fmt.Errorf("Error in getting DB instance")
	}

	filter := bson.D{{"username", username}}
	var result User

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		logger.Info().Msgf("Error in fetching: %s", err.Error())
		return nil, fmt.Errorf("user not found")
	}
	return &result, nil

}

// DeleteUserByID is a public method to remove tge user record matching privided ID
func (uc UserCollection) DeleteUserByID(id string) (string, error) {
	// Get Logger
	logger := common.GetLoggerInstance()
	// Get DB Connection
	collection, err := common.GetDBConnection()
	if err != nil {
		return "", fmt.Errorf("Error in getting DB instance")
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Info().Msg("Invalid ObjectID")
		return "", fmt.Errorf("Invalid Object ID")
	}

	filter := bson.D{{"_id", userID}}

	result, err1 := collection.DeleteMany(context.TODO(), filter)
	if err1 != nil {
		logger.Info().Msg(err1.Error())
		return "", fmt.Errorf("Failed to delete: %s", err1)
	}

	logger.Info().Msgf("DeleteMany removed %v document(s)\n", result.DeletedCount)
	if result.DeletedCount == 0 {
		return "", fmt.Errorf("No user found")
	}
	return "User Details Deleted successfully", nil
}

// CreateUser is a public method to create new user in the database
func (uc UserCollection) CreateUser(user *User) (*User, error) {
	// Get Logger
	logger := common.GetLoggerInstance()
	// Get DB Connection
	collection, err := common.GetDBConnection()
	if err != nil {
		return nil, fmt.Errorf("Error in getting DB instance")
	}

	// Insert document into DB
	user.ID = primitive.NewObjectID()
	// Encrypt the password
	encryptPsw, err := utils.EncryptPassword(user.Password)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}
	user.Password = encryptPsw
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, fmt.Errorf("Error in user creation:%s", err)
	}

	logger.Info().Msgf("Created a new User, user ID: %v", insertResult.InsertedID)

	if err != nil {
		logger.Error().Msg(err.Error())
		return nil, err
	}

	return user, nil
}

func (uc *UserCollection) addUser(u User) (*User, error) {
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

package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var UserDBCollection *mongo.Collection

func init() {
	fmt.Println("Inside init...")
	Client = InitDB("mongodb+srv://siriusbot:siriusbot@cluster0-2rdzx.mongodb.net/sample?retryWrites=true")
	UserDBCollection = Client.Database("sample").Collection("users")
}

func InitDB(connectionString string) *mongo.Client {

	fmt.Println("Connecting to Mongo DB....")
	Client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// defer client.Disconnect(ctx)

	fmt.Println("connected to mongoDB")
	return Client
}

func DisconnectDB() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer Client.Disconnect(ctx)
	fmt.Println("DB disconnected!")
}

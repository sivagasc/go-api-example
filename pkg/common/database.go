package common

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectToDB(connectionString, databaseName, collectionName string) (*mongo.Collection, error) {

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
	collection := Client.Database(databaseName).Collection(collectionName)

	return collection, nil
}

func DisconnectDB() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer Client.Disconnect(ctx)
	fmt.Println("DB disconnected!")
}

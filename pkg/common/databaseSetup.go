package common

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Client variable stores the mongoDB Connection
var Client *mongo.Client

//Collection variable stores the mongoDB collection
var Collection *mongo.Collection

//ConnectToDB Method connect to MongoDB
func ConnectToDB(connectionString, databaseName, collectionName string) error {

	logger.Info().Msg("Connecting to Mongo DB....")
	var err error
	Client, err = mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Println(err)
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = Client.Connect(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	// defer client.Disconnect(ctx)

	logger.Info().Msg("connected to mongoDB")
	Collection = Client.Database(databaseName).Collection(collectionName)

	return nil
}

//GetDBConnection Get the Database connection
func GetDBConnection() (*mongo.Collection, error) {
	logger.Debug().Msg("Get DB Connection..")
	// Check the connection
	err := Client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}
	return Collection, nil
}

//DisconnectDB to disconnect the DB connection
func DisconnectDB() {
	logger.Debug().Msg("Disconnect DB Connection..")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer Client.Disconnect(ctx)
	fmt.Println("DB disconnected!")
}

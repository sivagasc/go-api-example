package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/m/pkg/auth"
	"example.com/m/pkg/controllers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {

	// Load ENV file
	viper.SetConfigFile(".env")

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	// if we type assert to other type it will throw an error
	dbURL, ok := viper.Get("DB_URL").(string)
	if !ok {
		dbURL = "" // assign dummy value
		fmt.Println(dbURL)
	}

	// if we type assert to other type it will throw an error
	dbName, ok := viper.Get("DATABASE_NAME").(string)
	if !ok {
		dbName = "" // assign dummy value
		fmt.Println(dbName)
	}

	// if we type assert to other type it will throw an error
	collectionName, ok := viper.Get("COLLECTION_NAME").(string)
	if !ok {
		collectionName = "" // assign dummy value
		fmt.Println(collectionName)
	}

	// Connect to database
	controllers.Connect_database(dbURL, dbName, collectionName)

	// Read API key from command line flag if provided.
	var apiKey string
	flag.StringVar(&apiKey, "apikey", "", "API key")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", controllers.Hello)
	r.HandleFunc("/api/v1/users", controllers.ListUsers)
	r.HandleFunc("/api/v1/users/{id}", controllers.GetUser)
	r.HandleFunc("/api/v1/dbusers", controllers.Get_AllUsers).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/dbusers", controllers.Create_Users).Methods(http.MethodPost)

	r.HandleFunc("/api/v1/dbusers/{id}", controllers.Get_User).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/dbusers/{id}", controllers.Delete_User).Methods(http.MethodDelete)

	// The simple API key security is optional.
	// If a key is provided, we will protect all routes containing "/api/".
	log.Println("API Key:", apiKey)
	if apiKey != "" {
		akm := auth.APIKeyMiddleware{Path: "/api/"}
		akm.InitializeKey(apiKey)
		r.Use(akm.Middleware)
	}
	log.Println("Server is running on :8090")
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(":8090", loggedRouter))
}

package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sivagasc/go-api-example/pkg/auth"
	"github.com/sivagasc/go-api-example/pkg/controllers"
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

	dbURL, ok := viper.Get("DB_URL").(string)
	if !ok {
		log.Fatalln("DB_URL missing in env file")
		os.Exit(1)
	}

	dbName, ok := viper.Get("DATABASE_NAME").(string)
	if !ok {
		log.Fatalln("DATABASE_NAME missing in env file")
		os.Exit(1)
	}

	collectionName, ok := viper.Get("COLLECTION_NAME").(string)
	if !ok {
		log.Fatalln("COLLECTION_NAME missing in env file")
		os.Exit(1)
	}

	// Connect to database
	controllers.Connect_database(dbURL, dbName, collectionName)

	// Read API key from command line flag if provided.
	var apiKey string
	flag.StringVar(&apiKey, "apikey", "", "API key")
	flag.Parse()

	r := mux.NewRouter()

	s := r.PathPrefix("/api/v1").Subrouter()

	r.HandleFunc("/", controllers.Hello)

	s.HandleFunc("/users", controllers.ListUsers)
	s.HandleFunc("/users/{id}", controllers.GetUser)
	s.HandleFunc("/dbusers", controllers.Get_AllUsers).Methods(http.MethodGet)
	s.HandleFunc("/dbusers", controllers.Create_Users).Methods(http.MethodPost)

	s.HandleFunc("/dbusers/{id}", controllers.Get_User).Methods(http.MethodGet)
	s.HandleFunc("/dbusers/{id}", controllers.Delete_User).Methods(http.MethodDelete)

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

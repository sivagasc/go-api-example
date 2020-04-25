package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sivagasc/go-api-example/pkg/auth"
	"github.com/sivagasc/go-api-example/pkg/common"
	"github.com/sivagasc/go-api-example/pkg/controllers"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ENV_FILE        string = ".env"
	DB_URL          string = "DB_URL"
	DATABASE_NAME   string = "DATABASE_NAME"
	COLLECTION_NAME string = "COLLECTION_NAME"
	API_KEY         string = "apikey"
	API_KEY_DOC     string = "API key"
	API_PATH_PREFIX string = "/api/v1"
)

var collection *mongo.Collection

func main() {

	// Load ENV file
	viper.SetConfigFile(ENV_FILE)

	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	var dbURL, dbName, collectionName string
	var ok bool

	if dbURL, ok = viper.Get(DB_URL).(string); !ok {
		log.Fatalln("DB_URL missing in env file")
	}

	if dbName, ok = viper.Get(DATABASE_NAME).(string); !ok {
		log.Fatalln("DATABASE_NAME missing in env file")
	}

	if collectionName, ok = viper.Get(COLLECTION_NAME).(string); !ok {
		log.Fatalln("COLLECTION_NAME missing in env file")
	}

	// Connect to database
	collection, err = common.ConnectToDB(dbURL, dbName, collectionName)
	if err != nil {
		log.Fatalln("Error in DB Connection")
	}

	// Read API key from command line flag if provided.
	var apiKey string
	flag.StringVar(&apiKey, API_KEY, "", API_KEY_DOC)
	flag.Parse()

	r := mux.NewRouter()
	s := r.PathPrefix(API_PATH_PREFIX).Subrouter()

	r.HandleFunc("/", controllers.Hello)
	// sample request
	s.HandleFunc("/users", controllers.ListUsers)
	s.HandleFunc("/users/{id}", controllers.GetUser)
	// request with mongoDB CRUD operation
	s.Handle("/dbusers", controllers.Get_AllUsers(collection)).Methods(http.MethodGet)
	s.Handle("/dbusers", controllers.Create_Users(collection)).Methods(http.MethodPost)
	s.Handle("/dbusers/{id}", controllers.Get_User(collection)).Methods(http.MethodGet)
	s.Handle("/dbusers/{id}", controllers.Delete_User(collection)).Methods(http.MethodDelete)

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

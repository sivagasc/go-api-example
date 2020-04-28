package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sivagasc/go-api-example/pkg/auth"
	"github.com/sivagasc/go-api-example/pkg/common"
	chandler "github.com/sivagasc/go-api-example/pkg/handlers"
	"github.com/sivagasc/go-api-example/pkg/services"
	"github.com/sivagasc/go-api-example/pkg/services/users"

	// lr "github.com/sivagasc/go-api-example/pkg/util/logger"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ENV_FILE        string = ".env"
	DB_URL          string = "DB_URL"
	DATABASE_NAME   string = "DATABASE_NAME"
	COLLECTION_NAME string = "COLLECTION_NAME"
	ENVIRONMENT     string = "ENVIRONMENT"
	API_KEY         string = "apikey"
	API_KEY_DOC     string = "API key"
	API_PATH_PREFIX string = "/api/v1"
	DEVELOPMENT_ENV string = "Development"
	PRODUCTION_ENV  string = "Production"
	// LOG_LEVEL       string = "LOG_LEVEL"
	// LOG_LEVEL_INFO  string = "Info"
	// LOG_LEVEL_DEBUG string = "Debug"
	// LOG_LEVEL_ERROR string = "Error"
)

var collection *mongo.Collection

func main() {

	// Load ENV file
	viper.SetConfigFile(ENV_FILE)

	// Logger Config
	logger := &log.Logger
	// logger := lr.New(true)

	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal().Msgf("Error while reading config file: %s", err.Error())
	}

	var dbURL, dbName, collectionName, env string
	var ok bool

	if env, ok = viper.Get(ENVIRONMENT).(string); !ok {
		logger.Fatal().Msg("Env missing in env file")
	}
	if env == PRODUCTION_ENV {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger.Info().Msg("*** Production Configuration ***")
	} else {
		zl := logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		logger = &zl
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger.Info().Msg("*** Non-production Configuration ***")
		logger.Debug().Msg("*** Debug Logging Enabled ***")
	}

	if dbURL, ok = viper.Get(DB_URL).(string); !ok {
		logger.Fatal().Msg("DB_URL missing in env file")
	}

	if dbName, ok = viper.Get(DATABASE_NAME).(string); !ok {
		logger.Fatal().Msg("DATABASE_NAME missing in env file")
	}

	if collectionName, ok = viper.Get(COLLECTION_NAME).(string); !ok {
		logger.Fatal().Msg("COLLECTION_NAME missing in env file")
	}

	// Connect to database
	collection, err = common.ConnectToDB(dbURL, dbName, collectionName)
	if err != nil {
		logger.Fatal().Msg("Error in DB Connection")
	}

	srvcEnv := &services.Env{
		Collection: collection, // Shared database connection goes here
		Log:        logger,
	}

	// Service interface
	usersService, err := users.NewUsersSvc(collection, logger)

	// Read API key from command line flag if provided.
	var apiKey string
	flag.StringVar(&apiKey, API_KEY, "", API_KEY_DOC)
	flag.Parse()

	r := mux.NewRouter()
	s := r.PathPrefix(API_PATH_PREFIX).Subrouter()

	r.Handle("/", chandler.Hello(srvcEnv))
	// request with mongoDB CRUD operation
	s.Handle("/users", chandler.GetAllUsers(srvcEnv, usersService)).Methods(http.MethodGet)
	s.Handle("/users", chandler.CreateUsers(srvcEnv, usersService)).Methods(http.MethodPost)
	s.Handle("/users/{id}", chandler.GetUser(srvcEnv, usersService)).Methods(http.MethodGet)
	s.Handle("/users/{id}", chandler.DeleteUser(srvcEnv, usersService)).Methods(http.MethodDelete)

	// Authentication
	a := r.PathPrefix(API_PATH_PREFIX + "/auth").Subrouter()
	a.Handle("/login", chandler.TokenAuth(srvcEnv, usersService)).Methods(http.MethodGet)
	a.Handle("/refresh-token", chandler.RefreshToken(srvcEnv)).Methods(http.MethodGet)
	a.Handle("/logout", chandler.Logout(srvcEnv)).Methods(http.MethodGet)

	// The simple API key security is optional.
	// If a key is provided, we will protect all routes containing "/api/".
	logger.Debug().Msg("API Key:" + apiKey)
	if apiKey != "" {
		akm := auth.APIKeyMiddleware{Path: "/api/"}
		akm.InitializeKey(apiKey)
		r.Use(akm.Middleware)
	}

	logger.Info().Msg("Server is running on :8090")
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	logger.Fatal().Msg(http.ListenAndServe(":8090", loggedRouter).Error())
}

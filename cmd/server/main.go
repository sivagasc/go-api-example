package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/sivagasc/go-api-example/pkg/auth"
	"github.com/sivagasc/go-api-example/pkg/common"
	"github.com/sivagasc/go-api-example/pkg/config"
	chandler "github.com/sivagasc/go-api-example/pkg/handlers"
	"github.com/sivagasc/go-api-example/pkg/services/users"

	"github.com/spf13/viper"
)

//Common const
const (
	EnvFile        string = "env"
	DatabaseURL    string = "DB_URL"
	DatabaseName   string = "DATABASE_NAME"
	CollectionName string = "COLLECTION_NAME"
	APIKey         string = "apikey"
	APIKeyDoc      string = "API key"
	APIPathPrefix  string = "/api/v1"
)

func getEnvData(filename string) (config.Configurations, error) {

	viper.SetConfigName(filename)
	viper.AddConfigPath(".")
	var configuration config.Configurations

	if err := viper.ReadInConfig(); err != nil {
		return configuration, err
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		return configuration, err
	}
	return configuration, nil
}

func main() {

	// Logger Config
	logger := &log.Logger

	// Load Env Details
	var envConfig config.Configurations

	envConfig, err := getEnvData(EnvFile)
	if err != nil {
		logger.Fatal().Msg("Error in reading Env file")
	}

	var dbURL, dbName, collectionName, env, logOutput string

	if env = envConfig.Server.Environment; env == "" {
		logger.Fatal().Msg("Env missing in env file")
	}
	if logOutput = envConfig.Logger.OutputPath; logOutput == "" {
		logger.Fatal().Msg("Env missing in env file")
	}

	// Load Custom Logger
	logger = common.SetupLoggerInstance(logOutput, env)

	if dbURL = envConfig.Database.URL; dbURL == "" {
		logger.Fatal().Msg("DB_URL missing in env file")
	}

	if dbName = envConfig.Database.DBName; dbName == "" {
		logger.Fatal().Msg("DATABASE_NAME missing in env file")
	}

	if collectionName = envConfig.Database.CollectionName; collectionName == "" {
		logger.Fatal().Msg("COLLECTION_NAME missing in env file")
	}

	// Connect to database
	err = common.ConnectToDB(dbURL, dbName, collectionName)
	if err != nil {
		logger.Fatal().Msg("Error in DB Connection")
	}

	// Service interface
	usersService, err := users.NewUsersSvc()

	r := mux.NewRouter()
	s := r.PathPrefix(APIPathPrefix).Subrouter()
	r.Handle("/", chandler.Hello())

	// request with mongoDB CRUD operation
	s.Handle("/users", chandler.GetAllUsers(usersService)).Methods(http.MethodGet)
	s.Handle("/users", chandler.CreateUsers(usersService)).Methods(http.MethodPost)
	s.Handle("/users/{id}", chandler.GetUser(usersService)).Methods(http.MethodGet)
	s.Handle("/users/{id}", chandler.DeleteUser(usersService)).Methods(http.MethodDelete)

	// Authentication
	a := r.PathPrefix(APIPathPrefix + "/auth").Subrouter()
	a.Handle("/login", chandler.TokenAuth(usersService)).Methods(http.MethodGet)
	a.Handle("/refresh-token", chandler.RefreshToken()).Methods(http.MethodGet)
	a.Handle("/logout", chandler.Logout()).Methods(http.MethodGet)

	// Add middleware authentication check
	akm := auth.APIKeyMiddleware{Path: "/api/v1/users"}
	r.Use(akm.Middleware)

	logger.Info().Msg("Server is running on :8090")
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	logger.Fatal().Msg(http.ListenAndServe(":8090", loggedRouter).Error())
}

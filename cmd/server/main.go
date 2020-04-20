package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kkeuning/go-api-example/pkg/auth"
	"github.com/kkeuning/go-api-example/pkg/models"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<h1>Hello, World!</h1>\n")
}

func getUser(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id != "" {
		userID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected id as an integer.")
			return
		}
		u, err := models.Users.GetUserByID(userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "User not found.")
			return
		}
		juser, err := json.Marshal(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Something went wrong.")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(juser)
		return
	}
}

func listUsers(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id != "" {
		userID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected id as an integer.")
			return
		}
		for _, u := range models.Users.GetUsers() {
			if u.ID == userID {
				juser, err := json.Marshal(u)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Something went wrong.")
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write(juser)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User not found.")
		return
	}
	// No id specified, list all users.
	jusers, err := json.Marshal(models.Users.GetUsers())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something went wrong.")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jusers)
	return
}

func main() {
	// Read API key from command line flag if provided.
	var apiKey string
	flag.StringVar(&apiKey, "apikey", "", "API key")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", hello)
	r.HandleFunc("/api/v1/users", listUsers)
	r.HandleFunc("/api/v1/users/{id}", getUser)

	// The simple API key security is optional.
	// If a key is provided, we will protect all routes containing "/api/".
	if apiKey != "" {
		akm := auth.APIKeyMiddleware{Path: "/api/"}
		akm.InitializeKey(apiKey)
		r.Use(akm.Middleware)
	}

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(":8090", loggedRouter))
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"example.com/m/pkg/auth"
	"example.com/m/pkg/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

// Database handler

func get_AllUsers(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	allUsers, err := models.AllUsers(models.UserDBCollection)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error in all Users return.")
		return
	}
	usersJson, err := json.Marshal(allUsers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Cannot encode to JSON ", err)
		return
	}
	fmt.Fprintf(os.Stdout, "%s", usersJson)

	w.Write([]byte(usersJson))

	return
}

func get_User(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Expected id as an input.")
		return
	}

	u, err := models.DBUsers.GetUserByID(id, models.UserDBCollection)
	fmt.Println("Here")
	if err != nil {
		fmt.Println("inside if")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User not found.")
		return
	}
	fmt.Println("Not returned yet")
	juser, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something went wrong.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(juser)
	return

}

func delete_User(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Expected id as an input.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	message, err := models.DBUsers.DeleteUserByID(id, models.UserDBCollection)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something went wrong in user deletion.")
		return
	}

	w.Write([]byte(`{"message":"` + message + `"}`))
	return
}

func create_Users(w http.ResponseWriter, req *http.Request) {
	var user models.DBUser
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := models.DBUsers.CreateUser(user, models.UserDBCollection)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something went wrong in user deletion.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"` + message + `"}`))
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
	r.HandleFunc("/api/v1/dbusers", get_AllUsers).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/dbusers", create_Users).Methods(http.MethodPost)

	r.HandleFunc("/api/v1/dbusers/{id}", get_User).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/dbusers/{id}", delete_User).Methods(http.MethodDelete)

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

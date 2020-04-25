// Package classification User Operations API
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// Documentation for User Operations API
//
// Schemas: http
// BasePath: /v1
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// -application/json
// swagger:meta

package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sivagasc/go-api-example/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func Get_AllUsers(collection *mongo.Collection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		allUsers, err := models.AllUsers(collection)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
			return
		}

		usersJson, err := json.Marshal(allUsers)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal("Cannot encode to JSON ", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(usersJson))
		return
	})
}

func Get_User(collection *mongo.Collection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id := mux.Vars(req)["id"]
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected id as an input.")
			return
		}

		u, err := models.DBUsers.GetUserByID(id, collection)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, err.Error())
			return
		}

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
	})
}

func Delete_User(collection *mongo.Collection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id := mux.Vars(req)["id"]
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected id as an input.")
			return
		}

		message, err := models.DBUsers.DeleteUserByID(id, collection)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: "+err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"` + message + `"}`))
		return
	})
}

func Create_Users(collection *mongo.Collection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var user models.DBUser

		err := json.NewDecoder(req.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		message, err := models.DBUsers.CreateUser(user, collection)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"` + message + `"}`))
	})
}

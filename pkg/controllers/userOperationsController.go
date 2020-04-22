package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/m/pkg/common"
	"example.com/m/pkg/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserDBCollection *mongo.Collection

func Get_AllUsers(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	allUsers, err := models.AllUsers(UserDBCollection)

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

func Get_User(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Expected id as an input.")
		return
	}

	u, err := models.DBUsers.GetUserByID(id, UserDBCollection)

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(juser)
	return

}

func Delete_User(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Expected id as an input.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	message, err := models.DBUsers.DeleteUserByID(id, UserDBCollection)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something went wrong in user deletion.")
		return
	}

	w.Write([]byte(`{"message":"` + message + `"}`))
	return
}

func Create_Users(w http.ResponseWriter, req *http.Request) {
	var user models.DBUser
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := models.DBUsers.CreateUser(user, UserDBCollection)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something went wrong in user deletion.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"` + message + `"}`))
}

func Connect_database(dbURL, dbName, collectionName string) {
	UserDBCollection = common.ConnectToDB(dbURL, dbName, collectionName)
	fmt.Printf("DB Connection established")
}

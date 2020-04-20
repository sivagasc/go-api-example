package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

var users = []user{
	{
		ID:        1,
		FirstName: "Rob",
		LastName:  "Pike",
	},
	{
		ID:        2,
		FirstName: "Ken",
		LastName:  "Thompson",
	},
	{
		ID:        3,
		FirstName: "Robert",
		LastName:  "Griesemer",
	},
}

type user struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<h1>Hello, World!</h1>\n")
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
		for _, u := range users {
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
	jusers, err := json.Marshal(users)
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
	// Initial example is using DefaultServeMux
	// Will be updated to use gorilla/mux later to add path variables
	http.HandleFunc("/", hello)
	http.HandleFunc("/api/v1/users", listUsers)
	http.ListenAndServe(":8090", nil)
}

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/sivagasc/go-api-example/pkg/models"
	"github.com/sivagasc/go-api-example/pkg/services"
	"github.com/sivagasc/go-api-example/pkg/services/users"
)

// Hello is a simple Hello handler method
func Hello(e *services.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "<h1>Hello, World!</h1>\n")
	})
}
func validateToken(tokenString string) {
	// Token from another example.  This token is expired
	// var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"
	tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	if token.Valid {
		fmt.Println("You look nice today")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}

	// Output: Timing is everything
}

// GetAllUsers method used to retrieve all the user details from the database
func GetAllUsers(env *services.Env, usersSvc users.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		env.Log.Info().Msg("*** Get All Users")

		//Set authorization key in context
		ctx := context.WithValue(req.Context(), services.AuthorizationKey, req.Header.Get("Authorization"))

		payload := &users.ListPayload{}

		allUsers, err := usersSvc.List(ctx, payload)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
			return
		}

		usersJSON, err := json.Marshal(allUsers.Users)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			env.Log.Info().Msgf("Cannot encode to JSON, Error: %s", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(usersJSON))
		return
	})
}

// GetUser method is used to get an invidividual user from a database
func GetUser(env *services.Env, usersSvc users.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), services.AcceptTypeKey, req.Header.Get("Accept"))

		id := mux.Vars(req)["id"]
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			env.Log.Info().Msg("Expected id as an input.")
			return
		}

		payload := &users.ShowPayload{
			ID: id,
		}
		u, err := usersSvc.Show(ctx, payload)
		// u, err := models.Users.GetUserByID(id, env.Collection)
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

// DeleteUser method used to delete a user from a database
// ID must be added into the request param
func DeleteUser(env *services.Env, usersSvc users.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), services.AcceptTypeKey, req.Header.Get("Accept"))

		id := mux.Vars(req)["id"]
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected id as an input.")
			return
		}
		payload := &users.DeletePayload{
			ID: id,
		}
		message, err := usersSvc.Delete(ctx, payload)

		// message, err := models.Users.DeleteUserByID(id, env.Collection)
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

// CreateUsers method is used to Create a new User in the database
// The user details should be passed through request param
// Format {'firstname':'...','lastname':'...'}
func CreateUsers(env *services.Env, usersSvc users.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), services.AcceptTypeKey, req.Header.Get("Accept"))

		var user *models.User

		err := json.NewDecoder(req.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		user, err = usersSvc.Create(ctx, user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
			return
		}

		juser, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Something went wrong.")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(juser)
	})
}

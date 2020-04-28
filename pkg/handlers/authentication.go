package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sivagasc/go-api-example/pkg/auth"
	"github.com/sivagasc/go-api-example/pkg/models"
	"github.com/sivagasc/go-api-example/pkg/services"
	"github.com/sivagasc/go-api-example/pkg/services/users"
)

// TokenAuth is a JWT Authentication method
func TokenAuth(env *services.Env, usersSvc users.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestUser := new(models.UserAuthentication)
		decoder := json.NewDecoder(req.Body)
		decoder.Decode(&requestUser)
		fmt.Println(requestUser)
		responseStatus, token, expirationTime := auth.Login(requestUser, env.Collection, env.Log)
		w.WriteHeader(responseStatus)

		// Setting token and expiration time in cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: expirationTime,
		})

		if responseStatus == http.StatusOK {
			w.Header().Set("Content-Type", "application/json")
			response := map[string]string{"status": "success", "token": token}
			resJSON, _ := json.Marshal(response)
			w.Write([]byte(resJSON))
		}

		return
	})
}

// RefreshToken is to JWT User Refresh session
func RefreshToken(e *services.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "<h1>Refresh token</h1>\n")
	})
}

// Logout is to kill the user session
func Logout(e *services.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "<h1>Logout</h1>\n")
	})
}

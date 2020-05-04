package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sivagasc/go-api-example/pkg/auth"
	"github.com/sivagasc/go-api-example/pkg/models"
	"github.com/sivagasc/go-api-example/pkg/services/users"
	"github.com/sivagasc/go-api-example/pkg/utils"
)

// TokenAuth is a JWT Authentication method
func TokenAuth(usersSvc users.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestUser := new(models.UserAuthentication)
		decoder := json.NewDecoder(req.Body)
		decoder.Decode(&requestUser)
		responseStatus, token, expirationTime := auth.Login(requestUser)

		// Setting token and expiration time in cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: expirationTime,
		})

		if responseStatus == http.StatusOK {
			response := map[string]string{"status": "success", "token": token}
			utils.RespondJSON(w, responseStatus, response)
			return
		}

		utils.RespondJSON(w, responseStatus, "")
		return
	})
}

// RefreshToken is to JWT User Refresh session
func RefreshToken() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "<h1>Refresh token</h1>\n")
	})
}

// Logout is to kill the user session
func Logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "<h1>Logout</h1>\n")
	})
}

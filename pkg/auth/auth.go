package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/sivagasc/go-api-example/pkg/common"
	"github.com/sivagasc/go-api-example/pkg/models"
)

// APIKeyMiddleware is a simplified shared secret authentication for api requests
type APIKeyMiddleware struct {
	apiKey []byte
	Path   string
}

// // InitializeKey will load a hash that could be safely persisted instead of the actual key
// func (akm *APIKeyMiddleware) InitializeKey(key string) {
// 	k := []byte(key)
// 	bCryptKey, err := bcrypt.GenerateFromPassword(k, bcrypt.DefaultCost)
// 	if err != nil {
// 		log.Fatal("Key initialization error.")
// 	}
// 	akm.apiKey = bCryptKey
// }

// KeyIsValid checks that userAPIKey matches userAPIKey hash
// func (akm *APIKeyMiddleware) KeyIsValid(userAPIKey string) bool {
// 	log.Println("userkeystring:", userAPIKey)
// 	// if e := bcrypt.CompareHashAndPassword(akm.apiKey, []byte(userAPIKey)); e != nil {
// 	// 	log.Println("Key Mismatch")
// 	// 	return false
// 	// }

// 	return true
// }

// Middleware function, which will be called for each api request
func (akm *APIKeyMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if akm.Path == "" || strings.Contains(r.URL.Path, akm.Path) {
			// Get Logger
			logger := common.GetLoggerInstance()

			logger.Info().Msgf("Authorization required for %s", r.URL.Path)
			key := r.Header.Get("Authorization")

			userAuth := InitJWTAuthentication()
			var httpCode int
			var flag bool

			if httpCode, flag = userAuth.ValidateToken(key); flag { //akm.KeyIsValid(key)
				logger.Info().Msg("User is authorized.")
				// Pass down the request to the next middleware (or final handler)
				next.ServeHTTP(w, r)
				return
			}
			// Write an error and return to stop the handler chain
			//http.Error(w, "Forbidden", http.StatusForbidden)
			w.WriteHeader(httpCode)
			return
		}
		// Auth not required for non-api endpoints
		// Just pass the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})
}

//Login method authenticate the user and provide JWT token
func Login(requestUser *models.UserAuthentication) (int, string, time.Time) {
	// Get Logger
	logger := common.GetLoggerInstance()

	userAuth := InitJWTAuthentication()
	if userAuth.Authenticate(requestUser) {
		token, expirationTime, err := userAuth.GenerateToken(requestUser.Username)
		if err != nil {
			logger.Error().Msgf("Error:%s", err)
			return http.StatusInternalServerError, "Internal Server Error", time.Now()
		}
		return http.StatusOK, token, expirationTime
	}
	return http.StatusUnauthorized, "Unauthorized", time.Now()
}

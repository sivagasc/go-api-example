package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog"
	"github.com/sivagasc/go-api-example/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// JWTAuthentication ...
type JWTAuthentication struct {
	privateKey string
	PublicKey  string
}

var jwtKey = []byte("my_secret_key")

const (
	tokenDurationInMinute = 5
	expireOffset          = 3600
)

// Claims a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var authBackendInstance *JWTAuthentication = nil

// InitJWTAuthentication method initialize the JWTAuthentication model
func InitJWTAuthentication() *JWTAuthentication {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthentication{
			privateKey: "my_secret_key", //TODO: Need to change
			PublicKey:  "my_secret_key", //TODO: Need to change
		}
	}
	return authBackendInstance
}

// GenerateToken method generates the User JWT Token
func (backend *JWTAuthentication) GenerateToken(username string) (string, time.Time, error) {

	// Setting token expiration
	expirationTime := time.Now().Add(tokenDurationInMinute * time.Minute)

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //jwt.SigningMethodRS512

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "token error", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

//Authenticate method used to authenticate the users
func (backend *JWTAuthentication) Authenticate(user *models.UserAuthentication, collection *mongo.Collection, logger *zerolog.Logger) bool {
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testing"), 10)
	// testUser := models.UserAuthentication{
	// 	Username: "test",
	// 	Password: string(hashedPassword),
	// }

	dbuser, err := models.GetUserByUserName(user.Username, collection, logger)
	if err != nil {
		logger.Error().Msgf("Error on validating username:%s", err.Error())
		return false
	}
	userauth := models.UserAuthentication{
		Username: dbuser.UserName,
		Password: string(dbuser.Password),
	}

	return user.Username == userauth.Username && bcrypt.CompareHashAndPassword([]byte(userauth.Password), []byte(user.Password)) == nil
}

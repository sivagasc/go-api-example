package auth

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sivagasc/go-api-example/pkg/common"
	"github.com/sivagasc/go-api-example/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// JWTAuthentication ...
type JWTAuthentication struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	jwtKey     []byte
}

// var jwtKey = []byte("my_secret_key")

const (
	tokenDurationInMinute = 15
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
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
			jwtKey:     []byte("my_secret_key"),
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //jwt.SigningMethodRS512 or SigningMethodHS256

	// Create the JWT string
	tokenString, err := token.SignedString(backend.jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "token error", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

//Authenticate method used to authenticate the users
func (backend *JWTAuthentication) Authenticate(user *models.UserAuthentication) bool {
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testing"), 10)
	// testUser := models.UserAuthentication{
	// 	Username: "test",
	// 	Password: string(hashedPassword),
	// }

	// Get Logger
	logger := common.GetLoggerInstance()

	dbuser, err := models.GetUserByUserName(user.Username)
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

// ValidateToken ...
func (backend *JWTAuthentication) ValidateToken(token string) (int, bool) {

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	token = strings.Replace(token, "Bearer ", "", -1)

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return backend.jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {

			return http.StatusUnauthorized, false
		}
		return http.StatusBadRequest, false
	}
	if !tkn.Valid {
		return http.StatusUnauthorized, false
	}
	return http.StatusOK, true
}

func getPrivateKey() *rsa.PrivateKey {
	privateKeyFile, err := os.Open("/Users/sivanandhamp/Documents/Development/PlayGroundWorkspace/Go-Lang/go-api-example-master/keys/private_key")
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open("/Users/sivanandhamp/Documents/Development/PlayGroundWorkspace/Go-Lang/go-api-example-master/keys/public_key.pub")
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}

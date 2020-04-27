package services

import (
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type contextKey int

const (
	// AcceptTypeKey is the context key used to store the value of the HTTP
	// Accept-Type header
	AcceptTypeKey contextKey = iota + 1
	// AuthorizationKey is the context key used to store the value of the HTTP
	// Authorization Header
	AuthorizationKey
)

// Env ...
type Env struct {
	Collection *mongo.Collection
	Log        *zerolog.Logger
}

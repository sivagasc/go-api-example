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

package users

import (
	"context"

	"github.com/sivagasc/go-api-example/pkg/models"
)

type usersSvc struct {
	// collection *mongo.Collection
	// logger     *zerolog.Logger
}

// Service ...
type Service interface {
	Show(context.Context, *ShowPayload) (*models.User, error)
	Delete(context.Context, *DeletePayload) (string, error)
	Update(context.Context, *UpdatePayload) (string, error)
	List(context.Context, *ListPayload) (*models.UserCollection, error)
	Create(context.Context, *models.User) (*models.User, error)
}

// NewUsersSvc ...
func NewUsersSvc() (Service, error) {
	return &usersSvc{}, nil
}

func (us *usersSvc) Show(_ context.Context, payload *ShowPayload) (*models.User, error) {

	u, err := models.Users.GetUserByID(payload.ID)
	if err != nil {
		return nil, err
	}
	return u, nil

}

func (us *usersSvc) Delete(_ context.Context, payload *DeletePayload) (string, error) {

	message, err := models.Users.DeleteUserByID(payload.ID)
	if err != nil {
		return "Error on delete", err
	}
	return message, nil

}

func (us *usersSvc) Update(_ context.Context, _ *UpdatePayload) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (us *usersSvc) List(_ context.Context, _ *ListPayload) (*models.UserCollection, error) {
	allUsers, err := models.Users.AllUsers()
	if err != nil {
		return nil, err
	}
	return allUsers, nil
}

func (us *usersSvc) Create(_ context.Context, user *models.User) (*models.User, error) {

	userCollection, err := models.Users.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return userCollection, nil
}

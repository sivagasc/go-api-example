package users

import "github.com/sivagasc/go-api-example/pkg/models"

// ShowPayload ...
type ShowPayload struct {
	ID string
}

// ListPayload ...
type ListPayload struct {
	ID string
}

// DeletePayload ...
type DeletePayload struct {
	ID string
}

// UpdatePayload ...
type UpdatePayload struct {
	User models.User
}

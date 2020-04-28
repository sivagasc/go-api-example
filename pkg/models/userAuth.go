package models

// UserAuthentication is used to authenticate the user and pass JWT to client
type UserAuthentication struct {
	UUID     string `json:"uuid" form:"-"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

package model

// User is a struct that represents a user
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// NewUser is a function that creates a new user
func NewUser(id int, username, password, email string) User {
	return User{
		ID:       id,
		Username: username,
		Password: password,
		Email:    email,
	}
}

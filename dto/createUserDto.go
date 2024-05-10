package dto

// CreateUserRequestDto is a struct that represents the data that is required to create a user
type CreateUserRequestDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// CreateUserResponseDto is a struct that represents the response data that is returned after creating a user
type CreateUserResponseDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewCreateUserResponseDto(username, email string) CreateUserResponseDto {
	return CreateUserResponseDto{
		Username: username,
		Email:    email,
	}
}

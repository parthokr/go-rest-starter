package handler

import (
	"encoding/json"
	"go-rest-starter/dto"
	"go-rest-starter/server"
	"log"
	"net/http"
)

func HandleCreateUser(w http.ResponseWriter, r *http.Request) error {
	// unmarshal request body
	var user dto.CreateUserRequestDto
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding request body: %v", err)
		// generate error response for the consumer
		err := server.WriteJSON(w, 400, map[string]string{"error": "failed to decode request body"})
		// if error occurs while writing response, return the error (for developer)
		if err != nil {
			return err
		}
		// return nil to indicate that the error has been handled
		// ie, the consumer has been informed about the error
		return nil
	}

	// validate request body
	type fieldError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	validationErrors := make([]fieldError, 0)

	if user.Username == "" {
		validationErrors = append(validationErrors, fieldError{"username", "username is required"})
	}
	if user.Email == "" {
		validationErrors = append(validationErrors, fieldError{"email", "email is required"})
	}
	if user.Password == "" {
		validationErrors = append(validationErrors, fieldError{"password", "password is required"})
	}
	if len(validationErrors) > 0 {
		err := server.WriteJSON(w, 400, map[string]interface{}{"error": "invalid fields", "fields": validationErrors})
		if err != nil {
			return err
		}
		return nil
	}
	// log request body
	log.Printf("Request body: %+v", user)
	// create user
	err := server.WriteJSON(w, 201, dto.NewCreateUserResponseDto(user.Username, user.Email))
	if err != nil {
		return err
	}
	return nil
}

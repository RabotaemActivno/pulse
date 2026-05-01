package dto

import "github.com/google/uuid"

type RegisterUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserOutput struct {
	ID uuid.UUID `json:"id"`
}

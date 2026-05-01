package domain

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email" validate:"email"`
	CreatedAt    time.Time `json:"created_at"`
	PasswordHash string    `json:"password_hash"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func NewUser(email, pass string) (User, error) {

	bytesPass := []byte(pass)
	passwordHashBytes, err := bcrypt.GenerateFromPassword(bytesPass, bcrypt.DefaultCost)

	if err != nil {
		return User{}, fmt.Errorf("password.Hash: %w", err)
	}

	user := User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(passwordHashBytes),
	}

	if err := user.Validate(); err != nil {
		return User{}, fmt.Errorf("p.Validate: %w", err)
	}

	return user, nil
}

func (u User) Validate() error {
	err := validate.Struct(u)
	if err != nil {
		return fmt.Errorf("validate struct: %w", err)
	}

	return nil
}

func (u User) CheckPassword(pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pass))
}

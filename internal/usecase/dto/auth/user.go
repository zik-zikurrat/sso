package auth

import (
	"time"

	"github.com/google/uuid"
)

type GetUser struct {
	ID        uuid.UUID
	Email     string
	Password  string
	Login     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateUserRepo struct {
	Email        string
	PasswordHash []byte
	Login        string
	Role         string
}

type CreateUser struct {
	Email    string
	Password string
	Login    string
	Role     string
}

type UpdateUser struct {
	ID       *uuid.UUID
	Email    *string
	Password *string
	Login    *string
	Role     *string
}

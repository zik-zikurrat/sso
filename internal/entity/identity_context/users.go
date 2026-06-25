package identitycontext

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Login        string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

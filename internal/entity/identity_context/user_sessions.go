package identitycontext

import (
	"time"

	"github.com/google/uuid"
)

type UserSessions struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

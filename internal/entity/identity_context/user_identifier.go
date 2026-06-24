package identitycontext

import "github.com/google/uuid"

type UserIdentifier struct {
	ID    *uuid.UUID
	Login *string
	Email *string
}

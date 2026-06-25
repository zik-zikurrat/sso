package entity

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID        uuid.UUID
	Name      string
	Method    string
	URL       string
	Secure    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ServiceIdentifier struct {
	ID   *uuid.UUID
	Name *string
}

// register all routes
// webhook from all services to SSO to register their routes

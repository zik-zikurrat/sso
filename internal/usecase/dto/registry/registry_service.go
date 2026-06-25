package registry

import (
	"time"

	"github.com/google/uuid"
)

type GetService struct {
	ID        uuid.UUID
	Name      string
	Method    string
	URL       string
	Secure    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateService struct {
	Name     string
	Metadata map[string]string
	// Method   string
	// URL      string
	// Secure   bool
}

type UpdateService struct {
	ID     uuid.UUID
	Name   *string
	Method *string
	URL    *string
	Secure *bool
}

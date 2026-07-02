package registry

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrServiceNotFound = errors.New("service not found")

type Endpoint struct {
	ID        uuid.UUID
	Method    string
	URL       string
	Secure    bool
	CreatedAt time.Time
}

type CreateService struct {
	Name      string
	Endpoints []Endpoint
}

type ServiceWithEndpoints struct {
	ID        uuid.UUID
	Name      string
	Endpoints []Endpoint
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateService struct {
	ID   uuid.UUID
	Name *string
}

type UpdateEndpoint struct {
	ID     uuid.UUID
	Method *string
	URL    *string
	Secure *bool
}

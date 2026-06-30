package registry

import (
	"sso/internal/entity"
	"time"

	"github.com/google/uuid"
)

type Endpoint struct {
	ID        uuid.UUID
	Method    string
	URL       string
	Secure    bool
	CreatedAt time.Time
}

type GetService struct {
	Service   entity.Service
	Endpoints []entity.Endpoint
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

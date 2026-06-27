package registry

import (
	"sso/internal/entity"
	"time"

	"github.com/google/uuid"
)

type Endpoint struct {
	Method string
	URL    string
	Secure bool
}

type GetService struct {
	Service   entity.Service
	Endpoints []entity.Endpoint
}

type CreateService struct {
	Name      string
	Endpoints []Endpoint
}

type ListService struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateService struct {
	ID   uuid.UUID
	Name *string
}

type UpdteEndpoint struct {
	ID     uuid.UUID
	Method *string
	URL    *string
	Secure *bool
}

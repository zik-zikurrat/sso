package entity

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID        uuid.UUID
	Name      string
	Endpoints []Endpoint
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Endpoint struct {
	Method    string
	URL       string
	Secure    bool
	CreatedAt time.Time
}

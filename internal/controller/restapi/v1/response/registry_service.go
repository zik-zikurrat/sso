package response

import (
	"sso/internal/entity"

	"github.com/google/uuid"
)

type Metadata struct {
	Method string `json:"method"`
	URL    string `json:"url"`
}

type Service struct {
	ID       uuid.UUID `json:"id"`
	Metadata Metadata  `json:"metadata"`
}

func ToService(service entity.Service) Service {
	return Service{
		ID: service.ID,
		Metadata: Metadata{
			Method: service.Method,
			URL:    service.URL,
		},
	}
}

package response

import (
	"sso/internal/usecase/dto/registry"
	"time"

	"github.com/google/uuid"
)

type Endpoint struct {
	Method    string    `json:"method"`
	URL       string    `json:"url"`
	Secure    bool      `json:"secure"`
	CreatedAt time.Time `json:"created_at"`
}

type Service struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Endpoints []Endpoint `json:"endpoints"`
}

func ToServiceEndpoints(services []registry.ServiceWithEndpoints) []Service {
	outService := make([]Service, 0, len(services))
	for _, service := range services {
		endpoints := make([]Endpoint, 0, len(service.Endpoints))
		for _, e := range service.Endpoints {
			endpoints = append(endpoints, Endpoint{
				Method:    e.Method,
				URL:       e.URL,
				Secure:    e.Secure,
				CreatedAt: e.CreatedAt,
			})
		}
		outService = append(outService, Service{
			ID:        service.ID,
			Name:      service.Name,
			Endpoints: endpoints,
		})
	}
	return outService
}

type ServiceListItem struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToServiceList(in []registry.ServiceWithEndpoints) []ServiceListItem {
	out := make([]ServiceListItem, 0, len(in))
	for _, s := range in {
		out = append(out, ServiceListItem{
			ID:        s.ID,
			Name:      s.Name,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		})
	}
	return out
}

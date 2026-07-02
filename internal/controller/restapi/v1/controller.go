package v1

import (
	"context"
	"log/slog"
	"sso/internal/usecase/dto/registry"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type RegistryUseCase interface {
	RegisterService(ctx context.Context, in registry.CreateService) (string, error)
	ListServiceEndpoints(ctx context.Context) ([]registry.ServiceWithEndpoints, error)
	GetServiceEndpointsByServiceID(ctx context.Context, in uuid.UUID) (registry.ServiceWithEndpoints, error)
	DeleteService(ctx context.Context, serviceID uuid.UUID) error
	UpdateService(ctx context.Context, in registry.UpdateService) error
}

// V1 -.
type V1 struct {
	l *slog.Logger
	v *validator.Validate

	registry RegistryUseCase
}

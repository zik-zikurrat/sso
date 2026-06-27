package v1

import (
	"context"
	"log/slog"
	"sso/internal/entity"
	"sso/internal/usecase/dto/registry"

	"github.com/go-playground/validator/v10"
)

type RegistryUseCase interface {
	RegisterService(ctx context.Context, in registry.CreateService) (string, error)
	ListService(ctx context.Context) ([]registry.ListService, error)
	GetServiceByID(ctx context.Context, in entity.ServiceIdentifier) (entity.Service, error)
}

// V1 -.
type V1 struct {
	l *slog.Logger
	v *validator.Validate

	registry RegistryUseCase
}

package v1

import (
	"context"
	"log/slog"
	"sso/internal/entity"

	"github.com/go-playground/validator/v10"
)

type RegistryUseCase interface {
	ListService(ctx context.Context) ([]entity.Service, error)
	GetServiceByID(ctx context.Context, in entity.ServiceIdentifier) (entity.Service, error)
}

// V1 -.
type V1 struct {
	l *slog.Logger
	v *validator.Validate

	registry RegistryUseCase
}

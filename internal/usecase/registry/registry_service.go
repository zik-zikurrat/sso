package registry

import (
	"context"
	"log/slog"
	"sso/internal/usecase/dto/registry"

	"github.com/google/uuid"
)

type RegistryRepo interface {
	CreateService(ctx context.Context, in registry.CreateService) (string, error)
	// UpdateService(ctx context.Context, in registry.UpdateService) (uuid.UUID, error)
	// DeleteService(ctx context.Context, in uuid.UUID) error
	ListServiceEndpoints(ctx context.Context) ([]registry.ServiceWithEndpoints, error)
	GetServiceEndpointsByServiceID(ctx context.Context, in uuid.UUID) (registry.ServiceWithEndpoints, error)
}

type UseCase struct {
	l *slog.Logger
	r RegistryRepo
}

func NewRegistryUseCase(l *slog.Logger, r RegistryRepo) *UseCase {
	return &UseCase{
		l: l,
		r: r,
	}
}

func (uc *UseCase) RegisterService(ctx context.Context, in registry.CreateService) (string, error) {
	serviceID, err := uc.r.CreateService(ctx, in)
	if err != nil {
		uc.l.Error("failed to create service", slog.String("name", in.Name))
		return "", err
	}

	return serviceID, nil
}

func (uc *UseCase) ListServiceEndpoints(ctx context.Context) ([]registry.ServiceWithEndpoints, error) {
	services, err := uc.r.ListServiceEndpoints(ctx)
	if err != nil {
		uc.l.Error("failed to list service", slog.String("error", err.Error()))
		return nil, err
	}
	return services, nil
}

func (uc *UseCase) GetServiceEndpointsByServiceID(ctx context.Context, in uuid.UUID) (registry.ServiceWithEndpoints, error) {
	service, err := uc.r.GetServiceEndpointsByServiceID(ctx, in)
	if err != nil {
		uc.l.Error("failed to get service", slog.String("error", err.Error()))
		return registry.ServiceWithEndpoints{}, err
	}
	return service, nil
}

package registry

import (
	"context"
	"log/slog"
	"sso/internal/entity"
	"sso/internal/usecase/dto/registry"

	"github.com/google/uuid"
)

type RegistryRepo interface {
	CreateService(ctx context.Context, in registry.CreateService) error
	UpdateService(ctx context.Context, in registry.UpdateService) (uuid.UUID, error)
	DeleteService(ctx context.Context, in uuid.UUID) error
	ListService(ctx context.Context) ([]entity.Service, error)
	GetServiceByID(ctx context.Context, in entity.ServiceIdentifier) (entity.Service, error)
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

func (uc *UseCase) RegisterService(ctx context.Context, in registry.CreateService) error {
	if err := uc.r.CreateService(ctx, in); err != nil {
		uc.l.Error("failed to create service", slog.String("name", in.Name))
		return err
	}
	return nil
}

func (uc *UseCase) ListService(ctx context.Context) ([]entity.Service, error) {
	services, err := uc.r.ListService(ctx)
	if err != nil {
		uc.l.Error("failed to list service", slog.String("error", err.Error()))
		return nil, err
	}
	return services, nil
}

func (uc *UseCase) GetServiceByID(ctx context.Context, in entity.ServiceIdentifier) (entity.Service, error) {
	service, err := uc.r.GetServiceByID(ctx, in)
	if err != nil {
		uc.l.Error("failed to get service", slog.String("error", err.Error()))
		return entity.Service{}, err
	}
	return service, nil
}

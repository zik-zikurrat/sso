package persistent

import (
	"context"
	"fmt"
	"sso/internal/entity"
	"sso/internal/usecase/dto/registry"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const _defaultEntityCap = 64

type RegistryRepo struct {
	pool *pgxpool.Pool
}

func NewRegistryRepo(pool *pgxpool.Pool) *RegistryRepo {
	return &RegistryRepo{
		pool: pool,
	}
}

func (r *RegistryRepo) CreateService(ctx context.Context, in registry.CreateService) (string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("create service begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var serviceID uuid.UUID
	err = tx.QueryRow(ctx, insertServiceQuery, in.Name).Scan(&serviceID)
	if err != nil {
		return "", fmt.Errorf("create service: %w", err)
	}

	for i := 0; i < len(in.Endpoints); i++ {
		currEndpoint := in.Endpoints[i]
		if _, err := tx.Exec(ctx, insertEndpointQuery, serviceID, currEndpoint.Method, currEndpoint.URL, currEndpoint.Secure); err != nil {
			return "", fmt.Errorf("create service endpoint insert: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("create service commit: %w", err)
	}
	return serviceID.String(), nil
}

func (r *RegistryRepo) ListService(ctx context.Context) ([]registry.ListService, error) {
	rows, err := r.pool.Query(ctx, selectServiceQuery)
	if err != nil {
		return nil, fmt.Errorf("list service: %w", err)
	}
	defer rows.Close()

	out := make([]registry.ListService, 0, _defaultEntityCap)
	for rows.Next() {
		var e registry.ListService
		if err := rows.Scan(&e.ID, &e.Name, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, fmt.Errorf("list service scan: %w", err)
		}
		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list service rows: %w", err)
	}
	return out, nil
}

func (r *RegistryRepo) GetServiceByID(ctx context.Context, in entity.ServiceIdentifier) (entity.Service, error) {
	var s entity.Service
	err := r.pool.QueryRow(ctx, selectServiceQuery, in.ID).Scan(&s.ID, &s.Name, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return entity.Service{}, fmt.Errorf("get service: %w", err)
	}

	// догрузить endpoints этого сервиса
	rows, err := r.pool.Query(ctx, selectEndpointsByServiceQuery, s.ID)
	if err != nil {
		return entity.Service{}, fmt.Errorf("get service endpoints: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var e entity.Endpoint
		if err := rows.Scan(&e.Method, &e.URL, &e.Secure); err != nil {
			return entity.Service{}, fmt.Errorf("scan endpoint: %w", err)
		}
		s.Endpoints = append(s.Endpoints, e)
	}
	if err := rows.Err(); err != nil {
		return entity.Service{}, fmt.Errorf("endpoints rows: %w", err)
	}

	return s, nil
}

// func (r *RegistryRepo) UpdateService(ctx context.Context, in registry.UpdateService) (uuid.UUID, error) {
// }
// func (r *RegistryRepo) DeleteService(ctx context.Context, in uuid.UUID) error {

// }

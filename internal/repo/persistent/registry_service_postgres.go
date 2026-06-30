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

func (r *RegistryRepo) ListServiceEndpoints(ctx context.Context) ([]registry.ServiceWithEndpoints, error) {
	rows, err := r.pool.Query(ctx, selectServiceWithEndpointsQuery)
	if err != nil {
		return nil, fmt.Errorf("list service endpoints: %w", err)
	}
	defer rows.Close()

	out := make([]registry.ServiceWithEndpoints, 0, _defaultEntityCap)
	indexByID := make(map[uuid.UUID]int, _defaultEntityCap)
	for rows.Next() {
		var s registry.ServiceWithEndpoints
		var e registry.Endpoint
		if err := rows.Scan(&s.ID, &s.Name, &s.CreatedAt, &s.UpdatedAt, &e.ID, &e.Method, &e.URL, &e.Secure, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("list service endpoints scan: %w", err)
		}

		idx, ok := indexByID[s.ID]
		if !ok {
			s.Endpoints = make([]registry.Endpoint, 0, 1)
			out = append(out, s)
			idx = len(out) - 1
			indexByID[s.ID] = idx
		}
		out[idx].Endpoints = append(out[idx].Endpoints, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list service endpoints rows: %w", err)
	}
	return out, nil
}

func (r *RegistryRepo) GetServiceEndpointsByServiceID(ctx context.Context, in uuid.UUID) (entity.Service, error) {
	var s entity.Service
	err := r.pool.QueryRow(ctx, selectServiceWithEndpointsQuery, in.ID).Scan(&s.ID, &s.Name, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return entity.Service{}, fmt.Errorf("get service: %w", err)
	}

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

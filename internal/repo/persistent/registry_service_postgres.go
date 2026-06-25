package persistent

import (
	"context"
	"fmt"
	"sso/internal/entity"
	"sso/internal/usecase/dto/registry"

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

func (r *RegistryRepo) CreateService(ctx context.Context, in registry.CreateService) error {
	serviceName := in.Name
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("create service begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	for method, url := range in.Metadata {
		if _, err := tx.Exec(ctx, insertServicerQuery, serviceName, method, url); err != nil {
			return fmt.Errorf("create service insert: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("create service commit: %w", err)
	}
	return nil
}

func (r *RegistryRepo) ListService(ctx context.Context) ([]entity.Service, error) {
	rows, err := r.pool.Query(ctx, selectServiceQuery)
	if err != nil {
		return nil, fmt.Errorf("list service: %w", err)
	}
	defer rows.Close()

	out := make([]entity.Service, 0, _defaultEntityCap)
	for rows.Next() {
		var e entity.Service
		if err := rows.Scan(&e.ID, &e.Name, &e.Method, &e.URL, &e.Secure, &e.CreatedAt, &e.UpdatedAt); err != nil {
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
	err := r.pool.QueryRow(ctx, selectServiceQuery, in.ID).Scan(&s.ID, &s.Name, &s.Method, &s.URL, &s.Secure, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return entity.Service{}, nil
	}
	return s, nil
}

// func (r *RegistryRepo) UpdateService(ctx context.Context, in registry.UpdateService) (uuid.UUID, error) {
// }
// func (r *RegistryRepo) DeleteService(ctx context.Context, in uuid.UUID) error {

// }

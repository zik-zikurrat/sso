package persistent

import (
	"context"
	"fmt"
	identitycontext "sso/internal/entity/identity_context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepo struct {
	pool *pgxpool.Pool
}

func NewSessionRepo(pool *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{pool: pool}
}

func (r *SessionRepo) CreateSession(ctx context.Context, in identitycontext.UserSession) error {
	if _, err := r.pool.Exec(ctx, insertSessionQuery, in.UserID, in.TokenHash, in.ExpiresAt); err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

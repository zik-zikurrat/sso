package persistent

import (
	"context"
	"errors"
	"fmt"
	identitycontext "sso/internal/entity/identity_context"
	"sso/internal/usecase/dto/auth"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) CreateUser(ctx context.Context, in *auth.CreateUserRepo) (uuid.UUID, error) {
	var userID uuid.UUID
	err := r.pool.QueryRow(ctx, insertUserQuery, in.Email, in.PasswordHash, in.Login, in.Role).Scan(&userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create user: %w", err)
	}
	return userID, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, in auth.UpdateUser) (uuid.UUID, error) {
	var userID uuid.UUID
	err := r.pool.QueryRow(ctx, updateUserQuery, in.ID, in.Login, in.Email, in.Password).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, auth.ErrUserNotFound
		}
		return uuid.Nil, fmt.Errorf("update user: %w", err)
	}
	return userID, nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, deleteUserQuery, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return auth.ErrUserNotFound
	}
	return nil
}

func (r *UserRepo) GetByIdentifier(ctx context.Context, in identitycontext.UserIdentifier) (*identitycontext.User, error) {
	var u identitycontext.User
	err := r.pool.QueryRow(ctx, selectUserQuery, in.ID, in.Login, in.Email).
		Scan(&u.ID, &u.Login, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, auth.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &u, nil
}

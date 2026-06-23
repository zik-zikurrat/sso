package auth

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type UserRepo interface {
	CreateUser(ctx context.Context) (uuid.UUID, error)
}

type UseCase struct {
	repo UserRepo
	log  *slog.Logger
}

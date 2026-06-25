package auth

import (
	"context"
	"fmt"
	"log/slog"
	identitycontext "sso/internal/entity/identity_context"
	"sso/internal/usecase/dto/auth"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, in *auth.CreateUserRepo) (uuid.UUID, error)
	UpdateUser(ctx context.Context, in auth.UpdateUser) (uuid.UUID, error)
	DeleteUser(ctx context.Context, in uuid.UUID) error
	GetByIdentifier(ctx context.Context, in identitycontext.UserIdentifier) (*identitycontext.User, error)
}

type UseCase struct {
	l *slog.Logger
	r UserRepo
}

func NewUserUseCase(l *slog.Logger, r UserRepo) *UseCase {
	return &UseCase{
		l: l,
		r: r,
	}
}

func (uc *UseCase) Register(ctx context.Context, in auth.CreateUser) (uuid.UUID, error) {
	const op = "auth.Register"

	log := uc.l.With(
		slog.String("op", op),
		slog.String("email", in.Email),
	)

	log.Info("registering user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.String("error", err.Error()))
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	repoIN := auth.CreateUserRepo{
		Email:        in.Email,
		PasswordHash: passwordHash,
		Login:        in.Login,
		Role:         in.Role,
	}

	userID, err := uc.r.CreateUser(ctx, &repoIN)
	if err != nil {
		uc.l.Error("failed to save user", slog.String("errir", err.Error()))
		return uuid.Nil, err
	}
	return userID, nil
}

// func sendOtpToEmail(ctx context.Context, email string) error {
// 	otp, err := helper.RandomString(5)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (uc *UseCase) Login(ctx context.Context, email, password string, appID int32) (string, error) {

// }

// func (uc *UseCase) IsAdmin(ctx context.Context, userID string) (bool, error) {}
// func (uc *UseCase) IsDemo(ctx context.Context, userID string) (bool, error)  {}

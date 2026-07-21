package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"sso/internal/config"
	identitycontext "sso/internal/entity/identity_context"
	"sso/internal/helper"
	"sso/internal/usecase/dto/auth"
	"sso/pkg/cache"
	"sso/pkg/smtp"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const _tokenLength = 32
const _otpEpiryTime = 5

type UserRepo interface {
	CreateUser(ctx context.Context, in *auth.CreateUserRepo) (uuid.UUID, error)
	UpdateUser(ctx context.Context, in auth.UpdateUser) (uuid.UUID, error)
	DeleteUser(ctx context.Context, in uuid.UUID) error
	GetByIdentifier(ctx context.Context, in identitycontext.UserIdentifier) (*identitycontext.User, error)
}

type SessionRepo interface {
	CreateSession(ctx context.Context, in identitycontext.UserSession) error
}

type UseCase struct {
	l        *slog.Logger
	r        UserRepo
	s        SessionRepo
	c        *cache.Redis
	tokenTTL time.Duration
	smtpCfg  config.SMTPConfig
}

func NewUserUseCase(
	l *slog.Logger,
	r UserRepo,
	s SessionRepo,
	c *cache.Redis,
	tokenTTL time.Duration,
	smtpCfg config.SMTPConfig,
) *UseCase {
	return &UseCase{
		l:        l,
		r:        r,
		s:        s,
		c:        c,
		tokenTTL: tokenTTL,
		smtpCfg:  smtpCfg,
	}
}

func (uc *UseCase) Register(ctx context.Context, in auth.CreateUser) error {
	const op = "auth.Register"

	log := uc.l.With(
		slog.String("op", op),
		slog.String("email", in.Email),
	)

	log.Info("registering user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	repoIN := auth.CreateUserRepo{
		Email:        in.Email,
		PasswordHash: passwordHash,
		Login:        in.Login,
		Role:         in.Role,
	}

	otp, err := helper.RandomString(5)
	if err != nil {
		log.Error("failed to generate otp code", slog.String("error", err.Error()))
		return err
	}

	uc.c.Client.Set(ctx, otp, repoIN, time.Minute*_otpEpiryTime)
	err = smtp.SendEmailNotification(
		"Confirmation code to zik_zikurrat services",
		fmt.Sprintf("Your OTP code: %s", otp),
		[]string{in.Email},
		uc.smtpCfg,
		log,
	)
	if err != nil {
		log.Error("failed to send confirmation email", slog.String("error", err.Error()))
	}

	return nil
}

func (uc *UseCase) ConfirmOTP(ctx context.Context, otp string) (string, error) {
	var repoIN auth.CreateUserRepo
	userInfo := uc.c.Client.Get(ctx, otp)
	if err := json.Unmarshal(userInfo.Bytes(), &repoIN); err != nil {

	}
	userID, err := uc.r.CreateUser(ctx, &repoIN)
	if err != nil {
		log.Error("failed to save user", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

}

func (uc *UseCase) Login(ctx context.Context, email, password, appID string) (string, error) {
	const op = "auth.Login"

	log := uc.l.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.String("app_id", appID),
	)

	log.Info("logging in user")

	user, err := uc.r.GetByIdentifier(ctx, identitycontext.UserIdentifier{Email: &email})
	if err != nil {
		log.Warn("user not found", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, auth.ErrInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Warn("invalid password")
		return "", fmt.Errorf("%s: %w", op, auth.ErrInvalidCredentials)
	}

	token, err := helper.RandomString(_tokenLength)
	if err != nil {
		log.Error("failed to generate token", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	session := identitycontext.UserSession{
		UserID:    user.ID,
		TokenHash: hashToken(token),
		ExpiresAt: time.Now().Add(uc.tokenTTL),
	}
	if err := uc.s.CreateSession(ctx, session); err != nil {
		log.Error("failed to create session", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (uc *UseCase) IsAdmin(ctx context.Context, userID string) (bool, error) {
	return uc.hasRole(ctx, "auth.IsAdmin", userID, "admin")
}

func (uc *UseCase) IsDemo(ctx context.Context, userID string) (bool, error) {
	return uc.hasRole(ctx, "auth.IsDemo", userID, "demo")
}

func (uc *UseCase) hasRole(ctx context.Context, op, userID, role string) (bool, error) {
	log := uc.l.With(slog.String("op", op), slog.String("user_id", userID))

	id, err := uuid.Parse(userID)
	if err != nil {
		log.Warn("invalid user id", slog.String("error", err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	user, err := uc.r.GetByIdentifier(ctx, identitycontext.UserIdentifier{ID: &id})
	if err != nil {
		log.Error("failed to get user", slog.String("error", err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return user.Role == role, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

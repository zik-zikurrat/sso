package grpc

import (
	"context"
	"errors"
	"sso/internal/usecase/dto/auth"

	"buf.build/gen/go/zik-zikurrat-sso/sso/connectrpc/go/sso/v1/ssov1connect"
	ssov1 "buf.build/gen/go/zik-zikurrat-sso/sso/protocolbuffers/go/sso/v1"
	"connectrpc.com/connect"
)

type AuthUseCase interface {
	Register(ctx context.Context, in auth.CreateUser) (string, error)
	Login(ctx context.Context, email, password, appID string) (string, error)
	IsAdmin(ctx context.Context, userID string) (bool, error)
	IsDemo(ctx context.Context, userID string) (bool, error)
}

type AuthController struct {
	uc AuthUseCase
}

func NewAuthController(uc AuthUseCase) *AuthController {
	return &AuthController{uc: uc}
}

var _ ssov1connect.AuthServiceHandler = (*AuthController)(nil)

func (c *AuthController) Register(ctx context.Context, req *connect.Request[ssov1.RegisterRequest]) (*connect.Response[ssov1.RegisterResponse], error) {
	if req.Msg.GetEmail() == "" || req.Msg.GetPassword() == "" || req.Msg.GetLogin() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("login, email and password are required"))
	}

	id, err := c.uc.Register(ctx, auth.CreateUser{
		Login:    req.Msg.GetLogin(),
		Email:    req.Msg.GetEmail(),
		Password: req.Msg.GetPassword(),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&ssov1.RegisterResponse{UserId: id}), nil
}

func (c *AuthController) Login(ctx context.Context, req *connect.Request[ssov1.LoginRequest]) (*connect.Response[ssov1.LoginResponse], error) {
	if req.Msg.GetEmail() == "" || req.Msg.GetPassword() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("email and password are required"))
	}

	token, err := c.uc.Login(ctx, req.Msg.GetEmail(), req.Msg.GetPassword(), req.Msg.GetAppId())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&ssov1.LoginResponse{Token: token}), nil
}

func (c *AuthController) IsAdmin(ctx context.Context, req *connect.Request[ssov1.IsAdminRequest]) (*connect.Response[ssov1.IsAdminResponse], error) {
	if req.Msg.GetUserId() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("user_id is required"))
	}

	isAdmin, err := c.uc.IsAdmin(ctx, req.Msg.GetUserId())
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&ssov1.IsAdminResponse{IsAdmin: isAdmin}), nil
}

func (c *AuthController) IsDemo(ctx context.Context, req *connect.Request[ssov1.IsDemoRequest]) (*connect.Response[ssov1.IsDemoResponse], error) {
	if req.Msg.GetUserId() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("user_id is required"))
	}

	isDemo, err := c.uc.IsDemo(ctx, req.Msg.GetUserId())
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&ssov1.IsDemoResponse{IsDemo: isDemo}), nil
}

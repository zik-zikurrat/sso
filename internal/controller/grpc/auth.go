package grpc

import (
	"context"
	"sso/internal/usecase/dto/auth"

	"buf.build/gen/go/zik-zikurrat-sso/sso/connectrpc/go/sso/v1/ssov1connect"
	ssov1 "buf.build/gen/go/zik-zikurrat-sso/sso/protocolbuffers/go/sso/v1"
	"connectrpc.com/connect"
)

type AuthUseCase interface {
	Register(ctx context.Context, in auth.CreateUser) (string, error)
	Login(ctx context.Context, email, password string, appID int32) (string, error)
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

func (c *AuthController) IsAdmin(ctx context.Context, req *connect.Request[ssov1.IsAdminRequest]) (*connect.Response[ssov1.IsAdminResponse], error) {
	panic("NOT IMPLEMENTED")
}

func (c *AuthController) IsDemo(ctx context.Context, req *connect.Request[ssov1.IsDemoRequest]) (*connect.Response[ssov1.IsDemoResponse], error) {
	panic("NOT IMPLEMENTED")
}

func (c *AuthController) Login(ctx context.Context, req *connect.Request[ssov1.LoginRequest]) (*connect.Response[ssov1.LoginResponse], error) {
	panic("NOT IMPLEMENTED")
}

// Login, IsAdmin, IsDemo аналогично

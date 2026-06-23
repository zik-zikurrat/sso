package grpc

import (
	"context"
	"log/slog"

	ssov1 "buf.build/gen/go/zik-zikurrat-sso/sso/protocolbuffers/go/proto/sso/v1"
	"connectrpc.com/connect"
)

type AuthController struct {
	auth AuthUseCase
	log  *slog.Logger
}

func NewAuthController(auth AuthUseCase, log *slog.Logger) *AuthController {
	return &AuthController{auth: auth, log: log}
}

func (c *AuthController) Register(
	ctx context.Context,
	req *connect.Request[ssov1.RegisterRequest],
) (*connect.Response[ssov1.RegisterResponse], error) {
	// 1. достать данные из proto-сообщения
	// 2. позвать usecase (вся логика там)
	userID, err := c.auth.Register(ctx, req.Msg.GetLogin(), req.Msg.GetEmail(), req.Msg.GetPassword())
	if err != nil {
		// перевести доменную ошибку в connect.Error
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	// 3. обернуть результат в proto-ответ
	return connect.NewResponse(&ssov1.RegisterResponse{UserId: userID}), nil
}

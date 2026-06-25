package grpc

import (
	"context"
	"errors"
	"log/slog"
	"sso/internal/usecase/dto/registry"

	"buf.build/gen/go/zik-zikurrat-sso/sso/connectrpc/go/sso/v1/ssov1connect"
	ssov1 "buf.build/gen/go/zik-zikurrat-sso/sso/protocolbuffers/go/sso/v1"
	"connectrpc.com/connect"
)

type RegistryUseCase interface {
	RegisterService(ctx context.Context, in registry.CreateService) error
}

type RegistryController struct {
	l  *slog.Logger
	uc RegistryUseCase
}

func NewRegistryController(l *slog.Logger, uc RegistryUseCase) *RegistryController {
	return &RegistryController{uc: uc, l: l}
}

var _ ssov1connect.RegistryServiceHandler = (*RegistryController)(nil)

func (c *RegistryController) RegisterService(ctx context.Context, req *connect.Request[ssov1.RegisterServiceRequest]) (*connect.Response[ssov1.RegisterServiceResponse], error) {
	if err := validateRegisterService(req); err != nil {
		return nil, err
	}
	if err := c.uc.RegisterService(ctx, registry.CreateService{
		Name:     req.Msg.GetName(),
		Metadata: req.Msg.GetMetadata(),
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&ssov1.RegisterServiceResponse{Msg: "OK"}), nil
}

func validateRegisterService(req *connect.Request[ssov1.RegisterServiceRequest]) error {
	msg := req.Msg
	if msg.GetName() == "" {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("name of service is required"))
	}
	if msg.GetMetadata() == nil {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("metadata of service is required"))
	}
	return nil
}

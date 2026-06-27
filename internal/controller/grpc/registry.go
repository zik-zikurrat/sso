package grpc

import (
	"context"
	"errors"
	"log/slog"

	"buf.build/gen/go/zik-zikurrat-sso/sso/connectrpc/go/sso/v1/ssov1connect"
	ssov1 "buf.build/gen/go/zik-zikurrat-sso/sso/protocolbuffers/go/sso/v1"
	"connectrpc.com/connect"

	"sso/internal/usecase/dto/registry"
)

type RegistryUseCase interface {
	RegisterService(ctx context.Context, in registry.CreateService) (string, error)
}

type RegistryController struct {
	l  *slog.Logger
	uc RegistryUseCase
}

func NewRegistryController(l *slog.Logger, uc RegistryUseCase) *RegistryController {
	return &RegistryController{uc: uc, l: l}
}

var _ ssov1connect.RegistryServiceHandler = (*RegistryController)(nil)

func (c *RegistryController) RegisterService(
	ctx context.Context,
	req *connect.Request[ssov1.RegisterServiceRequest],
) (*connect.Response[ssov1.RegisterServiceResponse], error) {
	if err := validateRegisterService(req); err != nil {
		return nil, err
	}

	endpoints := make([]registry.Endpoint, 0, len(req.Msg.GetEndpoints()))
	for _, e := range req.Msg.GetEndpoints() {
		endpoints = append(endpoints, registry.Endpoint{
			Method: e.GetMethod(),
			URL:    e.GetUrl(),
			Secure: e.GetSecure(),
		})
	}

	serviceID, err := c.uc.RegisterService(ctx, registry.CreateService{
		Name:      req.Msg.GetName(),
		Endpoints: endpoints,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&ssov1.RegisterServiceResponse{
		ServiceId: serviceID,
	}), nil
}

func validateRegisterService(req *connect.Request[ssov1.RegisterServiceRequest]) error {
	msg := req.Msg
	if msg.GetName() == "" {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("name of service is required"))
	}
	if len(msg.GetEndpoints()) == 0 {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("endpoints are required"))
	}
	return nil
}

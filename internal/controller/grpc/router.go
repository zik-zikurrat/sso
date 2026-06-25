package grpc

import (
	"net/http"

	"buf.build/gen/go/zik-zikurrat-sso/sso/connectrpc/go/sso/v1/ssov1connect"
)

func NewMux(auth *AuthController, registry *RegistryController) http.Handler {
	mux := http.NewServeMux()

	authPath, authH := ssov1connect.NewAuthServiceHandler(auth)
	mux.Handle(authPath, authH)

	regPath, regH := ssov1connect.NewRegistryServiceHandler(registry)
	mux.Handle(regPath, regH)

	return mux
}

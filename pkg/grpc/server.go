package grpc

import (
	"context"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func New(mux http.Handler, addr string) *Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)
	srv.Protocols = protocols

	return &Server{httpServer: srv}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

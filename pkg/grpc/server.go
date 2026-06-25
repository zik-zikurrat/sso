package grpc

import (
	"context"
	"net"
	"net/http"
	"sso/internal/config"
	"strconv"
)

type Server struct {
	httpServer *http.Server
	notify     chan error
}

func New(mux http.Handler, cfg *config.Config) *Server {
	addr := net.JoinHostPort(cfg.Server.HOST, strconv.Itoa(cfg.GRPC.Port))
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)
	srv.Protocols = protocols

	return &Server{
		httpServer: srv,
		notify:     make(chan error, 1),
	}
}

func (s *Server) Start() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.notify <- err
		}
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

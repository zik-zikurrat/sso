package grpc

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"sso/internal/config"
	"strconv"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	notify     chan error
}

func New(mux http.Handler, cfg *config.Config, l *slog.Logger) *Server {
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
		logger:     l,
		notify:     make(chan error, 1),
	}
}

func (s *Server) Start() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.notify <- err
		}
	}()
	s.logger.Info("gRPC server - Server - Started")
}

func (s *Server) Notify() <-chan error {
	return s.notify
}
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("gRPC server - Server - Shutdown")
	return s.httpServer.Shutdown(ctx)
}

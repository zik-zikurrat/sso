package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sso/internal/config"
	"sso/internal/helper"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
)

// Server -.
type Server struct {
	cfg *config.Config
	ctx context.Context
	eg  *errgroup.Group

	App    *fiber.App
	notify chan error

	address string
	prefork bool
	logger  *slog.Logger
}

// New -.
func New(ctx context.Context, l *slog.Logger, cfg *config.Config) *Server {
	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(1)

	s := &Server{
		cfg:     cfg,
		ctx:     ctx,
		eg:      group,
		App:     nil,
		notify:  make(chan error, 1),
		prefork: false,
		address: helper.GetServerAddr(cfg),
		logger:  l,
	}

	app := fiber.New(fiber.Config{
		Prefork:      cfg.Server.Prefork,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
	})

	s.App = app

	return s

}

// Start -.
func (s *Server) Start() {
	s.eg.Go(func() error {
		err := s.App.Listen(s.address)
		if err != nil {
			s.notify <- err

			return err
		}
		return nil
	})
	s.logger.Info("restapi server - Server - Started")
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	var shutdownError []error

	err := s.App.ShutdownWithTimeout(time.Duration(s.cfg.Server.ShutdownTimeout) * time.Second)
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error("restapi server - Server - Shutdown - s.App.ShutdownWithTimeout", slog.String("error", err.Error()))

		shutdownError = append(shutdownError, err)

	}

	err = s.eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error("restapi server - Server - Shutdown - s.eg.Wait", slog.String("error", err.Error()))

		shutdownError = append(shutdownError, err)

	}

	s.logger.Info("restapi server - Server - Shutdown")

	return errors.Join(shutdownError...)

}

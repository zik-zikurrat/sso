package app

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/config"
	grpcCtrl "sso/internal/controller/grpc"
	"sso/internal/controller/restapi"
	"sso/internal/repo/persistent"
	"sso/internal/usecase/registry"
	grpcSrv "sso/pkg/grpc"
	"sso/pkg/httpserver"
	sloghandler "sso/pkg/logger"
	"sso/pkg/postgres"
	"syscall"
	"time"
)

func Run(cfg *config.Config) error {
	Migrate(cfg)
	log := SetupLogger(cfg.Logging)
	log.Info("starting application", slog.Any("config", cfg))
	pg, err := postgres.New(cfg, log)
	if err != nil {
		log.Error("app - Run - postgres.New", slog.String("error", err.Error()))
		return err
	}
	defer pg.Close()

	// Repo
	// userRepo := postgres.NewUserRepo(pg.Pool)
	registryRepo := persistent.NewRegistryRepo(pg.Pool)
	// UseCase
	// authUC := auth.NewUseCase(userRepo, log)
	registryUC := registry.NewRegistryUseCase(log, registryRepo)
	// Controller
	// authCtrl := grpccontroller.NewAuthController(authUC, log)
	registryCtrl := grpcCtrl.NewRegistryController(log, registryUC)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	grpcServer := grpcSrv.New(grpcCtrl.NewMux(registryCtrl), cfg)
	go grpcServer.Start()

	httpserver := httpserver.New(ctx, log, cfg)
	restapi.NewRouter(
		httpserver.App,
		cfg,
		pg.Pool,
		log,
		registryUC,
	)
	httpserver.Start()

	select {
	case <-ctx.Done():
		log.Info("app - Run - shutdown signal received")
	case err := <-httpserver.Notify():
		log.Error("app - Run - httpserver.Notify", slog.String("error", err.Error()))
	case err := <-grpcServer.Notify():
		log.Error("app - Run - grpcserver.Notify", slog.String("error", err.Error()))
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout*time.Second)
	defer cancel()
	shutdownErr := []error{
		httpserver.Shutdown(),
		grpcServer.Shutdown(shutdownCtx),
	}
	return errors.Join(shutdownErr...)
}

func SetupLogger(cfg config.LoggingConfig) *slog.Logger {
	var level slog.Level

	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &sloghandler.Options{
		Level: level,
	}

	var output io.Writer
	if cfg.Discard {
		output = io.Discard
	} else {
		output = os.Stdout
	}

	var handler slog.Handler
	switch cfg.Format {
	case "json":
		handler = sloghandler.NewPrettyJSONHandler(output, opts)
	default:
		handler = sloghandler.NewPrettyJSONHandler(output, opts)
	}

	return slog.New(handler)
}

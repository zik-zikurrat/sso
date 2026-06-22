package app

import (
	"io"
	"log/slog"
	"os"
	"sso/internal/config"
	sloghandler "sso/pkg/logger"
)

func Run(cfg *config.Config) error {
	Migrate(cfg)
	log := SetupLogger(cfg.Logging)
	log.Info("starting application", slog.Any("config", cfg))
	return nil
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

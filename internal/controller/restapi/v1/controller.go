package v1

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	l *slog.Logger
	v *validator.Validate
}

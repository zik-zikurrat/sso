package v1

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewSSORoutes(
	apiV1Group fiber.Router,
	l *slog.Logger,
) {
	r := &V1{
		l: l,
		v: validator.New(validator.WithRequiredStructEnabled()),
	}
	trainingGroup := apiV1Group.Group("/training")

	{
		// RegistryService
		trainingGroup.Post("/structure", r.CreateStructure)
		trainingGroup.Get("/structure", r.ListStructure)
		trainingGroup.Get("/structure/:id", r.GetStructure)
		trainingGroup.Patch("/structure/:id", r.UpdateStructure)
		trainingGroup.Delete("/structure/:id", r.DeleteStructure)
	}
}

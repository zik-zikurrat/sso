package v1

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewSSORoutes(
	apiV1Group fiber.Router,
	l *slog.Logger,

	registry RegistryUseCase,
) {
	r := &V1{
		l:        l,
		v:        validator.New(validator.WithRequiredStructEnabled()),
		registry: registry,
	}
	registryGroup := apiV1Group.Group("/registry")

	{
		// RegistryService
		registryGroup.Get("/service", r.ListServiceEndpoints)
		registryGroup.Get("/service/:id", r.GetServiceEndpointsByServiceID)
		// registryGroup.Patch("/service/:id", r.UpdateService)
		// registryGroup.Delete("/service/:id", r.DeleteService)
	}
}

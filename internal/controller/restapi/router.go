package restapi

import (
	"log/slog"
	"sso/internal/config"
	"sso/internal/controller/restapi/middleware"
	v1 "sso/internal/controller/restapi/v1"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewRouter -.
// Swagger spec:
// @title       SSO
// @description Single Sign On
// @version     1.0
// @host        localhost:3033
// @BasePath    /v1
func NewRouter(
	app *fiber.App,
	cfg *config.Config,
	pool *pgxpool.Pool,
	l *slog.Logger,
	registry v1.RegistryUseCase,
) {

	app.Use(recover.New())
	// Prometheus metrics TODO
	// Swagger TODO
	// app.Get("/swagger/*", swagger.HandlerDefault)

	apiV1Group := app.Group("/v1")
	apiV1Group.Get("/health", func(c *fiber.Ctx) error {
		if err := pool.Ping(c.Context()); err != nil {
			return c.SendStatus(503)
		}
		return c.SendStatus(200)
	})
	{
		// Tracing
		apiV1Group.Use(middleware.TracingMiddleware())
		// Options
		apiV1Group.Use(middleware.LoggerMiddleware(l))
		// Cors
		apiV1Group.Use(cors.New())
		v1.NewSSORoutes(
			apiV1Group,
			l,
			registry,
		)
	}
}

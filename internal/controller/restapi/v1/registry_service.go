package v1

import (
	"log/slog"
	"net/http"
	"sso/internal/controller/restapi/v1/response"
	"sso/internal/entity"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary     List service
// @Description List service
// @ID          listService
// @Tags  	    listService
// @Accept      json
// @Produce     json
// @Success     200
// @Failure     500 {object} response.Error
// @Router      /registry/service [get]
func (r *V1) ListService(ctx *fiber.Ctx) error {
	services, err := r.registry.ListService(ctx.UserContext())
	if err != nil {
		r.l.Error("restapi - v1 - service", slog.String("error", err.Error()))
		return errorResponse(ctx, http.StatusInternalServerError, "error while getting services")
	}
	return ctx.Status(http.StatusOK).JSON(response.ToServiceList(services))
}

// @Summary     Get service
// @Description Get service
// @ID          getService
// @Tags  	    getService
// @Accept      json
// @Produce     json
// @Success     200
// @Failure     500 {object} response.Error
// @Router      /registry/service [get]
func (r *V1) GetServiceByID(ctx *fiber.Ctx) error {
	uuidID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid service id")
	}

	service, err := r.registry.GetServiceByID(ctx.UserContext(), entity.ServiceIdentifier{ID: &uuidID})
	if err != nil {
		r.l.Error("restapi - v1 - service", slog.String("error", err.Error()))
		return errorResponse(ctx, http.StatusInternalServerError, "error while getting service")
	}

	return ctx.Status(http.StatusOK).JSON(response.ToService(service))
}

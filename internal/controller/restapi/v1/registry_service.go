package v1

import (
	"errors"
	"log/slog"
	"net/http"
	"sso/internal/controller/restapi/v1/request"
	"sso/internal/controller/restapi/v1/response"
	"sso/internal/usecase/dto/registry"

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
func (r *V1) ListServiceEndpoints(ctx *fiber.Ctx) error {
	services, err := r.registry.ListServiceEndpoints(ctx.UserContext())
	if err != nil {
		r.l.Error("restapi - v1 - service", slog.String("error", err.Error()))
		return errorResponse(ctx, http.StatusInternalServerError, "error while getting services")
	}
	return ctx.Status(http.StatusOK).JSON(response.ToServiceEndpoints(services))
}

// @Summary     Get service
// @Description Get service with its endpoints by service id
// @ID          getServiceEndpointsByServiceID
// @Tags  	    getService
// @Accept      json
// @Produce     json
// @Param       id  path string true "Service ID"
// @Success     200
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /registry/service/{id} [get]
func (r *V1) GetServiceEndpointsByServiceID(ctx *fiber.Ctx) error {
	serviceID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid service id")
	}

	service, err := r.registry.GetServiceEndpointsByServiceID(ctx.UserContext(), serviceID)
	if err != nil {
		r.l.Error("restapi - v1 - service", slog.String("error", err.Error()))
		return errorResponse(ctx, http.StatusInternalServerError, "error while getting service")
	}

	return ctx.Status(http.StatusOK).JSON(response.ToServiceEndpoint(service))
}

// @Summary     Delete service
// @Description Delete service and all its endpoints by service id
// @ID          deleteService
// @Tags  	    deleteService
// @Accept      json
// @Produce     json
// @Param       id  path string true "Service ID"
// @Success     204
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /registry/service/{id} [delete]
func (r *V1) DeleteService(ctx *fiber.Ctx) error {
	serviceID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid service id")
	}

	if err := r.registry.DeleteService(ctx.UserContext(), serviceID); err != nil {
		if errors.Is(err, registry.ErrServiceNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "service not found")
		}
		r.l.Error("restapi - v1 - service", slog.String("error", err.Error()))
		return errorResponse(ctx, http.StatusInternalServerError, "error while deleting service")
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Update service
// @Description Update service
// @ID          updateService
// @Tags  	    updateService
// @Accept      json
// @Produce     json
// @Param       id  path string true "Service ID"
// @Success     200
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /registry/service/{id} [patch]
func (r *V1) UpdateService(ctx *fiber.Ctx) error {
	var req request.UpdateService
	serviceID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid service id")
	}
	if err := ctx.BodyParser(&req); err != nil {
		r.l.Error("restapi - v1 - service", slog.String("error", err.Error()))
		return errorResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}

	if err := r.registry.UpdateService(ctx.UserContext(), registry.UpdateService{
		ID:   serviceID,
		Name: req.Name,
	}); err != nil {
		if errors.Is(err, registry.ErrServiceNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "service not found")
		}
		r.l.Error("restapi - v1 - service", slog.String("error", err.Error()))
		return errorResponse(ctx, http.StatusInternalServerError, "error while updating service")
	}

	return ctx.SendStatus(http.StatusOK)
}

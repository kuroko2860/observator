package handler

import (
	"github.com/labstack/echo/v4"
	"kuroko.com/analystics/internal/model"
)

// @Summary		Get All Operations From Service
// @Description	Get All Operations From Service
// @Tags			service
// @Accept			json
// @Produce		json
// @Param			service_name	query		string	true	"Service Name"
// @Success		200				{object}	[]string
// @Failure		500				{object}	model.Error
// @Router			/services/:service_name/operations [get]
func (h *Handler) GetAllOperationsFromServiceHandler(c echo.Context) error {
	serviceName := c.Param("service_name")

	res, err := h.service.GetAllOperationsFromService(c.Request().Context(), serviceName)
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, res)
}

// @Summary		Get All Operations Count From Service
// @Description	Get All Operations Count From Service
// @Tags			service
// @Accept			json
// @Produce		json
// @Param			service_name	query		string	true	"Service Name"
// @Param			from			query		string	true	"From"
// @Param			to				query		string	true	"To"
// @Success		200				{object}	[]string
// @Failure		500				{object}	model.Error
// @Router			/operations-count [get]
// func (h *Handler) GetAllOperationsCountFromServiceHandler(c echo.Context) error {
// 	serviceName := c.QueryParam("service_name")
// 	from := c.QueryParam("from")
// 	to := c.QueryParam("to")

// 	res, err := h.service.GetAllOperationsCountFromService(c.Request().Context(), serviceName, from, to)
// 	if err != nil {
// 		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
// 	}

// 	return c.JSON(200, res)
// }

// @Summary		Get All Services
// @Description	Get All Services
// @Tags			service
// @Accept			json
// @Produce		json
// @Success		200				{object}	[]string
// @Failure		500				{object}	model.Error
// @Router			/services [get]
func (h *Handler) GetAllServicesHandler(c echo.Context) error {
	res, err := h.service.GetAllServices(c.Request().Context())
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, res)
}

// @Summary		Get Service Detail
// @Description	Get Service Detail
// @Tags			service
// @Accept			json
// @Produce		json
// @Param			service_name	query		string	true	"Service Name"
// @Param			from			query		string	true	"From"
// @Param			to				query		string	true	"To"
// @Success		200				{object}	model.ServiceDetail
// @Failure		500				{object}	model.Error
// @Router			/services/:service_name [get]
func (h *Handler) GetServiceDetailHandler(c echo.Context) error {
	serviceName := c.Param("service_name")
	from := c.QueryParam("from")
	to := c.QueryParam("to")

	res, err := h.service.GetServiceDetailService(c.Request().Context(), serviceName, from, to)
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, res)
}

// @Summary		Get Http Service Api
// @Description	Get Http Service Api
// @Tags			service
// @Accept			json
// @Produce		json
// @Param			service_name	query		string	true	"Service Name"
// @Param			from			query		string	true	"From"
// @Param			to				query		string	true	"To"
// @Success		200				{object}	any
// @Failure		500				{object}	model.Error
// @Router			/http-service-api [get]
// func (h *Handler) GetHttpServiceApiHandler(c echo.Context) error {
// 	serviceName := c.QueryParam("service_name")
// 	from := c.QueryParam("from")
// 	to := c.QueryParam("to")

// 	res, err := h.service.GetHttpServiceApiService(c.Request().Context(), serviceName, from, to)
// 	if err != nil {
// 		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
// 	}

// 	return c.JSON(200, res)
// }

// @Summary		Get Service Endpoint
// @Description	Get Service Endpoint
// @Tags			service
// @Accept			json
// @Produce		json
// @Param			service_name	query		string	true	"Service Name"
// @Success		200				{object}	[]string
// @Failure		500				{object}	model.Error
// @Router			/services/:service_name/endpoints [get]
func (h *Handler) GetServiceEndpointHandler(c echo.Context) error {
	serviceName := c.Param("service_name")

	res, err := h.service.GetServiceEndpointService(c.Request().Context(), serviceName)
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, res)
}

// @Summary		Get Top Called Service
// @Description	Get Top Called Service
// @Tags			service
// @Accept			json
// @Produce		json
// @Param			from			query		string	true	"From"
// @Param			to				query		string	true	"To"
// @Param			limit			query		string	true	"Limit"
// @Success		200				{object}	map[string]int
// @Failure		500				{object}	model.Error
// @Router			/services/top-called [get]
func (h *Handler) GetTopCalledServiceHandler(c echo.Context) error {
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	limit := c.QueryParam("limit")

	res, err := h.service.GetTopCalledService(c.Request().Context(), from, to, limit)
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}

	return c.JSON(200, res)

}

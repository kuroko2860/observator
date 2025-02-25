package handler

import (
	"github.com/labstack/echo/v4"
	"kuroko.com/analystics/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) RegisterRoutes(v1 *echo.Group) {
	v1.GET("/traces", h.getAllTracesOfPath)

	v1.GET("/api-statistics", h.GetApiStatisticHandler)
	v1.GET("/long-api", h.GetLongApiHandler)
	v1.GET("/called-api", h.GetCalledApiHandler)
	v1.GET("/top-called-api", h.GetTopCalledApiHandler)
	v1.GET("/http-api", h.GetHttpApiByServiceHandler)

	v1.GET("/paths", h.GetAllPathFromHopHandler)       // handle db data
	v1.GET("/path-detail", h.GetPathDetailByIdHandler) // handle db data
	v1.GET("/hop-detail", h.GetHopDetailByIdHandler)
	v1.GET("/long-path", h.GetLongPathHandler)

	v1.GET("/operations", h.GetAllOperationsFromServiceHandler)
	v1.GET("/operations-count", h.GetAllOperationsCountFromServiceHandler)
	v1.GET("/services", h.GetAllServicesHandler)
	v1.GET("/services/:service_name", h.GetServiceDetailHandler)
	v1.GET("/http-service-api", h.GetHttpServiceApiHandler)
	v1.GET("/service-endpoint", h.GetServiceEndpointHandler)
	v1.GET("/top-called-service", h.GetTopCalledServiceHandler)

	v1.GET("/get-alert", h.GetAlertHandler)
	v1.GET("/uri-list", h.GetUriListHandler)
	v1.PATCH("/ignore-alert/:id", h.IgnoreAlertHandler)
	v1.GET("/online-time", h.OnlineTimeHandler)
	v1.GET("/online-user", h.OnlineUserHandler)
	v1.GET("/service-statistic", h.ServiceStatisticHandler)
	v1.GET("/uri-statistic", h.UriStatisticHandler)
	v1.GET("/usage", h.UsageHandler)

}

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
	// user view specific path then click view traces and view specific trace
	v1.GET("/paths/:path_id/traces", h.getAllTracesOfPath)
	v1.GET("/traces/:trace_id", h.getTraceById)

	v1.POST("/paths", h.GetAllPathFromOperationsHandler)
	v1.GET("/paths/:path_id", h.GetPathDetailByIdHandler)
	v1.GET("/paths/long", h.GetLongPathHandler)
	v1.GET("/hops/:hop_id", h.GetHopDetailByIdHandler)

	v1.GET("/services/:service_name/operations", h.GetAllOperationsFromServiceHandler)
	v1.GET("/services", h.GetAllServicesHandler)
	v1.GET("/services/top-called", h.GetTopCalledServiceHandler)
	v1.GET("/services/:service_name", h.GetServiceDetailHandler)
	v1.GET("/services/:service_name/endpoints", h.GetServiceEndpointHandler)
	// v1.GET("/http-service-api", h.GetHttpServiceApiHandler)
	// v1.GET("/operations-count", h.GetAllOperationsCountFromServiceHandler)

	v1.GET("/api-statistics", h.GetApiStatisticHandler)
	v1.GET("/api-statistics/long", h.GetLongApiHandler)
	v1.GET("/api-statistics/user-called", h.GetCalledApiHandler)
	v1.GET("/api-statistics/top-called", h.GetTopCalledApiHandler)
	v1.GET("/get-alert", h.GetAlertHandler)
	v1.GET("/uri-list", h.GetUriListHandler)
	v1.PATCH("/ignore-alert/:id", h.IgnoreAlertHandler)
	// v1.GET("/online-time", h.OnlineTimeHandler)
	// v1.GET("/online-user", h.OnlineUserHandler)
	v1.GET("/service-statistic", h.ServiceStatisticHandler)
	v1.GET("/uri-statistic", h.UriStatisticHandler)
	v1.GET("/usage", h.UsageHandler)
	// v1.GET("/http-api", h.GetHttpApiByServiceHandler)

	// Log query routes
	// Elasticsearch logs
	v1.GET("/logs/elasticsearch/trace/:trace_id", h.GetElasticsearchLogsByTraceId)
	v1.GET("/logs/elasticsearch/span/:span_id", h.GetElasticsearchLogsBySpanId)
	v1.GET("/logs/elasticsearch/trace/:trace_id/span/:span_id", h.GetElasticsearchLogsByTraceAndSpanId)
	v1.GET("/logs/elasticsearch/service/:service_name", h.GetElasticsearchLogsByService)
	v1.GET("/logs/elasticsearch/time-range", h.GetElasticsearchLogsByTimeRange)
	
	// MongoDB logs
	v1.GET("/logs/mongodb/trace/:trace_id", h.GetMongoDBLogsByTraceId)
	v1.GET("/logs/mongodb/span/:span_id", h.GetMongoDBLogsBySpanId)
	v1.GET("/logs/mongodb/trace/:trace_id/span/:span_id", h.GetMongoDBLogsByTraceAndSpanId)
}

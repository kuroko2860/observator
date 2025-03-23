package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"kuroko.com/analystics/internal/model"
)

// @Summary      Get logs by trace ID from Elasticsearch
// @Description  Retrieve all logs with the given trace ID from Elasticsearch
// @Tags         logs
// @Accept       json
// @Produce      json
// @Param        trace_id  path     string  true  "Trace ID"
// @Param        from      query    int     false "Starting offset for results"
// @Param        size      query    int     false "Number of results to return"
// @Success      200       {array}  map[string]interface{}
// @Failure      500       {object} model.Error
// @Router       /logs/elasticsearch/trace/{trace_id} [get]
func (h *Handler) GetElasticsearchLogsByTraceId(c echo.Context) error {
	traceId := c.Param("trace_id")
	from, size := getPaginationParams(c)
	
	logs, err := h.service.QueryLogsByTraceID(c.Request().Context(), traceId, from, size)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}
	return c.JSON(http.StatusOK, logs)
}

// @Summary      Get logs by span ID from Elasticsearch
// @Description  Retrieve all logs with the given span ID from Elasticsearch
// @Tags         logs
// @Accept       json
// @Produce      json
// @Param        span_id   path     string  true  "Span ID"
// @Param        from      query    int     false "Starting offset for results"
// @Param        size      query    int     false "Number of results to return"
// @Success      200       {array}  map[string]interface{}
// @Failure      500       {object} model.Error
// @Router       /logs/elasticsearch/span/{span_id} [get]
func (h *Handler) GetElasticsearchLogsBySpanId(c echo.Context) error {
	spanId := c.Param("span_id")
	from, size := getPaginationParams(c)
	
	logs, err := h.service.QueryLogsBySpanID(c.Request().Context(), spanId, from, size)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}
	return c.JSON(http.StatusOK, logs)
}

// @Summary      Get logs by trace ID and span ID from Elasticsearch
// @Description  Retrieve all logs with both the given trace ID and span ID from Elasticsearch
// @Tags         logs
// @Accept       json
// @Produce      json
// @Param        trace_id  path     string  true  "Trace ID"
// @Param        span_id   path     string  true  "Span ID"
// @Param        from      query    int     false "Starting offset for results"
// @Param        size      query    int     false "Number of results to return"
// @Success      200       {array}  map[string]interface{}
// @Failure      500       {object} model.Error
// @Router       /logs/elasticsearch/trace/{trace_id}/span/{span_id} [get]
func (h *Handler) GetElasticsearchLogsByTraceAndSpanId(c echo.Context) error {
	traceId := c.Param("trace_id")
	spanId := c.Param("span_id")
	from, size := getPaginationParams(c)
	
	logs, err := h.service.QueryLogsByTraceAndSpanID(c.Request().Context(), traceId, spanId, from, size)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}
	return c.JSON(http.StatusOK, logs)
}

// @Summary      Get logs by service name from Elasticsearch
// @Description  Retrieve all logs from the specified service from Elasticsearch
// @Tags         logs
// @Accept       json
// @Produce      json
// @Param        service_name  path     string  true  "Service Name"
// @Param        from          query    int     false "Starting offset for results"
// @Param        size          query    int     false "Number of results to return"
// @Success      200           {array}  map[string]interface{}
// @Failure      500           {object} model.Error
// @Router       /logs/elasticsearch/service/{service_name} [get]
func (h *Handler) GetElasticsearchLogsByService(c echo.Context) error {
	serviceName := c.Param("service_name")
	from, size := getPaginationParams(c)
	
	logs, err := h.service.QueryLogsByService(c.Request().Context(), serviceName, from, size)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}
	return c.JSON(http.StatusOK, logs)
}

// @Summary      Get logs by time range from Elasticsearch
// @Description  Retrieve all logs within the specified time range from Elasticsearch
// @Tags         logs
// @Accept       json
// @Produce      json
// @Param        start_time  query    int64  true  "Start time (Unix timestamp in milliseconds)"
// @Param        end_time    query    int64  true  "End time (Unix timestamp in milliseconds)"
// @Param        from        query    int    false "Starting offset for results"
// @Param        size        query    int    false "Number of results to return"
// @Success      200         {array}  map[string]interface{}
// @Failure      500         {object} model.Error
// @Router       /logs/elasticsearch/time-range [get]
func (h *Handler) GetElasticsearchLogsByTimeRange(c echo.Context) error {
	startTimeStr := c.QueryParam("start_time")
	endTimeStr := c.QueryParam("end_time")
	
	startTime, err := strconv.ParseInt(startTimeStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{
			Message: "Invalid start_time parameter",
			Code:    http.StatusBadRequest,
		})
	}
	
	endTime, err := strconv.ParseInt(endTimeStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Error{
			Message: "Invalid end_time parameter",
			Code:    http.StatusBadRequest,
		})
	}
	
	from, size := getPaginationParams(c)
	
	logs, err := h.service.QueryLogsByTimeRange(c.Request().Context(), startTime, endTime, from, size)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}
	return c.JSON(http.StatusOK, logs)
}

// @Summary      Get logs by trace ID from MongoDB
// @Description  Retrieve all logs with the given trace ID from MongoDB
// @Tags         logs
// @Accept       json
// @Produce      json
// @Param        trace_id  path     string  true  "Trace ID"
// @Success      200       {array}  model.HttpLogEntry
// @Failure      500       {object} model.Error
// @Router       /logs/mongodb/trace/{trace_id} [get]
func (h *Handler) GetMongoDBLogsByTraceId(c echo.Context) error {
	traceId := c.Param("trace_id")
	logs, err := h.service.FindHttpLogEntriesByTraceId(c.Request().Context(), traceId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}
	return c.JSON(http.StatusOK, logs)
}

// @Summary      Get logs by span ID from MongoDB
// @Description  Retrieve all logs with the given span ID from MongoDB
// @Tags         logs
// @Accept       json
// @Produce      json
// @Param        span_id   path     string  true  "Span ID"
// @Success      200       {array}  model.HttpLogEntry
// @Failure      500       {object} model.Error
// @Router       /logs/mongodb/span/{span_id} [get]
func (h *Handler) GetMongoDBLogsBySpanId(c echo.Context) error {
	spanId := c.Param("span_id")
	logs, err := h.service.FindHttpLogEntriesBySpanId(c.Request().Context(), spanId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}
	return c.JSON(http.StatusOK, logs)
}

// @Summary      Get logs by trace ID and span ID from MongoDB
// @Description  Retrieve all logs with both the given trace ID and span ID from MongoDB
// @Tags         logs
// @Accept       json
// @Produce      json
// @Param        trace_id  path     string  true  "Trace ID"
// @Param        span_id   path     string  true  "Span ID"
// @Success      200       {array}  model.HttpLogEntry
// @Failure      500       {object} model.Error
// @Router       /logs/mongodb/trace/{trace_id}/span/{span_id} [get]
func (h *Handler) GetMongoDBLogsByTraceAndSpanId(c echo.Context) error {
	traceId := c.Param("trace_id")
	spanId := c.Param("span_id")
	logs, err := h.service.FindHttpLogEntriesByTraceAndSpanId(c.Request().Context(), traceId, spanId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Error{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}
	return c.JSON(http.StatusOK, logs)
}

// Helper function to get pagination parameters from the request
func getPaginationParams(c echo.Context) (int, int) {
	fromStr := c.QueryParam("from")
	sizeStr := c.QueryParam("size")
	
	from := 0
	if fromVal, err := strconv.Atoi(fromStr); err == nil {
		from = fromVal
	}
	
	size := 10
	if sizeVal, err := strconv.Atoi(sizeStr); err == nil && sizeVal > 0 {
		size = sizeVal
	}
	
	return from, size
}
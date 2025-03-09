package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"kuroko.com/analystics/internal/model"
)

// @Summary		Get all traces of path
// @Description	Get all traces of path
// @Tags			traces
// @Accept			json
// @Produce		json
// @Param			path_id			path		string	true	"Path Id"
// @Param			from		query		string	true	"from"
// @Param			to		query		string	true	"to"
// @Success		200		{object}	[][]model.Span
// @Failure		400		{object}	model.Error
// @Router			/paths/:path_id/traces [get]
func (h *Handler) getAllTracesOfPath(c echo.Context) error {
	pathIdStr := c.Param("path_id")
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	pathId, _ := strconv.Atoi(pathIdStr)
	res, err := h.service.GetAllTracesOfPath(c.Request().Context(), uint32(pathId), from, to)
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}
	return c.JSON(200, res)
}

func (h *Handler) getTraceById(c echo.Context) error {
	traceId := c.Param("trace_id")
	res, err := h.service.GetTraceById(c.Request().Context(), traceId)
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}
	return c.JSON(200, res)
}

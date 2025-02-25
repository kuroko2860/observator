package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"kuroko.com/analystics/internal/model"
)

//		@Summary		Get all traces of path
//		@Description	Get all traces of path
//		@Tags			traces
//		@Accept			json
//		@Produce		json
//	   @Param			path_id			query		uint32	true	"Path Id"
//		@Success		200		{object}	[][]model.Span
//		@Failure		400		{object}	model.Error
//		@Router			/traces [get]
func (h *Handler) getAllTracesOfPath(c echo.Context) error {
	pathIdStr := c.QueryParam("path_id")
	pathId, _ := strconv.Atoi(pathIdStr)
	res, err := h.service.GetAllTracesOfPath(c.Request().Context(), uint32(pathId))
	if err != nil {
		return c.JSON(500, model.Error{Message: err.Error(), Code: 500})
	}
	return c.JSON(200, res)
}

package handler

import "github.com/labstack/echo/v4"

func (h *Handler) getAllTracesOfPath(c echo.Context) error {
	pathId := c.QueryParam("pathId")
	spans := h.service.GetAllTracesOfPath(c.Request().Context(), pathId)
	return c.JSON(200, spans)
}

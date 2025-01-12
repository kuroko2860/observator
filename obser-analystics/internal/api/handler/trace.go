package handler

import "github.com/labstack/echo/v4"

func (h *Handler) getAllSpans(c echo.Context) error {
	spans := h.service.GetAllSpans(c.Request().Context())
	return c.JSON(200, spans)
}

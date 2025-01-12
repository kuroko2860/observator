package handler

import "github.com/labstack/echo/v4"

func (h *Handler) getPathById(c echo.Context) error {
	id := c.Param("id")
	path, err := h.service.GetPathById(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(200, path)
}

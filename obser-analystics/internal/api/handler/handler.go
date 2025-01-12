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
	v1.GET("/paths/:id", h.getPathById)

}

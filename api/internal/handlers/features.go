package handlers

import (
	"api/internal/models"
	"net/http"

	"github.com/labstack/echo"
)

// GetFeatures ...
func (h *Handler) GetFeatures(ctx echo.Context) error {
	features := models.GetFeatures(h.Config)
	return ctx.JSON(http.StatusOK, features)
}

package handlers

import (
	"api/internal/models"
	"net/http"

	"github.com/labstack/echo"
)

// GetVersion ...
func (h *Handler) GetVersion(ctx echo.Context) error {
	version := models.GetVersion()
	return ctx.JSON(http.StatusOK, version)
}

// GetFeatures ...
func (h *Handler) GetFeatures(ctx echo.Context) error {
	features := models.GetFeatures(h.Config)
	return ctx.JSON(http.StatusOK, features)
}

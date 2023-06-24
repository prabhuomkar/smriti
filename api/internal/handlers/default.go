package handlers

import (
	"api/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetVersion ...
func (h *Handler) GetVersion(ctx echo.Context) error {
	version := models.GetVersion()
	return ctx.JSON(http.StatusOK, version)
}

// GetFeatures ...
func (h *Handler) GetFeatures(ctx echo.Context) error {
	cfgFeatures := models.GetFeatures(h.Config)
	features, _ := ctx.Get("features").(models.Features)

	features.Favourites = features.Favourites && cfgFeatures.Favourites
	features.Hidden = features.Hidden && cfgFeatures.Hidden
	features.Trash = features.Trash && cfgFeatures.Trash
	features.Albums = features.Albums && cfgFeatures.Albums
	features.Explore = features.Explore && cfgFeatures.Explore
	features.Places = features.Places && cfgFeatures.Places
	features.Things = features.Things && cfgFeatures.Things
	features.People = features.People && cfgFeatures.People
	features.Sharing = features.Sharing && cfgFeatures.Sharing

	return ctx.JSON(http.StatusOK, features)
}

// GetDisk ...
func (h *Handler) GetDisk(ctx echo.Context) error {
	disk := models.GetDisk(h.Config)
	return ctx.JSON(http.StatusOK, disk)
}

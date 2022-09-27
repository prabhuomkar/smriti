package handlers

import (
	"api/internal/models"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// GetFavouriteMediaItems ...
func (h *Handler) GetFavouriteMediaItems(ctx echo.Context) error {
	offset, limit := getOffsetAndLimit(ctx)
	favourites := []models.MediaItem{}
	result := h.DB.Where("is_favourite=true").Find(&favourites).Offset(offset).Limit(limit)
	if result.Error != nil {
		log.Printf("error getting favourite mediaitems: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, favourites)
}

// GetHiddenMediaItems ...
func (h *Handler) GetHiddenMediaItems(ctx echo.Context) error {
	offset, limit := getOffsetAndLimit(ctx)
	hidden := []models.MediaItem{}
	result := h.DB.Where("is_hidden=true").Find(&hidden).Offset(offset).Limit(limit)
	if result.Error != nil {
		log.Printf("error getting hidden mediaitems: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, hidden)
}

// GetDeletedMediaItems ...
func (h *Handler) GetDeletedMediaItems(ctx echo.Context) error {
	offset, limit := getOffsetAndLimit(ctx)
	deleted := []models.MediaItem{}
	result := h.DB.Where("is_deleted=true").Find(&deleted).Offset(offset).Limit(limit)
	if result.Error != nil {
		log.Printf("error getting deleted mediaitems: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, deleted)
}

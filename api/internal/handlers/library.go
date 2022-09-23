package handlers

import (
	"api/internal/models"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// GetFavouriteMediaItems ...
func (h *Handler) GetFavouriteMediaItems(ctx echo.Context) error {
	favourites := []models.MediaItem{}
	err := h.DB.Select(&favourites, "SELECT * FROM mediaitems WHERE is_favourite=true")
	if err != nil {
		log.Printf("error getting favourites: %+v", err)
		return echo.ErrInternalServerError
	}
	return ctx.JSON(http.StatusOK, favourites)
}

// GetHiddenMediaItems ...
func (h *Handler) GetHiddenMediaItems(ctx echo.Context) error {
	hidden := []models.MediaItem{}
	err := h.DB.Select(&hidden, "SELECT * FROM mediaitems WHERE is_hidden=true")
	if err != nil {
		log.Printf("error getting hidden: %+v", err)
		return echo.ErrInternalServerError
	}
	return ctx.JSON(http.StatusOK, hidden)
}

// GetDeletedMediaItems ...
func (h *Handler) GetDeletedMediaItems(ctx echo.Context) error {
	deleted := []models.MediaItem{}
	err := h.DB.Select(&deleted, "SELECT * FROM mediaitems WHERE is_deleted=true")
	if err != nil {
		log.Printf("error getting deleted: %+v", err)
		return echo.ErrInternalServerError
	}
	return ctx.JSON(http.StatusOK, deleted)
}

package handlers

import (
	"api/internal/models"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// GetMediaItemPlaces ...
func (h *Handler) GetMediaItemPlaces(ctx echo.Context) error {
	return nil
}

// GetMediaItemThings ...
func (h *Handler) GetMediaItemThings(ctx echo.Context) error {
	return nil
}

// GetMediaItemPeople ...
func (h *Handler) GetMediaItemPeople(ctx echo.Context) error {
	return nil
}

// GetMediaItem ...
func (h *Handler) GetMediaItem(ctx echo.Context) error {
	return nil
}

// UpdateMediaItem ...
func (h *Handler) UpdateMediaItem(ctx echo.Context) error {
	return nil
}

// DeleteMediaItem ...
func (h *Handler) DeleteMediaItem(ctx echo.Context) error {
	return nil
}

// GetMediaItems ...
func (h *Handler) GetMediaItems(ctx echo.Context) error {
	mediaItems := []models.MediaItem{}
	err := h.DB.Select(&mediaItems, "SELECT * FROM mediaitems WHERE (is_hidden=false OR is_deleted=false)")
	if err != nil {
		log.Printf("error getting mediaitems: %+v", err)
		return echo.ErrInternalServerError
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// UploadMediaItems ...
func (h *Handler) UploadMediaItems(ctx echo.Context) error {
	return nil
}

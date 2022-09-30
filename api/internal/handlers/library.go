package handlers

import (
	"api/internal/models"
	"log"
	"net/http"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// GetFavouriteMediaItems ...
func (h *Handler) GetFavouriteMediaItems(ctx echo.Context) error {
	offset, limit := getOffsetAndLimit(ctx)
	favourites := []models.MediaItem{}
	result := h.DB.Where("is_favourite=true AND is_deleted=false").Find(&favourites).Offset(offset).Limit(limit)
	if result.Error != nil {
		log.Printf("error getting favourite mediaitems: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, favourites)
}

// AddFavouriteMediaItems ...
func (h *Handler) AddFavouriteMediaItems(ctx echo.Context) error {
	mediaItems, err := getMediaItems(ctx)
	if err != nil {
		return err
	}
	mediaItemIDs := make([]uuid.UUID, len(mediaItems))
	for idx, mediaItem := range mediaItems {
		mediaItemIDs[idx] = mediaItem.ID
	}
	result := h.DB.Model(&models.MediaItem{}).Where("id IN ?", mediaItemIDs).
		Updates(map[string]interface{}{"is_favourite": true})
	if result.Error != nil {
		log.Printf("error adding favourite mediaitems: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// RemoveFavouriteMediaItems ...
func (h *Handler) RemoveFavouriteMediaItems(ctx echo.Context) error {
	mediaItems, err := getMediaItems(ctx)
	if err != nil {
		return err
	}
	mediaItemIDs := make([]uuid.UUID, len(mediaItems))
	for idx, mediaItem := range mediaItems {
		mediaItemIDs[idx] = mediaItem.ID
	}
	result := h.DB.Model(&models.MediaItem{}).Where("id IN ?", mediaItemIDs).
		Updates(map[string]interface{}{"is_favourite": false})
	if result.Error != nil {
		log.Printf("error removing favourite mediaitems: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// GetHiddenMediaItems ...
func (h *Handler) GetHiddenMediaItems(ctx echo.Context) error {
	offset, limit := getOffsetAndLimit(ctx)
	hidden := []models.MediaItem{}
	result := h.DB.Where("is_hidden=true AND is_deleted=false").Find(&hidden).Offset(offset).Limit(limit)
	if result.Error != nil {
		log.Printf("error getting hidden mediaitems: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, hidden)
}

// AddHiddenMediaItems ...
func (h *Handler) AddHiddenMediaItems(ctx echo.Context) error {
	mediaItems, err := getMediaItems(ctx)
	if err != nil {
		return err
	}
	mediaItemIDs := make([]uuid.UUID, len(mediaItems))
	for idx, mediaItem := range mediaItems {
		mediaItemIDs[idx] = mediaItem.ID
	}
	result := h.DB.Model(&models.MediaItem{}).Where("id IN ?", mediaItemIDs).
		Updates(map[string]interface{}{"is_hidden": true})
	if result.Error != nil {
		log.Printf("error adding hidden mediaitems: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// RemoveHiddenMediaItems ...
func (h *Handler) RemoveHiddenMediaItems(ctx echo.Context) error {
	mediaItems, err := getMediaItems(ctx)
	if err != nil {
		return err
	}
	mediaItemIDs := make([]uuid.UUID, len(mediaItems))
	for idx, mediaItem := range mediaItems {
		mediaItemIDs[idx] = mediaItem.ID
	}
	result := h.DB.Model(&models.MediaItem{}).Where("id IN ?", mediaItemIDs).
		Updates(map[string]interface{}{"is_hidden": false})
	if result.Error != nil {
		log.Printf("error removing hidden mediaitems: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
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

// AddDeletedMediaItems ...
func (h *Handler) AddDeletedMediaItems(ctx echo.Context) error {
	mediaItems, err := getMediaItems(ctx)
	if err != nil {
		return err
	}
	mediaItemIDs := make([]uuid.UUID, len(mediaItems))
	for idx, mediaItem := range mediaItems {
		mediaItemIDs[idx] = mediaItem.ID
	}
	result := h.DB.Model(&models.MediaItem{}).Where("id IN ?", mediaItemIDs).
		Updates(map[string]interface{}{"is_deleted": true})
	if result.Error != nil {
		log.Printf("error adding deleted mediaitems: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// RemoveDeletedMediaItems ...
func (h *Handler) RemoveDeletedMediaItems(ctx echo.Context) error {
	mediaItems, err := getMediaItems(ctx)
	if err != nil {
		return err
	}
	mediaItemIDs := make([]uuid.UUID, len(mediaItems))
	for idx, mediaItem := range mediaItems {
		mediaItemIDs[idx] = mediaItem.ID
	}
	result := h.DB.Model(&models.MediaItem{}).Where("id IN ?", mediaItemIDs).
		Updates(map[string]interface{}{"is_deleted": false})
	if result.Error != nil {
		log.Printf("error removing deleted mediaitems: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

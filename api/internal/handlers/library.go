package handlers

import (
	"api/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/exp/slog"
)

// GetFavouriteMediaItems ...
func (h *Handler) GetFavouriteMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	favourites := []models.MediaItem{}
	result := h.DB.Where("user_id=? AND is_favourite=true AND is_deleted=false", userID).
		Find(&favourites).Offset(offset).Limit(limit)
	if result.Error != nil {
		slog.Error("error getting favourite mediaitems", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, favourites)
}

// AddFavouriteMediaItems ...
func (h *Handler) AddFavouriteMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	mediaItems, err := getMediaItems(ctx)
	if err != nil {
		return err
	}
	mediaItemIDs := make([]uuid.UUID, len(mediaItems))
	for idx, mediaItem := range mediaItems {
		mediaItemIDs[idx] = mediaItem.ID
	}
	result := h.DB.Model(&models.MediaItem{}).Where("user_id=? AND id IN ?", userID, mediaItemIDs).
		Updates(map[string]interface{}{"is_favourite": true})
	if result.Error != nil {
		slog.Error("error adding favourite mediaitems", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// RemoveFavouriteMediaItems ...
func (h *Handler) RemoveFavouriteMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	mediaItems, err := getMediaItems(ctx)
	if err != nil {
		return err
	}
	mediaItemIDs := make([]uuid.UUID, len(mediaItems))
	for idx, mediaItem := range mediaItems {
		mediaItemIDs[idx] = mediaItem.ID
	}
	result := h.DB.Model(&models.MediaItem{}).Where("user_id=? AND id IN ?", userID, mediaItemIDs).
		Updates(map[string]interface{}{"is_favourite": false})
	if result.Error != nil {
		slog.Error("error removing favourite mediaitems", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// GetHiddenMediaItems ...
func (h *Handler) GetHiddenMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	hidden := []models.MediaItem{}
	result := h.DB.Where("user_id=? AND is_hidden=true AND is_deleted=false", userID).
		Find(&hidden).Offset(offset).Limit(limit)
	if result.Error != nil {
		slog.Error("error getting hidden mediaitems", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, hidden)
}

// AddHiddenMediaItems ...
func (h *Handler) AddHiddenMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	mediaItems, err := getMediaItems(ctx)
	if err != nil {
		return err
	}
	mediaItemIDs := make([]uuid.UUID, len(mediaItems))
	for idx, mediaItem := range mediaItems {
		mediaItemIDs[idx] = mediaItem.ID
	}
	result := h.DB.Model(&models.MediaItem{}).Where("user_id=? AND id IN ?", userID, mediaItemIDs).
		Updates(map[string]interface{}{"is_hidden": true})
	if result.Error != nil {
		slog.Error("error adding hidden mediaitems", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// RemoveHiddenMediaItems ...
func (h *Handler) RemoveHiddenMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	mediaItems, err := getMediaItems(ctx)
	if err != nil {
		return err
	}
	mediaItemIDs := make([]uuid.UUID, len(mediaItems))
	for idx, mediaItem := range mediaItems {
		mediaItemIDs[idx] = mediaItem.ID
	}
	result := h.DB.Model(&models.MediaItem{}).Where("user_id=? AND id IN ?", userID, mediaItemIDs).
		Updates(map[string]interface{}{"is_hidden": false})
	if result.Error != nil {
		slog.Error("error removing hidden mediaitems", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// GetDeletedMediaItems ...
func (h *Handler) GetDeletedMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	deleted := []models.MediaItem{}
	result := h.DB.Where("user_id=? AND is_deleted=true", userID).Find(&deleted).Offset(offset).Limit(limit)
	if result.Error != nil {
		slog.Error("error getting deleted mediaitems", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, deleted)
}

// AddDeletedMediaItems ...
func (h *Handler) AddDeletedMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	mediaItems, err := getMediaItems(ctx)
	if err != nil {
		return err
	}
	mediaItemIDs := make([]uuid.UUID, len(mediaItems))
	for idx, mediaItem := range mediaItems {
		mediaItemIDs[idx] = mediaItem.ID
	}
	result := h.DB.Model(&models.MediaItem{}).Where("user_id=? AND id IN ?", userID, mediaItemIDs).
		Updates(map[string]interface{}{"is_deleted": true})
	if result.Error != nil {
		slog.Error("error adding deleted mediaitems", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// RemoveDeletedMediaItems ...
func (h *Handler) RemoveDeletedMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	mediaItems, err := getMediaItems(ctx)
	if err != nil {
		return err
	}
	mediaItemIDs := make([]uuid.UUID, len(mediaItems))
	for idx, mediaItem := range mediaItems {
		mediaItemIDs[idx] = mediaItem.ID
	}
	result := h.DB.Model(&models.MediaItem{}).Where("user_id=? AND id IN ?", userID, mediaItemIDs).
		Updates(map[string]interface{}{"is_deleted": false})
	if result.Error != nil {
		slog.Error("error removing deleted mediaitems", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

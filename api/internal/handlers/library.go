package handlers

import (
	"api/internal/models"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

const (
	queryGetFavouriteMediaItems    = `SELECT * FROM mediaitems WHERE user_id=$1 AND is_favourite=true AND is_deleted=false OFFSET $2 LIMIT $3`
	queryGetHiddenMediaItems       = `SELECT * FROM mediaitems WHERE user_id=$1 AND is_hidden=true AND is_deleted=false OFFSET $2 LIMIT $3`
	queryGetDeletedMediaItems      = `SELECT * FROM mediaitems WHERE user_id=$1 AND is_deleted=true OFFSET $2 LIMIT $3`
	queryAddFavouriteMediaItems    = `UPDATE mediaitems SET is_favourite=true WHERE user_id=$1 AND id=ANY($2)`
	queryRemoveFavouriteMediaItems = `UPDATE mediaitems SET is_favourite=false WHERE user_id=$1 AND id=ANY($2)`
	queryAddHiddenMediaItems       = `UPDATE mediaitems SET is_hidden=true WHERE user_id=$1 AND id=ANY($2)`
	queryRemoveHiddenMediaItems    = `UPDATE mediaitems SET is_hidden=false WHERE user_id=$1 AND id=ANY($2)`
	queryAddDeletedMediaItems      = `UPDATE mediaitems SET is_deleted=true WHERE user_id=$1 AND id=ANY($2)`
	queryRemoveDeletedMediaItems   = `UPDATE mediaitems SET is_deleted=false WHERE user_id=$1 AND id=ANY($2)`
)

// GetFavouriteMediaItems ...
func (h *Handler) GetFavouriteMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	favourites := []models.MediaItem{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetFavouriteMediaItems, userID, offset, limit)
	if err != nil {
		slog.Error("error getting favourite mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		mediaItem, err := models.ScanRowsToMediaItem(rows)
		if err != nil {
			slog.Error("error scanning favourite mediaitem", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		favourites = append(favourites, mediaItem)
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
	_, err = h.DB.Exec(ctx.Request().Context(), queryAddFavouriteMediaItems, userID, mediaItemIDs)
	if err != nil {
		slog.Error("error adding favourite mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	_, err = h.DB.Exec(ctx.Request().Context(), queryRemoveFavouriteMediaItems, userID, mediaItemIDs)
	if err != nil {
		slog.Error("error removing favourite mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

// GetHiddenMediaItems ...
func (h *Handler) GetHiddenMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	hidden := []models.MediaItem{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetHiddenMediaItems, userID, offset, limit)
	if err != nil {
		slog.Error("error getting hidden mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		mediaItem, err := models.ScanRowsToMediaItem(rows)
		if err != nil {
			slog.Error("error scanning hidden mediaitem", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		hidden = append(hidden, mediaItem)
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
	_, err = h.DB.Exec(ctx.Request().Context(), queryAddHiddenMediaItems, userID, mediaItemIDs)
	if err != nil {
		slog.Error("error adding hidden mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	_, err = h.DB.Exec(ctx.Request().Context(), queryRemoveHiddenMediaItems, userID, mediaItemIDs)
	if err != nil {
		slog.Error("error removing hidden mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

// GetDeletedMediaItems ...
func (h *Handler) GetDeletedMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	deleted := []models.MediaItem{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetDeletedMediaItems, userID, offset, limit)
	if err != nil {
		slog.Error("error getting deleted mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		mediaItem, err := models.ScanRowsToMediaItem(rows)
		if err != nil {
			slog.Error("error scanning deleted mediaitem", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		deleted = append(deleted, mediaItem)
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
	_, err = h.DB.Exec(ctx.Request().Context(), queryAddDeletedMediaItems, userID, mediaItemIDs)
	if err != nil {
		slog.Error("error adding deleted mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	_, err = h.DB.Exec(ctx.Request().Context(), queryRemoveDeletedMediaItems, userID, mediaItemIDs)
	if err != nil {
		slog.Error("error removing deleted mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

package handlers

import (
	"api/internal/models"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

const (
	queryGetSharedAlbumMediaItems = `SELECT * FROM mediaitems WHERE id IN (SELECT mediaitem_id FROM album_mediaitems` +
		` WHERE shared = true AND album_id=$1) AND is_hidden = false AND is_deleted = false ORDER BY created_at DESC OFFSET $2 LIMIT $3`
	queryGetSharedAlbum = `SELECT a.*, m.* FROM albums a LEFT JOIN mediaitems m ON a.cover_mediaitem_id = m.id` +
		` WHERE a.shared = true AND a.id=$1 GROUP BY a.id, m.id`
)

// GetSharedAlbumMediaItems ...
func (h *Handler) GetSharedAlbumMediaItems(ctx echo.Context) error {
	offset, limit := getOffsetAndLimit(ctx)
	uid, err := getSharedAlbumID(ctx)
	if err != nil {
		return err
	}
	mediaItems := []models.MediaItem{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetSharedAlbumMediaItems, uid, offset, limit)
	if err != nil {
		slog.Error("error getting shared album mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		mediaItem, err := models.ScanRowsToMediaItem(rows)
		if err != nil {
			slog.Error("error scanning shared album mediaitem", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		mediaItems = append(mediaItems, mediaItem)
	}

	return ctx.JSON(http.StatusOK, mediaItems)
}

// GetSharedAlbum ...
func (h *Handler) GetSharedAlbum(ctx echo.Context) error {
	uid, err := getSharedAlbumID(ctx)
	if err != nil {
		return err
	}
	sharedAlbum := models.Album{CoverMediaItem: &models.MediaItem{}}
	err = h.DB.QueryRow(ctx.Request().Context(), queryGetSharedAlbum, uid).Scan(&sharedAlbum.ID,
		&sharedAlbum.UserID, &sharedAlbum.Name, &sharedAlbum.Description, &sharedAlbum.IsShared,
		&sharedAlbum.IsHidden, &sharedAlbum.MediaItemsCount, &sharedAlbum.CoverMediaItemID, &sharedAlbum.CreatedAt,
		&sharedAlbum.UpdatedAt, &sharedAlbum.CoverMediaItem.ID, &sharedAlbum.CoverMediaItem.UserID,
		&sharedAlbum.CoverMediaItem.Filename, &sharedAlbum.CoverMediaItem.Hash, &sharedAlbum.CoverMediaItem.Description,
		&sharedAlbum.CoverMediaItem.MimeType, &sharedAlbum.CoverMediaItem.SourceURL,
		&sharedAlbum.CoverMediaItem.PreviewURL, &sharedAlbum.CoverMediaItem.ThumbnailURL,
		&sharedAlbum.CoverMediaItem.Placeholder, &sharedAlbum.CoverMediaItem.IsFavourite,
		&sharedAlbum.CoverMediaItem.IsHidden, &sharedAlbum.CoverMediaItem.IsDeleted,
		&sharedAlbum.CoverMediaItem.Status, &sharedAlbum.CoverMediaItem.MediaItemType,
		&sharedAlbum.CoverMediaItem.MediaItemCategory, &sharedAlbum.CoverMediaItem.Width,
		&sharedAlbum.CoverMediaItem.Height, &sharedAlbum.CoverMediaItem.CreationTime,
		&sharedAlbum.CoverMediaItem.CameraMake, &sharedAlbum.CoverMediaItem.CameraModel,
		&sharedAlbum.CoverMediaItem.FocalLength, &sharedAlbum.CoverMediaItem.ApertureFnumber,
		&sharedAlbum.CoverMediaItem.IsoEquivalent, &sharedAlbum.CoverMediaItem.ExposureTime,
		&sharedAlbum.CoverMediaItem.Megapixels, &sharedAlbum.CoverMediaItem.Latitude,
		&sharedAlbum.CoverMediaItem.Longitude, &sharedAlbum.CoverMediaItem.FPS, &sharedAlbum.CoverMediaItem.EXIFData,
		&sharedAlbum.CoverMediaItem.Keywords, &sharedAlbum.CoverMediaItem.CreatedAt, &sharedAlbum.CoverMediaItem.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "shared link not found")
		}
		slog.Error("error getting shared album", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, sharedAlbum)
}

func getSharedAlbumID(ctx echo.Context) (uuid.UUID, error) {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting shared album id", "error", err)

		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid shared link")
	}

	return uid, err
}

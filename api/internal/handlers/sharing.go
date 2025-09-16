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
	queryGetSharedAlbum = `SELECT a.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,` +
		` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM albums a LEFT JOIN mediaitems m` +
		` ON a.cover_mediaitem_id=m.id WHERE a.shared = true AND a.id=$1`
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

	sharedAlbum := models.Album{}
	coverMediaItem := &models.CoverMediaItem{}

	err = h.DB.QueryRow(ctx.Request().Context(), queryGetSharedAlbum, uid).Scan(&sharedAlbum.ID,
		&sharedAlbum.UserID, &sharedAlbum.Name, &sharedAlbum.Description, &sharedAlbum.IsShared,
		&sharedAlbum.IsHidden, &sharedAlbum.MediaItemsCount, &sharedAlbum.CoverMediaItemID, &sharedAlbum.CreatedAt,
		&sharedAlbum.UpdatedAt, &coverMediaItem.ID, &coverMediaItem.UserID, &coverMediaItem.SourceURL,
		&coverMediaItem.PreviewURL, &coverMediaItem.ThumbnailURL, &coverMediaItem.Placeholder,
		&coverMediaItem.MediaItemType, &coverMediaItem.MediaItemCategory, &coverMediaItem.Width, &coverMediaItem.Height)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "shared link not found")
		}
		slog.Error("error getting shared album", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if sharedAlbum.CoverMediaItemID != nil {
		sharedAlbum.CoverMediaItem = coverMediaItem
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

package handlers

import (
	"api/internal/models"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type ( // AlbumRequest ...
	AlbumRequest struct {
		Name             *string `json:"name"`
		Description      *string `json:"description"`
		IsShared         *bool   `json:"shared"`
		IsHidden         *bool   `json:"hidden"`
		CoverMediaItemID *string `json:"coverMediaItemId"`
	}

	// MediaItemsRequest ...
	MediaItemsRequest struct {
		MediaItems []string `json:"mediaItems" required:"true"`
	}
)

const (
	queryGetAlbumMediaItems = `SELECT * FROM mediaitems WHERE id IN (SELECT mediaitem_id FROM album_mediaitems WHERE user_id=$1` +
		` AND album_id=$2) AND is_hidden=false AND is_deleted=false ORDER BY created_at DESC OFFSET $3 LIMIT $4`
	queryAddAlbumMediaItems          = `INSERT INTO album_mediaitems (album_id, mediaitem_id) VALUES ($1, $2)`
	queryRemoveAlbumMediaItems       = `DELETE FROM album_mediaitems WHERE album_id=$1 AND mediaitem_id=$2`
	queryGetAlbumMediaItemIDAndCount = `SELECT mediaitem_id, COUNT(*) OVER() AS mediaitems_count FROM album_mediaitems WHERE` +
		` album_id=$1 LIMIT 1`
	queryUpdateAlbumMediaItems = `UPDATE albums SET mediaitems_count=$1, cover_mediaitem_id=CASE WHEN cover_mediaitem_id` +
		` IS NULL THEN $2 ELSE cover_mediaitem_id END WHERE user_id=$3 AND id=$4`
	queryGetAlbum = `SELECT a.*, m.* FROM albums a LEFT JOIN mediaitems m ON a.cover_mediaitem_id=m.id` +
		` WHERE a.user_id=$1 AND a.id=$2 GROUP BY a.id, m.id`
	queryGetAlbums = `SELECT a.*, m.* FROM albums a LEFT JOIN mediaitems m ON a.cover_mediaitem_id=m.id` +
		` WHERE a.user_id=$1 AND a.is_hidden=false AND is_shared=$2 GROUP BY a.id, m.id ORDER BY %s` +
		` OFFSET $3 LIMIT $4`
	queryUpdateAlbum = `UPDATE albums SET name=$3, description=$4, is_shared=$5, is_hidden=$6,` +
		` cover_mediaitem_id=$7, updated_at=$8 WHERE user_id=$1 AND id=$2`
	queryDeleteAlbum = `DELETE FROM albums WHERE user_id=$1 AND id=$2`
	queryCreateAlbum = `INSERT INTO albums (id, user_id, name, description, is_shared, is_hidden,` +
		` created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
)

// GetAlbumMediaItems ...
func (h *Handler) GetAlbumMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	uid, err := getAlbumID(ctx)
	if err != nil {
		return err
	}
	mediaItems := []models.MediaItem{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetAlbumMediaItems, userID, uid, offset, limit)
	if err != nil {
		slog.Error("error getting album mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		mediaItem, err := models.ScanRowsToMediaItem(rows)
		if err != nil {
			slog.Error("error scanning album mediaitem", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		mediaItems = append(mediaItems, mediaItem)
	}

	return ctx.JSON(http.StatusOK, mediaItems)
}

// AddAlbumMediaItems ...
func (h *Handler) AddAlbumMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	uid, err := getAlbumID(ctx)
	if err != nil {
		return err
	}
	mediaItems, err := getMediaItems(ctx)
	if err != nil {
		return err
	}
	atx, err := h.DB.Begin(ctx.Request().Context())
	if err != nil {
		slog.Error("error starting transaction for adding album mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer func() {
		if err != nil {
			_ = atx.Rollback(ctx.Request().Context())
		}
	}()
	for _, mediaItem := range mediaItems {
		_, err = atx.Exec(ctx.Request().Context(), queryAddAlbumMediaItems, uid, mediaItem.ID)
		if err != nil {
			slog.Error("error adding album mediaitem", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	coverMediaItemID := uuid.Nil
	mediaItemsCount := 0
	err = atx.QueryRow(ctx.Request().Context(), queryGetAlbumMediaItemIDAndCount, uid).
		Scan(&coverMediaItemID, &mediaItemsCount)
	if err != nil && err != pgx.ErrNoRows {
		slog.Error("error getting album mediaitem id and count", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	_, err = atx.Exec(ctx.Request().Context(), queryUpdateAlbumMediaItems, mediaItemsCount, coverMediaItemID, userID, uid)
	if err != nil {
		slog.Error("error updating album", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err = atx.Commit(ctx.Request().Context()); err != nil {
		slog.Error("error committing transaction for adding album mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

// RemoveAlbumMediaItems ...
func (h *Handler) RemoveAlbumMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	uid, err := getAlbumID(ctx)
	if err != nil {
		return err
	}
	mediaItems, err := getMediaItems(ctx)
	if err != nil {
		return err
	}
	atx, err := h.DB.Begin(ctx.Request().Context())
	if err != nil {
		slog.Error("error starting transaction for removing album mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer func() {
		if err != nil {
			_ = atx.Rollback(ctx.Request().Context())
		}
	}()
	for _, mediaItem := range mediaItems {
		_, err = atx.Exec(ctx.Request().Context(), queryRemoveAlbumMediaItems, uid, mediaItem.ID)
		if err != nil {
			slog.Error("error removing album mediaitem", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	coverMediaItemID := uuid.Nil
	mediaItemsCount := 0
	err = atx.QueryRow(ctx.Request().Context(), queryGetAlbumMediaItemIDAndCount, uid).
		Scan(&coverMediaItemID, &mediaItemsCount)
	if err != nil && err != pgx.ErrNoRows {
		slog.Error("error getting album mediaitem id and count", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	_, err = atx.Exec(ctx.Request().Context(), queryUpdateAlbumMediaItems, mediaItemsCount,
		coverMediaItemID, userID, uid)
	if err != nil {
		slog.Error("error updating album", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err = atx.Commit(ctx.Request().Context()); err != nil {
		slog.Error("error committing transaction for removing album mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

// GetAlbum ...
func (h *Handler) GetAlbum(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	uid, err := getAlbumID(ctx)
	if err != nil {
		return err
	}
	album := models.Album{CoverMediaItem: &models.MediaItem{}}
	err = h.DB.QueryRow(ctx.Request().Context(), queryGetAlbum, userID, uid).Scan(&album.ID, &album.UserID,
		&album.Name, &album.Description, &album.IsShared, &album.IsHidden, &album.MediaItemsCount,
		&album.CoverMediaItemID, &album.CreatedAt, &album.UpdatedAt, &album.CoverMediaItem.ID,
		&album.CoverMediaItem.UserID, &album.CoverMediaItem.Filename, &album.CoverMediaItem.Hash,
		&album.CoverMediaItem.Description, &album.CoverMediaItem.MimeType, &album.CoverMediaItem.SourceURL,
		&album.CoverMediaItem.PreviewURL, &album.CoverMediaItem.ThumbnailURL, &album.CoverMediaItem.Placeholder,
		&album.CoverMediaItem.IsFavourite, &album.CoverMediaItem.IsHidden, &album.CoverMediaItem.IsDeleted,
		&album.CoverMediaItem.Status, &album.CoverMediaItem.MediaItemType, &album.CoverMediaItem.MediaItemCategory,
		&album.CoverMediaItem.Width, &album.CoverMediaItem.Height, &album.CoverMediaItem.CreationTime,
		&album.CoverMediaItem.CameraMake, &album.CoverMediaItem.CameraModel, &album.CoverMediaItem.FocalLength,
		&album.CoverMediaItem.ApertureFnumber, &album.CoverMediaItem.IsoEquivalent, &album.CoverMediaItem.ExposureTime,
		&album.CoverMediaItem.Megapixels, &album.CoverMediaItem.Latitude, &album.CoverMediaItem.Longitude,
		&album.CoverMediaItem.FPS, &album.CoverMediaItem.EXIFData, &album.CoverMediaItem.Keywords,
		&album.CoverMediaItem.CreatedAt, &album.CoverMediaItem.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "album not found")
		}
		slog.Error("error getting album", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, album)
}

// UpdateAlbum ...
func (h *Handler) UpdateAlbum(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	uid, err := getAlbumID(ctx)
	if err != nil {
		return err
	}
	album, err := getAlbum(ctx)
	if err != nil {
		return err
	}
	album.ID = uid
	album.UserID = userID
	album.UpdatedAt = time.Now()
	_, err = h.DB.Exec(ctx.Request().Context(), queryUpdateAlbum, userID, uid, album.Name, album.Description, album.IsShared, album.IsHidden, album.CoverMediaItemID, album.UpdatedAt)
	if err != nil {
		slog.Error("error updating album", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

// DeleteAlbum ...
func (h *Handler) DeleteAlbum(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	uid, err := getAlbumID(ctx)
	if err != nil {
		return err
	}
	_, err = h.DB.Exec(ctx.Request().Context(), queryDeleteAlbum, userID, uid)
	if err != nil {
		slog.Error("error deleting album", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

// GetAlbums ...
func (h *Handler) GetAlbums(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	shared := getAlbumShared(ctx)
	order := getAlbumSortOrder(ctx)
	albums := []models.Album{}
	rows, err := h.DB.Query(ctx.Request().Context(), fmt.Sprintf(queryGetAlbums, order), userID, shared, offset, limit)
	if err != nil {
		slog.Error("error getting albums", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		album, err := models.ScanRowsToAlbum(rows)
		if err != nil {
			slog.Error("error scanning album", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		albums = append(albums, album)
	}

	return ctx.JSON(http.StatusOK, albums)
}

// CreateAlbum ...
func (h *Handler) CreateAlbum(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	album, err := getAlbum(ctx)
	if err != nil {
		return err
	}
	album.ID = uuid.NewV4()
	album.UserID = userID
	_, err = h.DB.Exec(ctx.Request().Context(), queryCreateAlbum, album.ID, album.UserID, album.Name,
		album.Description, album.IsShared, album.IsHidden, album.CreatedAt, album.UpdatedAt)
	if err != nil {
		slog.Error("error creating album", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, album)
}

func getAlbumID(ctx echo.Context) (uuid.UUID, error) {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting album id", "error", err)

		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid album id")
	}

	return uid, err
}

func getMediaItems(ctx echo.Context) ([]*models.MediaItem, error) {
	mediaItemsRequest := new(MediaItemsRequest)
	err := ctx.Bind(mediaItemsRequest)
	if err != nil || len(mediaItemsRequest.MediaItems) == 0 {
		slog.Error("error getting album mediaitems", "error", err)

		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitems")
	}
	mediaItems := make([]*models.MediaItem, len(mediaItemsRequest.MediaItems))
	for idx, mediaItem := range mediaItemsRequest.MediaItems {
		uid, err := uuid.FromString(mediaItem)
		if err != nil {
			slog.Error("error getting album mediaitem id", "error", err)

			return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
		}
		mediaItems[idx] = &models.MediaItem{ID: uid}
	}

	return mediaItems, nil
}

func getAlbum(ctx echo.Context) (*models.Album, error) {
	albumRequest := new(AlbumRequest)
	err := ctx.Bind(albumRequest)
	if err != nil {
		slog.Error("error getting album", "error", err)

		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid album")
	}
	album := models.Album{
		Description: albumRequest.Description, IsShared: albumRequest.IsShared, IsHidden: albumRequest.IsHidden,
	}
	if albumRequest.Name != nil {
		album.Name = *albumRequest.Name
	}
	if albumRequest.CoverMediaItemID != nil {
		coverMediaItemID, err := uuid.FromString(*albumRequest.CoverMediaItemID)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid album cover mediaitem id")
		}
		album.CoverMediaItemID = &coverMediaItemID
	}
	if reflect.DeepEqual(models.Album{}, album) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid album")
	}

	return &album, nil
}

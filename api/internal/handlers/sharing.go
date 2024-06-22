package handlers

import (
	"api/internal/models"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

// GetSharedAlbumMediaItems ...
func (h *Handler) GetSharedAlbumMediaItems(ctx echo.Context) error {
	offset, limit := getOffsetAndLimit(ctx)
	uid, err := getSharedAlbumID(ctx)
	if err != nil {
		return err
	}
	sharedAlbum := new(models.Album)
	sharedAlbum.ID = uid
	mediaItems := []models.MediaItem{}
	err = h.DB.Model(&sharedAlbum).Offset(offset).Limit(limit).Association("MediaItems").Find(&mediaItems, "is_deleted=?", false)
	if err != nil {
		slog.Error("error getting shared album mediaitems", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	result := h.DB.Model(&models.Album{}).
		Where("is_shared=true AND id=?", uid).
		Preload("CoverMediaItem").
		First(&sharedAlbum)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "shared link not found")
		}
		slog.Error("error getting shared album", "error", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
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

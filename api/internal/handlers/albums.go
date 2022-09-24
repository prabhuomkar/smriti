package handlers

import (
	"api/internal/models"
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// GetAlbumMediaItems ...
func (h *Handler) GetAlbumMediaItems(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting album id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid album id")
	}
	mediaItems := []models.MediaItem{}
	err = h.DB.Select(&mediaItems, "SELECT * FROM album_mediaitems "+
		"INNER JOIN mediaitems ON album_mediaitems.mediaitem_id = mediaitems.id "+
		"WHERE album_mediaitems.album_id=$1", uid)
	if err != nil {
		log.Printf("error getting album mediaitems: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// AddAlbumMediaItems ...
func (h *Handler) AddAlbumMediaItems(ctx echo.Context) error {
	return nil
}

// RemoveAlbumMediaItems ...
func (h *Handler) RemoveAlbumMediaItems(ctx echo.Context) error {
	return nil
}

// GetAlbum ...
func (h *Handler) GetAlbum(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting album id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid album id")
	}
	album := models.Album{}
	err = h.DB.Get(&album, "SELECT * FROM albums WHERE id=$1", uid)
	if err != nil {
		log.Printf("error getting album: %+v", err)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "album not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, album)
}

// UpdateAlbum ...
func (h *Handler) UpdateAlbum(ctx echo.Context) error {
	return nil
}

// DeleteAlbum ...
func (h *Handler) DeleteAlbum(ctx echo.Context) error {
	return nil
}

// GetAlbums ...
func (h *Handler) GetAlbums(ctx echo.Context) error {
	albums := []models.Album{}
	err := h.DB.Select(&albums, "SELECT * FROM albums WHERE is_hidden=false")
	if err != nil {
		log.Printf("error getting albums: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, albums)
}

// CreateAlbum ...
func (h *Handler) CreateAlbum(ctx echo.Context) error {
	return nil
}

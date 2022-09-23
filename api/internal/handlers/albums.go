package handlers

import (
	"api/internal/models"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// GetAlbumMediaItems ...
func (h *Handler) GetAlbumMediaItems(ctx echo.Context) error {
	return nil
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
	return nil
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
		return echo.ErrInternalServerError
	}
	return ctx.JSON(http.StatusOK, albums)
}

// CreateAlbum ...
func (h *Handler) CreateAlbum(ctx echo.Context) error {
	return nil
}

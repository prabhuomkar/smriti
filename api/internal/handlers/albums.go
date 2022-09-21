package handlers

import (
	"errors"

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
	return errors.New("some error")
}

// CreateAlbum ...
func (h *Handler) CreateAlbum(ctx echo.Context) error {
	return nil
}

package handlers

import (
	"api/internal/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

type (
	// AlbumMediaItemRequest ...
	AlbumMediaItemRequest struct {
		MediaItems []string `json:"mediaItems" required:"true"`
	}
)

// GetAlbumMediaItems ...
func (h *Handler) GetAlbumMediaItems(ctx echo.Context) error {
	uid, err := getAlbumID(ctx)
	if err != nil {
		return err
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
	uid, err := getAlbumID(ctx)
	if err != nil {
		return err
	}
	mediaItems, err := getAlbumMediaItems(ctx)
	if err != nil {
		return err
	}
	fmt.Println(uid, mediaItems)
	return nil
}

// RemoveAlbumMediaItems ...
func (h *Handler) RemoveAlbumMediaItems(ctx echo.Context) error {
	uid, err := getAlbumID(ctx)
	if err != nil {
		return err
	}
	mediaItems, err := getAlbumMediaItems(ctx)
	if err != nil {
		return err
	}
	fmt.Println(uid, mediaItems)
	return nil
}

// GetAlbum ...
func (h *Handler) GetAlbum(ctx echo.Context) error {
	uid, err := getAlbumID(ctx)
	if err != nil {
		return err
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
	uid, err := getAlbumID(ctx)
	if err != nil {
		return err
	}
	album, err := getAlbum(ctx)
	if err != nil {
		return err
	}
	album.Update()
	albumCoverMediaItemID := uuid.FromStringOrNil(album.CoverMediaItemID)
	_, err = h.DB.Exec("UPDATE albums SET name=$2,description=$3,is_shared=$4,cover_mediaitem_id=$5,"+
		"cover_mediaitem_thumbnail_url=$6,is_hidden=$7,updated_at=$8 WHERE id=$1", uid, album.Name,
		album.Description, album.IsShared, albumCoverMediaItemID, album.CoverMediaItemThumbnailUrl,
		album.IsHidden, album.UpdatedAt)
	if err != nil {
		log.Printf("error updating album: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// DeleteAlbum ...
func (h *Handler) DeleteAlbum(ctx echo.Context) error {
	uid, err := getAlbumID(ctx)
	if err != nil {
		return err
	}
	_, err = h.DB.Exec("DELETE FROM albums WHERE id=$1", uid)
	if err != nil {
		log.Printf("error deleting album: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
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
	album, err := getAlbum(ctx)
	if err != nil {
		return err
	}
	album.NewID()
	album.Create()
	albumCoverMediaItemID := uuid.FromStringOrNil(album.CoverMediaItemID)
	_, err = h.DB.Exec("INSERT INTO albums VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", album.ID, album.Name,
		album.Description, album.IsShared, albumCoverMediaItemID, album.CoverMediaItemThumbnailUrl, 0,
		album.IsHidden, album.CreatedAt, album.UpdatedAt)
	if err != nil {
		log.Printf("error creating album: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, album)
}

func getAlbumID(ctx echo.Context) (uuid.UUID, error) {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting album id: %+v", err)
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid album id")
	}
	return uid, err
}

func getAlbumMediaItems(ctx echo.Context) ([]uuid.UUID, error) {
	albumMediaItemRequest := new(AlbumMediaItemRequest)
	err := ctx.Bind(albumMediaItemRequest)
	if err != nil || len(albumMediaItemRequest.MediaItems) == 0 {
		log.Printf("error getting album mediaitems: %+v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid album mediaitems")
	}
	uids := make([]uuid.UUID, len(albumMediaItemRequest.MediaItems))
	for idx, mediaItem := range albumMediaItemRequest.MediaItems {
		uid, err := uuid.FromString(mediaItem)
		if err != nil {
			log.Printf("error getting album mediaitem id: %+v", err)
			return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
		}
		uids[idx] = uid
	}
	return uids, nil
}

func getAlbum(ctx echo.Context) (*models.Album, error) {
	album := new(models.Album)
	err := ctx.Bind(album)
	if err != nil {
		log.Printf("error getting album: %+v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid album")
	}
	return album, nil
}

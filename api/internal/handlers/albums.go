package handlers

import (
	"api/internal/models"
	"errors"
	"log"
	"net/http"
	"reflect"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type (
	// AlbumRequest ...
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

// GetAlbumMediaItems ...
func (h *Handler) GetAlbumMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	uid, err := getAlbumID(ctx)
	if err != nil {
		return err
	}
	album := new(models.Album)
	album.ID = uid
	album.UserID = userID
	mediaItems := []models.MediaItem{}
	err = h.DB.Model(&album).Offset(offset).Limit(limit).Association("MediaItems").Find(&mediaItems)
	if err != nil {
		log.Printf("error getting album mediaitems: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	album := new(models.Album)
	album.ID = uid
	album.UserID = userID
	err = h.DB.Omit("MediaItems.*").Model(&album).Association("MediaItems").Append(mediaItems)
	if err != nil {
		log.Printf("error adding album mediaitems: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	mediaItemCount := int(h.DB.Model(&album).Association("MediaItems").Count())
	album.MediaItemsCount = &mediaItemCount
	album.CoverMediaItemID = &mediaItems[len(mediaItems)-1].ID
	result := h.DB.Model(&album).Omit("MediaItems").Updates(album)
	if result.Error != nil {
		log.Printf("error updating album: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
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
	album := new(models.Album)
	album.ID = uid
	album.UserID = userID
	err = h.DB.Omit("MediaItems.*").Model(&album).Association("MediaItems").Delete(mediaItems)
	if err != nil {
		log.Printf("error removing album mediaitems: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	newCoverMediaItem := models.MediaItem{}
	err = h.DB.Model(&album).Association("MediaItems").Find(&newCoverMediaItem)
	if err != nil {
		log.Printf("error getting new album cover mediaitem: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	mediaItemCount := int(h.DB.Model(&album).Association("MediaItems").Count())
	album.MediaItemsCount = &mediaItemCount
	album.CoverMediaItemID = &newCoverMediaItem.ID
	if newCoverMediaItem.ID == uuid.Nil {
		album.CoverMediaItemID = nil
	}
	result := h.DB.Model(&album).Omit("MediaItems").Updates(map[string]interface{}{
		"MediaItemsCount":  &mediaItemCount,
		"CoverMediaItemID": album.CoverMediaItemID,
	})
	if result.Error != nil {
		log.Printf("error updating album: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
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
	album := models.Album{}
	result := h.DB.Model(&models.Album{}).
		Where("id=? AND user_id=?", uid, userID).
		Preload("CoverMediaItem").
		First(&album)
	if result.Error != nil {
		log.Printf("error getting album: %+v", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "album not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
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
	result := h.DB.Model(&album).Updates(album)
	if result.Error != nil || result.RowsAffected != 1 {
		log.Printf("error updating album: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
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
	album := models.Album{ID: uid, UserID: userID}
	err = h.DB.Model(&album).Association("MediaItems").Clear()
	if err != nil {
		log.Printf("error deleting album: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	result := h.DB.Delete(&album)
	if result.Error != nil || result.RowsAffected != 1 {
		log.Printf("error deleting album: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// GetAlbums ...
func (h *Handler) GetAlbums(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	order := getAlbumSortOrder(ctx)
	albums := []models.Album{}
	result := h.DB.Model(&models.Album{}).
		Where("is_hidden=false AND user_id=?", userID).
		Preload("CoverMediaItem").
		Order(order).
		Find(&albums).
		Offset(offset).
		Limit(limit)
	if result.Error != nil {
		log.Printf("error getting albums: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
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
	if result := h.DB.Create(&album); result.Error != nil {
		log.Printf("error creating album: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
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

func getMediaItems(ctx echo.Context) ([]*models.MediaItem, error) {
	mediaItemsRequest := new(MediaItemsRequest)
	err := ctx.Bind(mediaItemsRequest)
	if err != nil || len(mediaItemsRequest.MediaItems) == 0 {
		log.Printf("error getting album mediaitems: %+v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitems")
	}
	mediaItems := make([]*models.MediaItem, len(mediaItemsRequest.MediaItems))
	for idx, mediaItem := range mediaItemsRequest.MediaItems {
		uid, err := uuid.FromString(mediaItem)
		if err != nil {
			log.Printf("error getting album mediaitem id: %+v", err)
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
		log.Printf("error getting album: %+v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid album")
	}
	album := models.Album{
		Description: albumRequest.Description,
		IsShared:    albumRequest.IsShared,
		IsHidden:    albumRequest.IsHidden,
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

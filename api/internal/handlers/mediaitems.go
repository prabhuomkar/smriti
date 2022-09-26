package handlers

import (
	"api/internal/models"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// GetMediaItemPlaces ...
func (h *Handler) GetMediaItemPlaces(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting mediaitem id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	places := []models.Place{}
	err = h.DB.Model(&mediaItem).Association("Places").Find(&places)
	if err != nil {
		log.Printf("error getting mediaitem places: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, places)
}

// GetMediaItemThings ...
func (h *Handler) GetMediaItemThings(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting mediaitem id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	things := []models.Thing{}
	err = h.DB.Model(&mediaItem).Association("Things").Find(&things)
	if err != nil {
		log.Printf("error getting mediaitem things: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, things)
}

// GetMediaItemPeople ...
func (h *Handler) GetMediaItemPeople(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting mediaitem id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	people := []models.People{}
	err = h.DB.Model(&mediaItem).Association("People").Find(&people)
	if err != nil {
		log.Printf("error getting mediaitem people: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, people)
}

// GetMediaItem ...
func (h *Handler) GetMediaItem(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting mediaitem id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := models.MediaItem{}
	result := h.DB.Where("id = ?", uid).First(&mediaItem)
	if result.Error != nil {
		log.Printf("error getting mediaitem: %+v", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "mediaitem not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItem)
}

// UpdateMediaItem ...
func (h *Handler) UpdateMediaItem(ctx echo.Context) error {
	return nil
}

// DeleteMediaItem ...
func (h *Handler) DeleteMediaItem(ctx echo.Context) error {
	return nil
}

// GetMediaItems ...
func (h *Handler) GetMediaItems(ctx echo.Context) error {
	mediaItems := []models.MediaItem{}
	result := h.DB.Where("is_hidden=false OR is_deleted=false").Find(&mediaItems)
	if result.Error != nil {
		log.Printf("error getting mediaitems: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// UploadMediaItems ...
func (h *Handler) UploadMediaItems(ctx echo.Context) error {
	uid := uuid.NewV4()
	err := h.mockCreateMediaItem(uid)
	if err != nil {
		log.Printf("error creating mock mediaitem: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	mediaItem := models.MediaItem{}
	result := h.DB.Where("id = ?", uid).First(&mediaItem)
	if result.Error != nil {
		log.Printf("error getting mediaitem: %+v", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "mediaitem not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItem)
}

func (h *Handler) mockCreateMediaItem(uid uuid.UUID) error {
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	mediaItem.Status = models.Ready
	mediaItem.MediaItemType = models.Photo
	result := h.DB.Create(&mediaItem)
	return result.Error
}

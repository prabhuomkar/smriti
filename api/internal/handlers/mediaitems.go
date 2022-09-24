package handlers

import (
	"api/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// GetMediaItemPlaces ...
func (h *Handler) GetMediaItemPlaces(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting mediaitem id: %+v", err)
		return echo.ErrBadRequest
	}
	places := []models.Place{}
	err = h.DB.Select(&places, "SELECT * FROM places "+
		"INNER JOIN place_mediaitems ON places.id = place_mediaitems.place_id "+
		"WHERE place_mediaitems.mediaitem_id=$1", uid)
	if err != nil {
		log.Printf("error getting mediaitem places: %+v", err)
		return echo.ErrInternalServerError
	}
	return ctx.JSON(http.StatusOK, places)
}

// GetMediaItemThings ...
func (h *Handler) GetMediaItemThings(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting mediaitem id: %+v", err)
		return echo.ErrBadRequest
	}
	things := []models.Thing{}
	err = h.DB.Select(&things, "SELECT * FROM things "+
		"INNER JOIN thing_mediaitems ON things.id = thing_mediaitems.thing_id "+
		"WHERE thing_mediaitems.mediaitem_id=$1", uid)
	if err != nil {
		log.Printf("error getting mediaitem things: %+v", err)
		return echo.ErrInternalServerError
	}
	return ctx.JSON(http.StatusOK, things)
}

// GetMediaItemPeople ...
func (h *Handler) GetMediaItemPeople(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting mediaitem id: %+v", err)
		return echo.ErrBadRequest
	}
	people := []models.People{}
	err = h.DB.Select(&people, "SELECT * FROM people "+
		"INNER JOIN people_mediaitems ON people.id = people_mediaitems.people_id "+
		"WHERE people_mediaitems.mediaitem_id=$1", uid)
	if err != nil {
		log.Printf("error getting mediaitem people: %+v", err)
		return echo.ErrInternalServerError
	}
	return ctx.JSON(http.StatusOK, people)
}

// GetMediaItem ...
func (h *Handler) GetMediaItem(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting mediaitem id: %+v", err)
		return echo.ErrBadRequest
	}
	mediaItem := models.MediaItem{}
	err = h.DB.Get(&mediaItem, "SELECT * FROM mediaitems WHERE id=$1", uid)
	if err != nil {
		log.Printf("error getting mediaitem: %+v", err)
		return echo.ErrInternalServerError
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
	err := h.DB.Select(&mediaItems, "SELECT * FROM mediaitems WHERE (is_hidden=false OR is_deleted=false)")
	if err != nil {
		log.Printf("error getting mediaitems: %+v", err)
		return echo.ErrInternalServerError
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// UploadMediaItems ...
func (h *Handler) UploadMediaItems(ctx echo.Context) error {
	uid := uuid.NewV4()
	err := h.mockCreateMediaItem(uid)
	if err != nil {
		log.Printf("error creating mock mediaitem: %+v", err)
		return echo.ErrInternalServerError
	}
	mediaItem := models.MediaItem{}
	err = h.DB.Get(&mediaItem, "SELECT * FROM mediaitems WHERE id=$1", uid)
	if err != nil {
		log.Printf("error getting mediaitem: %+v", err)
		return echo.ErrInternalServerError
	}
	return ctx.JSON(http.StatusOK, mediaItem)
}

func (h *Handler) mockCreateMediaItem(uid uuid.UUID) error {
	_, err := h.DB.Exec(h.DB.Rebind("INSERT into mediaitems VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"),
		uid, "IMG_284.jpg", "A sample image from wedding", "image/jpeg", "", "", "", false, false, false, models.Ready, models.Photo, 720, 480, time.Now(), "", "", "", "", "", "", nil, "", time.Now(), time.Now())

	return err
}

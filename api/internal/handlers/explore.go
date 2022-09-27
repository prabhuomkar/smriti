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

// GetPlaces ...
func (h *Handler) GetPlaces(ctx echo.Context) error {
	places := []models.Place{}
	result := h.DB.Model(&models.Place{}).Where("is_hidden=false").Preload("CoverMediaItem").Find(&places)
	if result.Error != nil {
		log.Printf("error getting places: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, places)
}

// GetPlace ...
func (h *Handler) GetPlace(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting place id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid place id")
	}
	place := models.Place{}
	result := h.DB.Model(&models.Place{}).Where("id = ?", uid).Preload("CoverMediaItem").First(&place)
	if result.Error != nil {
		log.Printf("error getting place: %+v", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "place not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, place)
}

// GetPlaceMediaItems ...
func (h *Handler) GetPlaceMediaItems(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting place id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid place id")
	}
	place := new(models.Place)
	place.ID = uid
	mediaItems := []models.MediaItem{}
	err = h.DB.Model(&place).Association("MediaItems").Find(&mediaItems)
	if err != nil {
		log.Printf("error getting place mediaitems: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// GetThings ...
func (h *Handler) GetThings(ctx echo.Context) error {
	things := []models.Thing{}
	result := h.DB.Model(&models.Thing{}).Where("is_hidden=false").Preload("CoverMediaItem").Find(&things)
	if result.Error != nil {
		log.Printf("error getting things: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, things)
}

// GetThing ...
func (h *Handler) GetThing(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting thing id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid thing id")
	}
	thing := models.Thing{}
	result := h.DB.Model(&models.Thing{}).Where("id = ?", uid).Preload("CoverMediaItem").First(&thing)
	if result.Error != nil {
		log.Printf("error getting thing: %+v", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "thing not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, thing)
}

// GetThingMediaItems ...
func (h *Handler) GetThingMediaItems(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting thing id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid thing id")
	}
	thing := new(models.Thing)
	thing.ID = uid
	mediaItems := []models.MediaItem{}
	err = h.DB.Model(&thing).Association("MediaItems").Find(&mediaItems)
	if err != nil {
		log.Printf("error getting thing mediaitems: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// GetPeople ...
func (h *Handler) GetPeople(ctx echo.Context) error {
	people := []models.People{}
	result := h.DB.Model(&models.People{}).Where("is_hidden=false").Preload("CoverMediaItem").Find(&people)
	if result.Error != nil {
		log.Printf("error getting people: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, people)
}

// GetPerson ...
func (h *Handler) GetPerson(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting person id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid person id")
	}
	person := models.People{}
	result := h.DB.Model(&models.People{}).Where("id = ?", uid).Preload("CoverMediaItem").First(&person)
	if result.Error != nil {
		log.Printf("error getting person: %+v", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "person not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, person)
}

// GetPeopleMediaItems ...
func (h *Handler) GetPeopleMediaItems(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting person id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid people id")
	}
	person := new(models.People)
	person.ID = uid
	mediaItems := []models.MediaItem{}
	err = h.DB.Model(&person).Association("MediaItems").Find(&mediaItems)
	if err != nil {
		log.Printf("error getting people mediaitems: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

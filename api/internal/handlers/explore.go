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
	// PeopleRequest ...
	PeopleRequest struct {
		Name             *string `json:"name"`
		IsHidden         *bool   `json:"hidden"`
		CoverMediaItemID *string `json:"coverMediaItemId"`
	}
)

// GetPlaces ...
func (h *Handler) GetPlaces(ctx echo.Context) error {
	offset, limit := getOffsetAndLimit(ctx)
	places := []models.Place{}
	result := h.DB.Model(&models.Place{}).
		Where("is_hidden=false").
		Preload("CoverMediaItem").
		Find(&places).
		Offset(offset).
		Limit(limit)
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
	offset, limit := getOffsetAndLimit(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting place id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid place id")
	}
	place := new(models.Place)
	place.ID = uid
	mediaItems := []models.MediaItem{}
	err = h.DB.Model(&place).Offset(offset).Limit(limit).Association("MediaItems").Find(&mediaItems)
	if err != nil {
		log.Printf("error getting place mediaitems: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// GetThings ...
func (h *Handler) GetThings(ctx echo.Context) error {
	offset, limit := getOffsetAndLimit(ctx)
	things := []models.Thing{}
	result := h.DB.Model(&models.Thing{}).
		Where("is_hidden=false").
		Preload("CoverMediaItem").
		Find(&things).
		Offset(offset).
		Limit(limit)
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
	offset, limit := getOffsetAndLimit(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting thing id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid thing id")
	}
	thing := new(models.Thing)
	thing.ID = uid
	mediaItems := []models.MediaItem{}
	err = h.DB.Model(&thing).Offset(offset).Limit(limit).Association("MediaItems").Find(&mediaItems)
	if err != nil {
		log.Printf("error getting thing mediaitems: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// UpdatePeople ...
func (h *Handler) UpdatePeople(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting people id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid people id")
	}
	people, err := getPeople(ctx)
	if err != nil {
		return err
	}
	people.ID = uid
	result := h.DB.Model(&people).Updates(people)
	if result.Error != nil || result.RowsAffected != 1 {
		log.Printf("error updating people: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// GetPeople ...
func (h *Handler) GetPeople(ctx echo.Context) error {
	offset, limit := getOffsetAndLimit(ctx)
	people := []models.People{}
	result := h.DB.Model(&models.People{}).
		Where("is_hidden=false").
		Preload("CoverMediaItem").
		Find(&people).
		Offset(offset).
		Limit(limit)
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
	offset, limit := getOffsetAndLimit(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting person id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid people id")
	}
	person := new(models.People)
	person.ID = uid
	mediaItems := []models.MediaItem{}
	err = h.DB.Model(&person).Offset(offset).Limit(limit).Association("MediaItems").Find(&mediaItems)
	if err != nil {
		log.Printf("error getting people mediaitems: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

func getPeople(ctx echo.Context) (*models.People, error) {
	peopleRequest := new(PeopleRequest)
	err := ctx.Bind(peopleRequest)
	if err != nil {
		log.Printf("error getting people: %+v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid people")
	}
	people := models.People{
		IsHidden: peopleRequest.IsHidden,
	}
	if peopleRequest.Name != nil {
		people.Name = *peopleRequest.Name
	}
	if peopleRequest.CoverMediaItemID != nil {
		coverMediaItemID, err := uuid.FromString(*peopleRequest.CoverMediaItemID)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid people cover mediaitem id")
		}
		people.CoverMediaItemID = coverMediaItemID
	}
	if reflect.DeepEqual(models.People{}, people) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid people")
	}
	return &people, nil
}

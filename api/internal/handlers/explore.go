package handlers

import (
	"api/internal/models"
	"log"
	"net/http"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// GetPlaces ...
func (h *Handler) GetPlaces(ctx echo.Context) error {
	places := []models.Place{}
	err := h.DB.Select(&places, "SELECT * FROM places WHERE is_hidden=false")
	if err != nil {
		log.Printf("error getting places: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, places)
}

// GetPlace ...
func (h *Handler) GetPlace(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting place id: %+v", err)
		return echo.ErrBadRequest
	}
	place := models.Place{}
	err = h.DB.Get(&place, "SELECT * FROM places WHERE id=$1", uid)
	if err != nil {
		log.Printf("error getting place: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, place)
}

// GetPlaceMediaItems ...
func (h *Handler) GetPlaceMediaItems(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting place id: %+v", err)
		return echo.ErrBadRequest
	}
	mediaItems := []models.MediaItem{}
	err = h.DB.Select(&mediaItems, "SELECT * FROM place_mediaitems "+
		"INNER JOIN mediaitems ON place_mediaitems.mediaitem_id = mediaitems.id "+
		"WHERE place_mediaitems.place_id=$1 AND (mediaitems.is_hidden=false OR mediaitems.is_deleted=false)", uid)
	if err != nil {
		log.Printf("error getting place mediaitems: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// GetThings ...
func (h *Handler) GetThings(ctx echo.Context) error {
	things := []models.Thing{}
	err := h.DB.Select(&things, "SELECT * FROM things WHERE is_hidden=false")
	if err != nil {
		log.Printf("error getting things: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, things)
}

// GetThing ...
func (h *Handler) GetThing(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting thing id: %+v", err)
		return echo.ErrBadRequest
	}
	thing := models.Thing{}
	err = h.DB.Get(&thing, "SELECT * FROM things WHERE id=$1", uid)
	if err != nil {
		log.Printf("error getting thing: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, thing)
}

// GetThingMediaItems ...
func (h *Handler) GetThingMediaItems(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting thing id: %+v", err)
		return echo.ErrBadRequest
	}
	mediaItems := []models.MediaItem{}
	err = h.DB.Select(&mediaItems, "SELECT * FROM thing_mediaitems "+
		"INNER JOIN mediaitems ON thing_mediaitems.mediaitem_id = mediaitems.id "+
		"WHERE thing_mediaitems.thing_id=$1 AND (mediaitems.is_hidden=false OR mediaitems.is_deleted=false)", uid)
	if err != nil {
		log.Printf("error getting thing mediaitems: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// GetPeople ...
func (h *Handler) GetPeople(ctx echo.Context) error {
	people := []models.People{}
	err := h.DB.Select(&people, "SELECT * FROM people WHERE is_hidden=false")
	if err != nil {
		log.Printf("error getting people: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, people)
}

// GetPerson ...
func (h *Handler) GetPerson(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting person id: %+v", err)
		return echo.ErrBadRequest
	}
	person := models.People{}
	err = h.DB.Get(&person, "SELECT * FROM people WHERE id=$1", uid)
	if err != nil {
		log.Printf("error getting person: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, person)
}

// GetPeopleMediaItems ...
func (h *Handler) GetPeopleMediaItems(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting person id: %+v", err)
		return echo.ErrBadRequest
	}
	mediaItems := []models.MediaItem{}
	err = h.DB.Select(&mediaItems, "SELECT * FROM people_mediaitems "+
		"INNER JOIN mediaitems ON people_mediaitems.mediaitem_id = mediaitems.id "+
		"WHERE people_mediaitems.people_id=$1 AND (mediaitems.is_hidden=false OR mediaitems.is_deleted=false)", uid)
	if err != nil {
		log.Printf("error getting people mediaitems: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

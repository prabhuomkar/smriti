package handlers

import (
	"api/internal/models"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// GetPlaces ...
func (h *Handler) GetPlaces(ctx echo.Context) error {
	places := []models.Place{}
	err := h.DB.Select(&places, "SELECT * FROM places WHERE is_hidden=false")
	if err != nil {
		log.Printf("error getting places: %+v", err)
		return echo.ErrInternalServerError
	}
	return ctx.JSON(http.StatusOK, places)
}

// GetPlace ...
func (h *Handler) GetPlace(ctx echo.Context) error {
	return nil
}

// GetPlaceMediaItems ...
func (h *Handler) GetPlaceMediaItems(ctx echo.Context) error {
	return nil
}

// GetThings ...
func (h *Handler) GetThings(ctx echo.Context) error {
	things := []models.Thing{}
	err := h.DB.Select(&things, "SELECT * FROM things WHERE is_hidden=false")
	if err != nil {
		log.Printf("error getting things: %+v", err)
		return echo.ErrInternalServerError
	}
	return ctx.JSON(http.StatusOK, things)
}

// GetThing ...
func (h *Handler) GetThing(ctx echo.Context) error {
	return nil
}

// GetThingMediaItems ...
func (h *Handler) GetThingMediaItems(ctx echo.Context) error {
	return nil
}

// GetPeople ...
func (h *Handler) GetPeople(ctx echo.Context) error {
	people := []models.People{}
	err := h.DB.Select(&people, "SELECT * FROM people WHERE is_hidden=false")
	if err != nil {
		log.Printf("error getting people: %+v", err)
		return echo.ErrInternalServerError
	}
	return ctx.JSON(http.StatusOK, people)
}

// GetPerson ...
func (h *Handler) GetPerson(ctx echo.Context) error {
	return nil
}

// GetPeopleMediaItems ...
func (h *Handler) GetPeopleMediaItems(ctx echo.Context) error {
	return nil
}

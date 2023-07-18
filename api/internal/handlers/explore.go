package handlers

import (
	"api/internal/models"
	"errors"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

type (
	// PeopleRequest ...
	PeopleRequest struct {
		Name             *string `json:"name"`
		IsHidden         *bool   `json:"hidden"`
		CoverMediaItemID *string `json:"coverMediaItemId"`
	}

	// MemoryMediaItem ...
	MemoryMediaItem struct {
		models.MediaItem
		Year string `json:"year" gorm:"column:creation_year"`
	}
)

// GetYearsAgoMediaItems ...
func (h *Handler) GetYearsAgoMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	month, date, err := getMonthAndDate(ctx)
	if err != nil {
		slog.Error("error getting month and date", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid month and date")
	}
	var memoryMediaItems []MemoryMediaItem
	result := h.DB.Raw("SELECT *, EXTRACT(year FROM creation_time) as creation_year "+
		"FROM mediaitems "+
		"WHERE user_id = ? AND EXTRACT(month FROM creation_time) = ? AND EXTRACT(day FROM creation_time) = ? "+
		"AND EXTRACT(year FROM creation_time) IN (SELECT EXTRACT(year FROM creation_time) FROM mediaitems) "+
		"ORDER BY creation_time", userID, month, date).Scan(&memoryMediaItems)
	if result.Error != nil {
		slog.Error("error getting years ago mediaitems", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, memoryMediaItems)
}

// GetPlaces ...
func (h *Handler) GetPlaces(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	places := []models.Place{}
	result := h.DB.Model(&models.Place{}).
		Where("user_id=? AND is_hidden=false", userID).
		Preload("CoverMediaItem").
		Find(&places).
		Offset(offset).
		Limit(limit)
	if result.Error != nil {
		slog.Error("error getting places", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, places)
}

// GetPlace ...
func (h *Handler) GetPlace(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting place id", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid place id")
	}
	place := models.Place{}
	result := h.DB.Model(&models.Place{}).
		Where("id=? AND user_id=?", uid, userID).
		Preload("CoverMediaItem").
		First(&place)
	if result.Error != nil {
		slog.Error("error getting place", slog.Any("error", result.Error))
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "place not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, place)
}

// GetPlaceMediaItems ...
func (h *Handler) GetPlaceMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting place id", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid place id")
	}
	place := new(models.Place)
	place.ID = uid
	place.UserID = userID
	mediaItems := []models.MediaItem{}
	err = h.DB.Model(&place).Offset(offset).Limit(limit).Association("MediaItems").Find(&mediaItems)
	if err != nil {
		slog.Error("error getting place mediaitems", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// GetThings ...
func (h *Handler) GetThings(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	things := []models.Thing{}
	result := h.DB.Model(&models.Thing{}).
		Where("user_id=? AND is_hidden=false", userID).
		Preload("CoverMediaItem").
		Find(&things).
		Offset(offset).
		Limit(limit)
	if result.Error != nil {
		slog.Error("error getting things", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, things)
}

// GetThing ...
func (h *Handler) GetThing(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting thing id", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid thing id")
	}
	thing := models.Thing{}
	result := h.DB.Model(&models.Thing{}).
		Where("id=? AND user_id=?", uid, userID).
		Preload("CoverMediaItem").
		First(&thing)
	if result.Error != nil {
		slog.Error("error getting thing", slog.Any("error", result.Error))
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "thing not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, thing)
}

// GetThingMediaItems ...
func (h *Handler) GetThingMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting thing id", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid thing id")
	}
	thing := new(models.Thing)
	thing.ID = uid
	thing.UserID = userID
	mediaItems := []models.MediaItem{}
	err = h.DB.Model(&thing).Offset(offset).Limit(limit).Association("MediaItems").Find(&mediaItems)
	if err != nil {
		slog.Error("error getting thing mediaitems", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// UpdatePerson ...
func (h *Handler) UpdatePerson(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting people id", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid people id")
	}
	people, err := getPeople(ctx)
	if err != nil {
		return err
	}
	people.ID = uid
	people.UserID = userID
	result := h.DB.Model(&people).Updates(people)
	if result.Error != nil || result.RowsAffected != 1 {
		slog.Error("error updating people", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// GetPeople ...
func (h *Handler) GetPeople(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	people := []models.People{}
	result := h.DB.Model(&models.People{}).
		Where("user_id=? AND is_hidden=false", userID).
		Preload("CoverMediaItem").
		Find(&people).
		Offset(offset).
		Limit(limit)
	if result.Error != nil {
		slog.Error("error getting people", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, people)
}

// GetPerson ...
func (h *Handler) GetPerson(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting person id", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid person id")
	}
	person := models.People{}
	result := h.DB.Model(&models.People{}).
		Where("id=? AND user_id=?", uid, userID).
		Preload("CoverMediaItem").
		First(&person)
	if result.Error != nil {
		slog.Error("error getting person", slog.Any("error", result.Error))
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "person not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, person)
}

// GetPeopleMediaItems ...
func (h *Handler) GetPeopleMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting person id", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid people id")
	}
	person := new(models.People)
	person.ID = uid
	person.UserID = userID
	mediaItems := []models.MediaItem{}
	err = h.DB.Model(&person).Offset(offset).Limit(limit).Association("MediaItems").Find(&mediaItems)
	if err != nil {
		slog.Error("error getting people mediaitems", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

func getPeople(ctx echo.Context) (*models.People, error) {
	peopleRequest := new(PeopleRequest)
	err := ctx.Bind(peopleRequest)
	if err != nil {
		slog.Error("error getting people", slog.Any("error", err))
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
		people.CoverMediaItemID = &coverMediaItemID
	}
	if reflect.DeepEqual(models.People{}, people) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid people")
	}
	return &people, nil
}

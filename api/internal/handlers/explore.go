package handlers

import (
	"api/internal/models"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type ( // PeopleRequest ...
	PeopleRequest struct {
		Name             *string `json:"name"`
		IsHidden         *bool   `json:"hidden"`
		CoverMediaItemID *string `json:"coverMediaItemId"`
	}

	// MemoryMediaItem ...
	MemoryMediaItem struct {
		models.MediaItem

		Year string `json:"year"`
	}
)

const (
	queryGetYearsAgoMediaItems = `SELECT *, EXTRACT(year FROM creation_time) as creation_year FROM mediaitems WHERE user_id=$1 AND` +
		` EXTRACT(month FROM creation_time)=$2 AND EXTRACT(day FROM creation_time)=$3 AND EXTRACT(year FROM creation_time)` +
		` IN (SELECT EXTRACT(year FROM creation_time) FROM mediaitems) ORDER BY creation_time`
	queryGetPlace = `SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,` +
		` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places p LEFT JOIN mediaitems m` +
		` ON p.cover_mediaitem_id=m.id WHERE p.user_id=$1 AND p.id=$2`
	queryGetPlaces = `SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,` +
		` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places p LEFT JOIN mediaitems m` +
		` ON p.cover_mediaitem_id=m.id WHERE p.user_id=$1 AND p.is_hidden=false ORDER BY p.created_at DESC` +
		` OFFSET $2 LIMIT $3`
	queryGetPlaceMediaItems = `SELECT * FROM mediaitems WHERE id IN (SELECT mediaitem_id FROM place_mediaitems` +
		` WHERE user_id=$1 AND place_id=$2) AND is_hidden=false ORDER BY created_at DESC OFFSET $3 LIMIT $4`
	queryGetPerson = `SELECT p.*, mf.* FROM people p LEFT JOIN mediaitem_faces mf ON p.cover_mediaitem_face_id=mf.id` +
		` WHERE p.user_id=$1 AND p.id=$2`
	queryGetPeople = `SELECT p.*, mf.* FROM people p LEFT JOIN mediaitem_faces mf ON p.cover_mediaitem_face_id=mf.id` +
		` WHERE p.user_id=$1 AND p.is_hidden=false ORDER BY p.created_at DESC OFFSET $2 LIMIT $3`
	queryGetPersonMediaItems = `SELECT * FROM mediaitems WHERE id IN (SELECT mediaitem_id FROM people_mediaitems` +
		` WHERE user_id=$1 AND people_id=$2) AND is_hidden=false ORDER BY created_at DESC OFFSET $3 LIMIT $4`
	queryUpdatePerson = `UPDATE people SET name=$3, is_hidden=$4, cover_mediaitem_id=$5, cover_mediaitem_face_id=$6,` +
		` updated_at=$7 WHERE user_id=$1 AND id=$2`
)

// GetYearsAgoMediaItems ...
func (h *Handler) GetYearsAgoMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	month, date, err := getMonthAndDate(ctx)
	if err != nil {
		slog.Error("error getting month and date", "error", err)

		return echo.NewHTTPError(http.StatusBadRequest, "invalid month and date")
	}
	var memoryMediaItems []MemoryMediaItem
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetYearsAgoMediaItems, userID, month, date)
	if err != nil {
		slog.Error("error getting years ago mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var memoryItem MemoryMediaItem
		err = rows.Scan(&memoryItem.ID, &memoryItem.UserID, &memoryItem.Filename, &memoryItem.Hash,
			&memoryItem.Description, &memoryItem.MimeType, &memoryItem.SourceURL, &memoryItem.PreviewURL,
			&memoryItem.ThumbnailURL, &memoryItem.Placeholder, &memoryItem.IsFavourite, &memoryItem.IsHidden,
			&memoryItem.IsDeleted, &memoryItem.Status, &memoryItem.MediaItemType, &memoryItem.MediaItemCategory,
			&memoryItem.Width, &memoryItem.Height, &memoryItem.CreationTime, &memoryItem.CameraMake,
			&memoryItem.CameraModel, &memoryItem.FocalLength, &memoryItem.ApertureFnumber, &memoryItem.IsoEquivalent,
			&memoryItem.ExposureTime, &memoryItem.Megapixels, &memoryItem.Latitude, &memoryItem.Longitude,
			&memoryItem.FPS, &memoryItem.EXIFData, &memoryItem.DetectedText, &memoryItem.Caption, &memoryItem.CreatedAt,
			&memoryItem.UpdatedAt, &memoryItem.Year)
		if err != nil {
			slog.Error("error scanning years ago mediaitem", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		memoryMediaItems = append(memoryMediaItems, memoryItem)
	}

	return ctx.JSON(http.StatusOK, memoryMediaItems)
}

// GetPlaces ...
func (h *Handler) GetPlaces(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	places := []models.Place{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetPlaces, userID, offset, limit)
	if err != nil {
		slog.Error("error getting places", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		place, err := models.ScanRowsToPlace(rows)
		if err != nil {
			slog.Error("error scanning place", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		places = append(places, place)
	}

	return ctx.JSON(http.StatusOK, places)
}

// GetPlace ...
func (h *Handler) GetPlace(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error getting place id", "error", err)

		return echo.NewHTTPError(http.StatusBadRequest, "invalid place id")
	}
	place := models.Place{}
	coverMediaItem := &models.CoverMediaItem{}

	err = h.DB.QueryRow(ctx.Request().Context(), queryGetPlace, userID, uid).Scan(&place.ID, &place.UserID,
		&place.Name, &place.Postcode, &place.Country, &place.Locality, &place.Area, &place.IsHidden,
		&place.CoverMediaItemID, &place.CreatedAt, &place.UpdatedAt, &coverMediaItem.ID,
		&coverMediaItem.UserID, &coverMediaItem.SourceURL, &coverMediaItem.PreviewURL,
		&coverMediaItem.ThumbnailURL, &coverMediaItem.Placeholder, &coverMediaItem.MediaItemType,
		&coverMediaItem.MediaItemCategory, &coverMediaItem.Width, &coverMediaItem.Height)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "place not found")
		}
		slog.Error("error getting place", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if place.CoverMediaItemID != nil {
		place.CoverMediaItem = coverMediaItem
	}

	return ctx.JSON(http.StatusOK, place)
}

// GetPlaceMediaItems ...
func (h *Handler) GetPlaceMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error getting place id", "error", err)

		return echo.NewHTTPError(http.StatusBadRequest, "invalid place id")
	}
	mediaItems := []models.MediaItem{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetPlaceMediaItems, userID, uid, offset, limit)
	if err != nil {
		slog.Error("error getting place mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		mediaItem, err := models.ScanRowsToMediaItem(rows)
		if err != nil {
			slog.Error("error scanning place mediaitem", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		mediaItems = append(mediaItems, mediaItem)
	}

	return ctx.JSON(http.StatusOK, mediaItems)
}

// UpdatePerson ...
func (h *Handler) UpdatePerson(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error getting people id", "error", err)

		return echo.NewHTTPError(http.StatusBadRequest, "invalid people id")
	}
	people, err := getPeople(ctx)
	if err != nil {
		return err
	}
	people.ID = uid
	people.UserID = userID
	people.UpdatedAt = time.Now()
	_, err = h.DB.Exec(ctx.Request().Context(), queryUpdatePerson, userID, uid, people.Name, people.IsHidden,
		people.CoverMediaItemID, people.CoverMediaItemFaceID, people.UpdatedAt)
	if err != nil {
		slog.Error("error updating person", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

// GetPeople ...
func (h *Handler) GetPeople(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	people := []models.People{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetPeople, userID, offset, limit)
	if err != nil {
		slog.Error("error getting people", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		person, err := models.ScanRowsToPerson(rows)
		if err != nil {
			slog.Error("error scanning person", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		people = append(people, person)
	}

	return ctx.JSON(http.StatusOK, people)
}

// GetPerson ...
func (h *Handler) GetPerson(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error getting person id", "error", err)

		return echo.NewHTTPError(http.StatusBadRequest, "invalid person id")
	}
	person := models.People{CoverMediaItemFace: &models.MediaitemFace{}}
	err = h.DB.QueryRow(ctx.Request().Context(), queryGetPerson, userID, uid).Scan(&person.ID, &person.UserID,
		&person.Name, &person.IsHidden, &person.CoverMediaItemID, &person.CoverMediaItemFaceID, &person.CreatedAt,
		&person.UpdatedAt, &person.CoverMediaItemFace.ID, &person.CoverMediaItemFace.MediaitemID,
		&person.CoverMediaItemFace.PeopleID, &person.CoverMediaItemFace.Embedding, &person.CoverMediaItemFace.Thumbnail)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "person not found")
		}
		slog.Error("error getting person", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, person)
}

// GetPersonMediaItems ...
func (h *Handler) GetPersonMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	id := ctx.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		slog.Error("error getting person id", "error", err)

		return echo.NewHTTPError(http.StatusBadRequest, "invalid people id")
	}
	mediaItems := []models.MediaItem{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetPersonMediaItems, userID, uid, offset, limit)
	if err != nil {
		slog.Error("error getting person mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		mediaItem, err := models.ScanRowsToMediaItem(rows)
		if err != nil {
			slog.Error("error scanning person mediaitem", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		mediaItems = append(mediaItems, mediaItem)
	}

	return ctx.JSON(http.StatusOK, mediaItems)
}

func getPeople(ctx echo.Context) (*models.People, error) {
	peopleRequest := new(PeopleRequest)
	err := ctx.Bind(peopleRequest)
	if err != nil {
		slog.Error("error getting people", "error", err)

		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid people")
	}
	people := models.People{IsHidden: peopleRequest.IsHidden}
	if peopleRequest.Name != nil {
		people.Name = *peopleRequest.Name
	}
	if peopleRequest.CoverMediaItemID != nil {
		coverMediaItemID, err := uuid.Parse(*peopleRequest.CoverMediaItemID)
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

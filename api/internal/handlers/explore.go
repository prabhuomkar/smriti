package handlers

import (
	"api/internal/models"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
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
		Year string `json:"year"`
	}
)

const (
	queryGetYearsAgoMediaItems = `SELECT *, EXTRACT(year FROM creation_time) as creation_year FROM mediaitems WHERE user_id=$1 AND EXTRACT(month FROM creation_time)=$2 AND EXTRACT(day FROM creation_time)=$3 AND EXTRACT(year FROM creation_time) IN (SELECT EXTRACT(year FROM creation_time) FROM mediaitems) ORDER BY creation_time`
	queryGetPlace              = `SELECT p.*, m.* FROM places p LEFT JOIN mediaitems m ON p.cover_mediaitem_id=m.id WHERE p.user_id=$1 AND p.id=$2 GROUP BY p.id, m.id`
	queryGetPlaces             = `SELECT p.*, m.* FROM places p LEFT JOIN mediaitems m ON p.cover_mediaitem_id=m.id WHERE p.user_id=$1 AND p.is_hidden=false GROUP BY p.id, m.id ORDER BY p.created_at DESC OFFSET $2 LIMIT $3`
	queryGetPlaceMediaItems    = `SELECT * FROM mediaitems WHERE id IN (SELECT mediaitem_id FROM place_mediaitems WHERE user_id=$1 AND place_id=$2) AND is_hidden=false ORDER BY created_at DESC OFFSET $3 LIMIT $4`
	queryGetThing              = `SELECT t.*, m.* FROM things t LEFT JOIN mediaitems m ON t.cover_mediaitem_id=m.id WHERE t.user_id=$1 AND t.id=$2 GROUP BY t.id, m.id`
	queryGetThings             = `SELECT t.*, m.* FROM things t LEFT JOIN mediaitems m ON t.cover_mediaitem_id=m.id WHERE t.user_id=$1 AND t.is_hidden=false GROUP BY t.id, m.id ORDER BY t.created_at DESC OFFSET $2 LIMIT $3`
	queryGetThingMediaItems    = `SELECT * FROM mediaitems WHERE id IN (SELECT mediaitem_id FROM thing_mediaitems WHERE user_id=$1 AND thing_id=$2) AND is_hidden=false ORDER BY created_at DESC OFFSET $3 LIMIT $4`
	queryGetPerson             = `SELECT p.*, mf.* FROM people p LEFT JOIN mediaitem_faces mf ON p.cover_mediaitem_face_id=mf.id WHERE p.user_id=$1 AND p.id=$2 GROUP BY p.id, mf.id`
	queryGetPeople             = `SELECT p.*, mf.* FROM people p LEFT JOIN mediaitem_faces mf ON p.cover_mediaitem_face_id=mf.id WHERE p.user_id=$1 AND p.is_hidden=false GROUP BY p.id, mf.id ORDER BY p.created_at DESC OFFSET $2 LIMIT $3`
	queryGetPersonMediaItems   = `SELECT * FROM mediaitems WHERE id IN (SELECT mediaitem_id FROM people_mediaitems WHERE user_id=$1 AND people_id=$2) AND is_hidden=false ORDER BY created_at DESC OFFSET $3 LIMIT $4`
	queryUpdatePerson          = `UPDATE people SET name=$3, is_hidden=$4, cover_mediaitem_id=$5, cover_mediaitem_face_id=$6, updated_at=$7 WHERE user_id=$1 AND id=$2`
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
		err = rows.Scan(&memoryItem.ID,
			&memoryItem.UserID,
			&memoryItem.Filename,
			&memoryItem.Hash,
			&memoryItem.Description,
			&memoryItem.MimeType,
			&memoryItem.SourceURL,
			&memoryItem.PreviewURL,
			&memoryItem.ThumbnailURL,
			&memoryItem.Placeholder,
			&memoryItem.IsFavourite,
			&memoryItem.IsHidden,
			&memoryItem.IsDeleted,
			&memoryItem.Status,
			&memoryItem.MediaItemType,
			&memoryItem.MediaItemCategory,
			&memoryItem.Width,
			&memoryItem.Height,
			&memoryItem.CreationTime,
			&memoryItem.CameraMake,
			&memoryItem.CameraModel,
			&memoryItem.FocalLength,
			&memoryItem.ApertureFnumber,
			&memoryItem.IsoEquivalent,
			&memoryItem.ExposureTime,
			&memoryItem.Latitude,
			&memoryItem.Longitude,
			&memoryItem.FPS,
			&memoryItem.EXIFData,
			&memoryItem.Keywords,
			&memoryItem.CreatedAt,
			&memoryItem.UpdatedAt,
			&memoryItem.Year)
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
		place := models.Place{CoverMediaItem: &models.MediaItem{}}
		err = rows.Scan(&place.ID,
			&place.UserID,
			&place.Name,
			&place.Postcode,
			&place.Country,
			&place.Locality,
			&place.Area,
			&place.IsHidden,
			&place.CoverMediaItemID,
			&place.CreatedAt,
			&place.UpdatedAt,
			&place.CoverMediaItem.ID,
			&place.CoverMediaItem.UserID,
			&place.CoverMediaItem.Filename,
			&place.CoverMediaItem.Hash,
			&place.CoverMediaItem.Description,
			&place.CoverMediaItem.MimeType,
			&place.CoverMediaItem.SourceURL,
			&place.CoverMediaItem.PreviewURL,
			&place.CoverMediaItem.ThumbnailURL,
			&place.CoverMediaItem.Placeholder,
			&place.CoverMediaItem.IsFavourite,
			&place.CoverMediaItem.IsHidden,
			&place.CoverMediaItem.IsDeleted,
			&place.CoverMediaItem.Status,
			&place.CoverMediaItem.MediaItemType,
			&place.CoverMediaItem.MediaItemCategory,
			&place.CoverMediaItem.Width,
			&place.CoverMediaItem.Height,
			&place.CoverMediaItem.CreationTime,
			&place.CoverMediaItem.CameraMake,
			&place.CoverMediaItem.CameraModel,
			&place.CoverMediaItem.FocalLength,
			&place.CoverMediaItem.ApertureFnumber,
			&place.CoverMediaItem.IsoEquivalent,
			&place.CoverMediaItem.ExposureTime,
			&place.CoverMediaItem.Latitude,
			&place.CoverMediaItem.Longitude,
			&place.CoverMediaItem.FPS,
			&place.CoverMediaItem.EXIFData,
			&place.CoverMediaItem.Keywords,
			&place.CoverMediaItem.CreatedAt,
			&place.CoverMediaItem.UpdatedAt,
		)
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
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting place id", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid place id")
	}
	place := models.Place{CoverMediaItem: &models.MediaItem{}}
	err = h.DB.QueryRow(ctx.Request().Context(), queryGetPlace, userID, uid).Scan(&place.ID,
		&place.UserID,
		&place.Name,
		&place.Postcode,
		&place.Country,
		&place.Locality,
		&place.Area,
		&place.IsHidden,
		&place.CoverMediaItemID,
		&place.CreatedAt,
		&place.UpdatedAt,
		&place.CoverMediaItem.ID,
		&place.CoverMediaItem.UserID,
		&place.CoverMediaItem.Filename,
		&place.CoverMediaItem.Hash,
		&place.CoverMediaItem.Description,
		&place.CoverMediaItem.MimeType,
		&place.CoverMediaItem.SourceURL,
		&place.CoverMediaItem.PreviewURL,
		&place.CoverMediaItem.ThumbnailURL,
		&place.CoverMediaItem.Placeholder,
		&place.CoverMediaItem.IsFavourite,
		&place.CoverMediaItem.IsHidden,
		&place.CoverMediaItem.IsDeleted,
		&place.CoverMediaItem.Status,
		&place.CoverMediaItem.MediaItemType,
		&place.CoverMediaItem.MediaItemCategory,
		&place.CoverMediaItem.Width,
		&place.CoverMediaItem.Height,
		&place.CoverMediaItem.CreationTime,
		&place.CoverMediaItem.CameraMake,
		&place.CoverMediaItem.CameraModel,
		&place.CoverMediaItem.FocalLength,
		&place.CoverMediaItem.ApertureFnumber,
		&place.CoverMediaItem.IsoEquivalent,
		&place.CoverMediaItem.ExposureTime,
		&place.CoverMediaItem.Latitude,
		&place.CoverMediaItem.Longitude,
		&place.CoverMediaItem.FPS,
		&place.CoverMediaItem.EXIFData,
		&place.CoverMediaItem.Keywords,
		&place.CoverMediaItem.CreatedAt,
		&place.CoverMediaItem.UpdatedAt,
	)
	if err != nil {
		slog.Error("error getting place", "error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "place not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
		mediaItem := models.MediaItem{}
		err = rows.Scan(&mediaItem.ID,
			&mediaItem.UserID,
			&mediaItem.Filename,
			&mediaItem.Hash,
			&mediaItem.Description,
			&mediaItem.MimeType,
			&mediaItem.SourceURL,
			&mediaItem.PreviewURL,
			&mediaItem.ThumbnailURL,
			&mediaItem.Placeholder,
			&mediaItem.IsFavourite,
			&mediaItem.IsHidden,
			&mediaItem.IsDeleted,
			&mediaItem.Status,
			&mediaItem.MediaItemType,
			&mediaItem.MediaItemCategory,
			&mediaItem.Width,
			&mediaItem.Height,
			&mediaItem.CreationTime,
			&mediaItem.CameraMake,
			&mediaItem.CameraModel,
			&mediaItem.FocalLength,
			&mediaItem.ApertureFnumber,
			&mediaItem.IsoEquivalent,
			&mediaItem.ExposureTime,
			&mediaItem.Latitude,
			&mediaItem.Longitude,
			&mediaItem.FPS,
			&mediaItem.EXIFData,
			&mediaItem.Keywords,
			&mediaItem.CreatedAt,
			&mediaItem.UpdatedAt)
		if err != nil {
			slog.Error("error scanning place mediaitem", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		mediaItems = append(mediaItems, mediaItem)
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// GetThings ...
func (h *Handler) GetThings(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	things := []models.Thing{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetThings, userID, offset, limit)
	if err != nil {
		slog.Error("error getting things", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		thing := models.Thing{CoverMediaItem: &models.MediaItem{}}
		err = rows.Scan(&thing.ID,
			&thing.UserID,
			&thing.Name,
			&thing.IsHidden,
			&thing.CoverMediaItemID,
			&thing.CreatedAt,
			&thing.UpdatedAt,
			&thing.CoverMediaItem.ID,
			&thing.CoverMediaItem.UserID,
			&thing.CoverMediaItem.Filename,
			&thing.CoverMediaItem.Hash,
			&thing.CoverMediaItem.Description,
			&thing.CoverMediaItem.MimeType,
			&thing.CoverMediaItem.SourceURL,
			&thing.CoverMediaItem.PreviewURL,
			&thing.CoverMediaItem.ThumbnailURL,
			&thing.CoverMediaItem.Placeholder,
			&thing.CoverMediaItem.IsFavourite,
			&thing.CoverMediaItem.IsHidden,
			&thing.CoverMediaItem.IsDeleted,
			&thing.CoverMediaItem.Status,
			&thing.CoverMediaItem.MediaItemType,
			&thing.CoverMediaItem.MediaItemCategory,
			&thing.CoverMediaItem.Width,
			&thing.CoverMediaItem.Height,
			&thing.CoverMediaItem.CreationTime,
			&thing.CoverMediaItem.CameraMake,
			&thing.CoverMediaItem.CameraModel,
			&thing.CoverMediaItem.FocalLength,
			&thing.CoverMediaItem.ApertureFnumber,
			&thing.CoverMediaItem.IsoEquivalent,
			&thing.CoverMediaItem.ExposureTime,
			&thing.CoverMediaItem.Latitude,
			&thing.CoverMediaItem.Longitude,
			&thing.CoverMediaItem.FPS,
			&thing.CoverMediaItem.EXIFData,
			&thing.CoverMediaItem.Keywords,
			&thing.CoverMediaItem.CreatedAt,
			&thing.CoverMediaItem.UpdatedAt,
		)
		if err != nil {
			slog.Error("error scanning thing", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		things = append(things, thing)
	}
	return ctx.JSON(http.StatusOK, things)
}

// GetThing ...
func (h *Handler) GetThing(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting thing id", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid thing id")
	}
	thing := models.Thing{CoverMediaItem: &models.MediaItem{}}
	err = h.DB.QueryRow(ctx.Request().Context(), queryGetThing, userID, uid).Scan(&thing.ID,
		&thing.UserID,
		&thing.Name,
		&thing.IsHidden,
		&thing.CoverMediaItemID,
		&thing.CreatedAt,
		&thing.UpdatedAt,
		&thing.CoverMediaItem.ID,
		&thing.CoverMediaItem.UserID,
		&thing.CoverMediaItem.Filename,
		&thing.CoverMediaItem.Hash,
		&thing.CoverMediaItem.Description,
		&thing.CoverMediaItem.MimeType,
		&thing.CoverMediaItem.SourceURL,
		&thing.CoverMediaItem.PreviewURL,
		&thing.CoverMediaItem.ThumbnailURL,
		&thing.CoverMediaItem.Placeholder,
		&thing.CoverMediaItem.IsFavourite,
		&thing.CoverMediaItem.IsHidden,
		&thing.CoverMediaItem.IsDeleted,
		&thing.CoverMediaItem.Status,
		&thing.CoverMediaItem.MediaItemType,
		&thing.CoverMediaItem.MediaItemCategory,
		&thing.CoverMediaItem.Width,
		&thing.CoverMediaItem.Height,
		&thing.CoverMediaItem.CreationTime,
		&thing.CoverMediaItem.CameraMake,
		&thing.CoverMediaItem.CameraModel,
		&thing.CoverMediaItem.FocalLength,
		&thing.CoverMediaItem.ApertureFnumber,
		&thing.CoverMediaItem.IsoEquivalent,
		&thing.CoverMediaItem.ExposureTime,
		&thing.CoverMediaItem.Latitude,
		&thing.CoverMediaItem.Longitude,
		&thing.CoverMediaItem.FPS,
		&thing.CoverMediaItem.EXIFData,
		&thing.CoverMediaItem.Keywords,
		&thing.CoverMediaItem.CreatedAt,
		&thing.CoverMediaItem.UpdatedAt,
	)
	if err != nil {
		slog.Error("error getting thing", "error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "thing not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
		slog.Error("error getting thing id", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid thing id")
	}
	mediaItems := []models.MediaItem{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetThingMediaItems, userID, uid, offset, limit)
	if err != nil {
		slog.Error("error getting thing mediaitems", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		mediaItem := models.MediaItem{}
		err = rows.Scan(&mediaItem.ID,
			&mediaItem.UserID,
			&mediaItem.Filename,
			&mediaItem.Hash,
			&mediaItem.Description,
			&mediaItem.MimeType,
			&mediaItem.SourceURL,
			&mediaItem.PreviewURL,
			&mediaItem.ThumbnailURL,
			&mediaItem.Placeholder,
			&mediaItem.IsFavourite,
			&mediaItem.IsHidden,
			&mediaItem.IsDeleted,
			&mediaItem.Status,
			&mediaItem.MediaItemType,
			&mediaItem.MediaItemCategory,
			&mediaItem.Width,
			&mediaItem.Height,
			&mediaItem.CreationTime,
			&mediaItem.CameraMake,
			&mediaItem.CameraModel,
			&mediaItem.FocalLength,
			&mediaItem.ApertureFnumber,
			&mediaItem.IsoEquivalent,
			&mediaItem.ExposureTime,
			&mediaItem.Latitude,
			&mediaItem.Longitude,
			&mediaItem.FPS,
			&mediaItem.EXIFData,
			&mediaItem.Keywords,
			&mediaItem.CreatedAt,
			&mediaItem.UpdatedAt)
		if err != nil {
			slog.Error("error scanning thing mediaitem", "error", err)
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
	uid, err := uuid.FromString(id)
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
	result, err := h.DB.Exec(ctx.Request().Context(), queryUpdatePerson, userID, uid,
		people.Name, people.IsHidden, people.CoverMediaItemID, people.CoverMediaItemFaceID, people.UpdatedAt)
	if !result.Update() || err != nil {
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
		person := models.People{CoverMediaItemFace: &models.MediaitemFace{}}
		err = rows.Scan(&person.ID,
			&person.UserID,
			&person.Name,
			&person.IsHidden,
			&person.CoverMediaItemID,
			&person.CoverMediaItemFaceID,
			&person.CreatedAt,
			&person.UpdatedAt,
			&person.CoverMediaItemFace.ID,
			&person.CoverMediaItemFace.MediaitemID,
			&person.CoverMediaItemFace.PeopleID,
			&person.CoverMediaItemFace.Embedding,
			&person.CoverMediaItemFace.Thumbnail,
		)
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
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting person id", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid person id")
	}
	person := models.People{CoverMediaItemFace: &models.MediaitemFace{}}
	err = h.DB.QueryRow(ctx.Request().Context(), queryGetPerson, userID, uid).Scan(&person.ID,
		&person.UserID,
		&person.Name,
		&person.IsHidden,
		&person.CoverMediaItemID,
		&person.CoverMediaItemFaceID,
		&person.CreatedAt,
		&person.UpdatedAt,
		&person.CoverMediaItemFace.ID,
		&person.CoverMediaItemFace.MediaitemID,
		&person.CoverMediaItemFace.PeopleID,
		&person.CoverMediaItemFace.Embedding,
		&person.CoverMediaItemFace.Thumbnail,
	)
	if err != nil {
		slog.Error("error getting person", "error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "person not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, person)
}

// GetPersonMediaItems ...
func (h *Handler) GetPersonMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
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
		mediaItem := models.MediaItem{}
		err = rows.Scan(&mediaItem.ID,
			&mediaItem.UserID,
			&mediaItem.Filename,
			&mediaItem.Hash,
			&mediaItem.Description,
			&mediaItem.MimeType,
			&mediaItem.SourceURL,
			&mediaItem.PreviewURL,
			&mediaItem.ThumbnailURL,
			&mediaItem.Placeholder,
			&mediaItem.IsFavourite,
			&mediaItem.IsHidden,
			&mediaItem.IsDeleted,
			&mediaItem.Status,
			&mediaItem.MediaItemType,
			&mediaItem.MediaItemCategory,
			&mediaItem.Width,
			&mediaItem.Height,
			&mediaItem.CreationTime,
			&mediaItem.CameraMake,
			&mediaItem.CameraModel,
			&mediaItem.FocalLength,
			&mediaItem.ApertureFnumber,
			&mediaItem.IsoEquivalent,
			&mediaItem.ExposureTime,
			&mediaItem.Latitude,
			&mediaItem.Longitude,
			&mediaItem.FPS,
			&mediaItem.EXIFData,
			&mediaItem.Keywords,
			&mediaItem.CreatedAt,
			&mediaItem.UpdatedAt)
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

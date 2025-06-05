package handlers

import (
	"api/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

const (
	HeaderUploadType         = "X-Smriti-Upload-Type"    // for resumable
	HeaderUploadCommand      = "X-Smriti-Upload-Command" // start, continue, finish
	HeaderUploadChunkOffset  = "X-Smriti-Upload-Chunk-Offset"
	HeaderUploadChunkSession = "X-Smriti-Upload-Chunk-Session"

	fileFlag       = os.O_WRONLY | os.O_APPEND | os.O_CREATE
	filePermission = 0o644

	queryGetMediaItemPlaces         = `SELECT p.*, m.* FROM places p LEFT JOIN mediaitems m ON p.cover_mediaitem_id=m.id WHERE p.user_id=$1 AND p.is_hidden=false AND p.id IN (SELECT place_id FROM place_mediaitems WHERE mediaitem_id=$2) GROUP BY p.id, m.id ORDER BY p.created_at DESC`
	queryGetMediaItemThings         = `SELECT t.*, m.* FROM things t LEFT JOIN mediaitems m ON t.cover_mediaitem_id=m.id WHERE t.user_id=$1 AND t.is_hidden=false AND t.id IN (SELECT thing_id FROM thing_mediaitems WHERE mediaitem_id=$2) GROUP BY t.id, m.id ORDER BY t.created_at DESC`
	queryGetMediaItemPeople         = `SELECT p.*, mf.* FROM people p LEFT JOIN mediaitem_faces mf ON p.cover_mediaitem_face_id=mf.id WHERE p.user_id=$1 AND p.is_hidden=false AND p.id IN (SELECT people_id FROM people_mediaitems WHERE mediaitem_id=$2) GROUP BY p.id, mf.id ORDER BY p.created_at DESC`
	queryGetMediaItemAlbums         = `SELECT a.*, m.* FROM albums a LEFT JOIN mediaitems m ON a.cover_mediaitem_id=m.id WHERE a.user_id=$1 AND a.is_hidden=false AND a.is_hidden=false AND a.id IN (SELECT album_id FROM album_mediaitems WHERE mediaitem_id=$2) GROUP BY a.id, m.id ORDER BY a.created_at DESC`
	queryGetMediaItem               = `SELECT * FROM mediaitems WHERE user_id=$1 AND id=$2`
	queryUpdateMediaItem            = `UPDATE mediaitems SET description=$3, is_favourite=$4, is_hidden=$5, updated_at=$6 WHERE user_id=$1 AND id=$2`
	queryDeleteMediaItem            = `UPDATE mediaitems SET is_deleted=true, updated_at=$3 WHERE user_id=$1 AND id=$2`
	queryGetMediaItems              = `SELECT * FROM mediaitems WHERE user_id=$1 AND is_hidden=false AND is_deleted=false %s ORDER BY created_at DESC OFFSET $2 LIMIT $3`
	queryGetPlaceNewCoverMediaItem  = `SELECT DISTINCT ON (pm.place_id) pm.place_id, pm.mediaitem_id FROM place_mediaitems pm JOIN places p ON p.id = pm.place_id WHERE p.user_id=$1 AND p.cover_mediaitem_id=$2 AND pm.mediaitem_id!=$2`
	queryGetThingNewCoverMediaItem  = `SELECT DISTINCT ON (tm.thing_id) tm.thing_id, tm.mediaitem_id FROM thing_mediaitems tm JOIN things t ON t.id = tm.thing_id WHERE t.user_id=$1 AND t.cover_mediaitem_id=$2 AND tm.mediaitem_id!=$2`
	queryGetPeopleNewCoverMediaItem = `SELECT DISTINCT ON (pm.people_id) pm.people_id, pm.mediaitem_id FROM people_mediaitems pm JOIN people p ON p.id = pm.people_id WHERE p.user_id=$1 AND p.cover_mediaitem_id=$2 AND pm.mediaitem_id!=$2`
	queryGetAlbumNewCoverMediaItem  = `SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems am JOIN albums a ON a.id = am.album_id WHERE a.user_id=$1 AND a.cover_mediaitem_id=$2 AND am.mediaitem_id!=$2`
	queryUpdatePlaceCoverMediaItem  = `UPDATE places SET cover_mediaitem_id = $3 WHERE user_id = $1 AND id = $2`
	queryUpdateThingCoverMediaItem  = `UPDATE things SET cover_mediaitem_id = $3 WHERE user_id = $1 AND id = $2`
	queryUpdateAlbumCoverMediaItem  = `UPDATE albums SET cover_mediaitem_id = $3 WHERE user_id = $1 AND id = $2`
	queryUpdatePeopleCoverMediaItem = `UPDATE people SET cover_mediaitem_id = $3 WHERE user_id = $1 AND id = $2`
	queryInsertMediaItem            = `INSERT INTO mediaitems (id, user_id, filename, mediaitem_type, mediaitem_category, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	queryUpdateMediaItemHash        = `UPDATE mediaitems SET hash=$3, updated_at=$4 WHERE user_id=$1 AND id=$2`
)

type (
	// MediaItemRequest ...
	MediaItemRequest struct {
		Description *string `json:"description"`
		IsFavourite *bool   `json:"favourite"`
		IsHidden    *bool   `json:"hidden"`
		IsDeleted   *bool   `json:"deleted"`
	}

	// MediaItemResponse ...
	MediaItemResponse struct {
		ID string `json:"id"`
	}
)

// GetMediaItemPlaces ...
func (h *Handler) GetMediaItemPlaces(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	places := []models.Place{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetMediaItemPlaces, userID, uid)
	if err != nil {
		slog.Error("error getting mediaitem places", "error", err)
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
			slog.Error("error scanning mediaitem place", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		places = append(places, place)
	}
	return ctx.JSON(http.StatusOK, places)
}

// GetMediaItemThings ...
func (h *Handler) GetMediaItemThings(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	things := []models.Thing{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetMediaItemThings, userID, uid)
	if err != nil {
		slog.Error("error getting mediaitem things", "error", err)
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
			slog.Error("error scanning mediaitem thing", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		things = append(things, thing)
	}
	return ctx.JSON(http.StatusOK, things)
}

// GetMediaItemPeople ...
func (h *Handler) GetMediaItemPeople(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	people := []models.People{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetMediaItemPeople, userID, uid)
	if err != nil {
		slog.Error("error getting mediaitem people", "error", err)
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
			slog.Error("error scanning mediaitem person", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		people = append(people, person)
	}
	return ctx.JSON(http.StatusOK, people)
}

// GetMediaItemAlbums ...
func (h *Handler) GetMediaItemAlbums(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	albums := []models.Album{}
	rows, err := h.DB.Query(ctx.Request().Context(), queryGetMediaItemAlbums, userID, uid)
	if err != nil {
		slog.Error("error getting mediaitem albums", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		album := models.Album{CoverMediaItem: &models.MediaItem{}}
		err := rows.Scan(
			&album.ID,
			&album.UserID,
			&album.Name,
			&album.Description,
			&album.IsShared,
			&album.IsHidden,
			&album.MediaItemsCount,
			&album.CoverMediaItemID,
			&album.CreatedAt,
			&album.UpdatedAt,
			&album.CoverMediaItem.ID,
			&album.CoverMediaItem.UserID,
			&album.CoverMediaItem.Filename,
			&album.CoverMediaItem.Hash,
			&album.CoverMediaItem.Description,
			&album.CoverMediaItem.MimeType,
			&album.CoverMediaItem.SourceURL,
			&album.CoverMediaItem.PreviewURL,
			&album.CoverMediaItem.ThumbnailURL,
			&album.CoverMediaItem.Placeholder,
			&album.CoverMediaItem.IsFavourite,
			&album.CoverMediaItem.IsHidden,
			&album.CoverMediaItem.IsDeleted,
			&album.CoverMediaItem.Status,
			&album.CoverMediaItem.MediaItemType,
			&album.CoverMediaItem.MediaItemCategory,
			&album.CoverMediaItem.Width,
			&album.CoverMediaItem.Height,
			&album.CoverMediaItem.CreationTime,
			&album.CoverMediaItem.CameraMake,
			&album.CoverMediaItem.CameraModel,
			&album.CoverMediaItem.FocalLength,
			&album.CoverMediaItem.ApertureFnumber,
			&album.CoverMediaItem.IsoEquivalent,
			&album.CoverMediaItem.ExposureTime,
			&album.CoverMediaItem.Latitude,
			&album.CoverMediaItem.Longitude,
			&album.CoverMediaItem.FPS,
			&album.CoverMediaItem.EXIFData,
			&album.CoverMediaItem.Keywords,
			&album.CoverMediaItem.CreatedAt,
			&album.CoverMediaItem.UpdatedAt,
		)
		if err != nil {
			slog.Error("error scanning mediaitem album", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		albums = append(albums, album)
	}
	return ctx.JSON(http.StatusOK, albums)
}

// GetMediaItem ...
func (h *Handler) GetMediaItem(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := models.MediaItem{}
	err = h.DB.QueryRow(ctx.Request().Context(), queryGetMediaItem, userID, uid).Scan(&mediaItem.ID,
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
		slog.Error("error getting mediaitem", "error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "mediaitem not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItem)
}

// UpdateMediaItem ...
func (h *Handler) UpdateMediaItem(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	uid, err := getMediaItemID(ctx)
	if err != nil {
		return err
	}
	mediaItem, err := getMediaItem(ctx)
	if err != nil {
		return err
	}
	mediaItem.ID = uid
	mediaItem.UserID = userID
	mediaItem.UpdatedAt = time.Now()
	result, err := h.DB.Exec(ctx.Request().Context(), queryUpdateMediaItem, mediaItem.UserID, mediaItem.ID, mediaItem.Description, mediaItem.IsFavourite, mediaItem.IsHidden, mediaItem.UpdatedAt)
	if !result.Update() || err != nil {
		slog.Error("error updating mediaItem", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// DeleteMediaItem ...
func (h *Handler) DeleteMediaItem(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	uid, err := getMediaItemID(ctx)
	if err != nil {
		return err
	}
	mediaItem := models.MediaItem{ID: uid, UserID: userID, UpdatedAt: time.Now()}
	result, err := h.DB.Exec(ctx.Request().Context(), queryDeleteMediaItem, mediaItem.UserID, mediaItem.ID, mediaItem.UpdatedAt)
	if !result.Update() || err != nil {
		slog.Error("error deleting mediaItem", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err = h.updateCoverMediaItems(ctx.Request().Context(), userID, uid)
	if err != nil {
		slog.Error("error updating associated cover mediaitems", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

func (h *Handler) updateCoverMediaItems(ctx context.Context, userID, mediaItemID uuid.UUID) error { //nolint: funlen,cyclop
	var (
		albumsToUpdate [][]uuid.UUID
		placesToUpdate [][]uuid.UUID
		thingsToUpdate [][]uuid.UUID
		peopleToUpdate [][]uuid.UUID
		err            error
	)
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction for updating cover mediaitem: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()
	albumsToUpdate, err = h.getEntityToUpdate(ctx, tx, "album", queryGetAlbumNewCoverMediaItem, userID, mediaItemID)
	if err != nil {
		return err
	}
	placesToUpdate, err = h.getEntityToUpdate(ctx, tx, "place", queryGetPlaceNewCoverMediaItem, userID, mediaItemID)
	if err != nil {
		return err
	}
	thingsToUpdate, err = h.getEntityToUpdate(ctx, tx, "thing", queryGetThingNewCoverMediaItem, userID, mediaItemID)
	if err != nil {
		return err
	}
	peopleToUpdate, err = h.getEntityToUpdate(ctx, tx, "people", queryGetPeopleNewCoverMediaItem, userID, mediaItemID)
	if err != nil {
		return err
	}
	for _, albumToUpdate := range albumsToUpdate {
		_, err = tx.Exec(ctx, queryUpdateAlbumCoverMediaItem, userID, albumToUpdate[0], albumToUpdate[1])
		if err != nil {
			return fmt.Errorf("error updating album cover mediaitem: %w", err)
		}
	}
	for _, placeToUpdate := range placesToUpdate {
		_, err = tx.Exec(ctx, queryUpdatePlaceCoverMediaItem, userID, placeToUpdate[0], placeToUpdate[1])
		if err != nil {
			return fmt.Errorf("error updating place cover mediaitem: %w", err)
		}
	}
	for _, thingToUpdate := range thingsToUpdate {
		_, err = tx.Exec(ctx, queryUpdateThingCoverMediaItem, userID, thingToUpdate[0], thingToUpdate[1])
		if err != nil {
			return fmt.Errorf("error updating thing cover mediaitem: %w", err)
		}
	}
	for _, personToUpdate := range peopleToUpdate {
		_, err = tx.Exec(ctx, queryUpdatePeopleCoverMediaItem, userID, personToUpdate[0], personToUpdate[1])
		if err != nil {
			return fmt.Errorf("error updating people cover mediaitem: %w", err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction for updating cover mediaitem: %w", err)
	}
	return nil
}

func (h *Handler) getEntityToUpdate(ctx context.Context, tx pgx.Tx, entity, query string, userID, mediaItemID uuid.UUID) ([][]uuid.UUID, error) {
	rows, err := tx.Query(ctx, query, userID, mediaItemID)
	if err != nil {
		return nil, fmt.Errorf("error getting %s new cover mediaitems: %w", entity, err)
	}
	defer rows.Close()
	var entitiesToUpdate [][]uuid.UUID
	for rows.Next() {
		entityToUpdate := make([]uuid.UUID, 2)
		err = rows.Scan(&entityToUpdate[0], &entityToUpdate[1])
		if err != nil {
			return nil, fmt.Errorf("error scanning %s new cover mediaitem: %w", entity, err)
		}
		entitiesToUpdate = append(entitiesToUpdate, entityToUpdate)
	}
	return entitiesToUpdate, nil
}

// GetMediaItems ...
func (h *Handler) GetMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	filters := getMediaItemFilters(ctx)
	mediaItems := []models.MediaItem{}
	rows, err := h.DB.Query(ctx.Request().Context(), fmt.Sprintf(queryGetMediaItems, filters), userID, offset, limit)
	if err != nil {
		slog.Error("error getting mediaitems", "error", err)
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
			slog.Error("error scanning mediaitems", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		mediaItems = append(mediaItems, mediaItem)
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// UploadMediaItems ...
func (h *Handler) UploadMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	command := "start, finish"
	session := ""
	var err error
	uploadType := ctx.Request().Header.Get(HeaderUploadType)
	if uploadType == "resumable" {
		command, session, err = validateChunk(ctx)
		if err != nil {
			return err
		}
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		slog.Error("error uploading mediaitem", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	openedFile, err := file.Open()
	if err != nil {
		slog.Error("error reading uploaded mediaitem", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer openedFile.Close()

	if strings.Contains(command, "start") {
		mediaItem := createNewMediaItem(userID, file.Filename)
		result, err := h.DB.Exec(ctx.Request().Context(), queryInsertMediaItem,
			mediaItem.ID, mediaItem.UserID, mediaItem.Filename,
			mediaItem.MediaItemType, mediaItem.MediaItemCategory, mediaItem.Status,
			mediaItem.CreatedAt, mediaItem.UpdatedAt)
		if !result.Insert() || err != nil {
			slog.Error("error inserting mediaitem", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = h.saveToDisk(ctx.Request().Context(), userID.String(), mediaItem.ID.String(),
			openedFile, strings.Contains(command, "finish"))
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusCreated, &MediaItemResponse{
			ID: mediaItem.ID.String(),
		})
	}

	err = h.saveToDisk(ctx.Request().Context(), userID.String(), session,
		openedFile, strings.Contains(command, "finish"))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

func (h *Handler) saveToDisk(ctx context.Context, userID, mediaItemID string, openedFile multipart.File, finish bool) error {
	dstFile, err := os.OpenFile(fmt.Sprintf("%s/%s", h.Config.DiskRoot, mediaItemID), fileFlag, filePermission)
	if err != nil {
		slog.Error("error opening file", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	_, err = io.Copy(dstFile, openedFile)
	if err != nil {
		slog.Error("error copying file", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if finish {
		err = h.generateHashForDuplicates(ctx, userID, mediaItemID, dstFile.Name())
		if err != nil {
			if strings.Contains(err.Error(), "violates unique constraint") {
				slog.Error("error due to duplicate mediaitem", "error", err)
				return echo.NewHTTPError(http.StatusConflict, "mediaitem already exists")
			}
			slog.Error("error while generating hash for mediaitem", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return nil
}

func (h *Handler) generateHashForDuplicates(ctx context.Context, userID, mediaItemID, filePath string) error {
	openedFile, err := os.Open(filePath)
	if err != nil {
		slog.Error("error opening file for generating hash", "error", err)
		return err
	}
	defer openedFile.Close()

	fileHash := sha256.New()
	if _, err := io.Copy(fileHash, openedFile); err != nil {
		slog.Error("error copying file for generating hash", "error", err)
		return err
	}

	mediaItemHash := hex.EncodeToString(fileHash.Sum(nil))

	mediaItem := new(models.MediaItem)
	mediaItem.ID = uuid.FromStringOrNil(mediaItemID)
	mediaItem.UserID = uuid.FromStringOrNil(userID)
	mediaItem.Hash = &mediaItemHash
	mediaItem.UpdatedAt = time.Now()
	result, err := h.DB.Exec(ctx, queryUpdateMediaItemHash, mediaItem.UserID, mediaItem.ID, mediaItem.Hash, mediaItem.UpdatedAt)
	if !result.Update() || err != nil {
		slog.Error("error updating mediaitem hash", "error", err)
		return err
	}

	return nil
}

func createNewMediaItem(userID uuid.UUID, fileName string) *models.MediaItem {
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uuid.NewV4()
	mediaItem.UserID = userID
	mediaItem.Filename = fileName
	mediaItem.MediaItemType = models.TypeUnknown
	mediaItem.MediaItemCategory = models.CategoryDefault
	mediaItem.Status = models.StatusUnspecified
	mediaItem.CreatedAt = time.Now()
	mediaItem.UpdatedAt = mediaItem.CreatedAt
	return mediaItem
}

func validateChunk(ctx echo.Context) (string, string, error) {
	command := ctx.Request().Header.Get(HeaderUploadCommand)
	offset, _ := strconv.Atoi(ctx.Request().Header.Get(HeaderUploadChunkOffset))

	if len(command) == 0 {
		slog.Error("error getting command for resumable upload")
		return "", "", echo.NewHTTPError(http.StatusBadRequest, "invalid command for resumable upload")
	}
	if command != "start" && offset == 0 {
		slog.Error("error getting chunk offset for resumable upload")
		return "", "", echo.NewHTTPError(http.StatusBadRequest, "invalid chunk offset for resumable upload")
	}
	session := ctx.Request().Header.Get(HeaderUploadChunkSession)
	if command != "start" && len(session) == 0 {
		slog.Error("error getting chunk session for resumable upload")
		return "", "", echo.NewHTTPError(http.StatusBadRequest, "invalid chunk session for resumable upload")
	}

	return command, session, nil
}

func getMediaItemID(ctx echo.Context) (uuid.UUID, error) {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		slog.Error("error getting mediaitem id", "error", err)
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	return uid, err
}

func getMediaItem(ctx echo.Context) (*models.MediaItem, error) {
	mediaItemRequest := new(MediaItemRequest)
	err := ctx.Bind(mediaItemRequest)
	if err != nil {
		slog.Error("error getting mediaitem", "error", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem")
	}
	mediaItem := models.MediaItem{
		Description: mediaItemRequest.Description,
		IsFavourite: mediaItemRequest.IsFavourite,
		IsHidden:    mediaItemRequest.IsHidden,
		IsDeleted:   mediaItemRequest.IsDeleted,
	}
	if reflect.DeepEqual(models.MediaItem{}, mediaItem) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem")
	}
	return &mediaItem, nil
}

package handlers

import (
	"api/internal/models"
	"api/pkg/services/api"
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

	queryGetMediaItemPlaces = `SELECT p.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,` +
		` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM places p LEFT JOIN mediaitems m ON p.cover_mediaitem_id=m.id` +
		` WHERE p.user_id=$1 AND p.is_hidden=false AND p.id IN (SELECT place_id FROM place_mediaitems` +
		` WHERE mediaitem_id=$2) ORDER BY p.created_at DESC`
	queryGetMediaItemThings = `SELECT t.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,` +
		` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM things t LEFT JOIN mediaitems m ON t.cover_mediaitem_id=m.id` +
		` WHERE t.user_id=$1 AND t.is_hidden=false AND t.id IN (SELECT thing_id FROM thing_mediaitems` +
		` WHERE mediaitem_id=$2) ORDER BY t.created_at DESC`
	queryGetMediaItemPeople = `SELECT p.*, mf.* FROM people p LEFT JOIN mediaitem_faces mf ON p.cover_mediaitem_face_id=mf.id` +
		` WHERE p.user_id=$1 AND p.is_hidden=false AND p.id IN (SELECT people_id FROM people_mediaitems` +
		` WHERE mediaitem_id=$2) ORDER BY p.created_at DESC`
	queryGetMediaItemAlbums = `SELECT a.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,` +
		` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM albums a LEFT JOIN mediaitems m ON a.cover_mediaitem_id=m.id` +
		` WHERE a.user_id=$1 AND a.is_hidden=false AND a.is_hidden=false AND a.id IN (SELECT album_id FROM album_mediaitems` +
		` WHERE mediaitem_id=$2) ORDER BY a.created_at DESC`
	queryGetMediaItem    = `SELECT * FROM mediaitems WHERE user_id=$1 AND id=$2`
	queryUpdateMediaItem = `UPDATE mediaitems SET description=$3, is_favourite=$4, is_hidden=$5,` +
		` updated_at=$6 WHERE user_id=$1 AND id=$2`
	queryDeleteMediaItem = `UPDATE mediaitems SET is_deleted=true, updated_at=$3 WHERE user_id=$1 AND id=$2`
	queryGetMediaItems   = `SELECT * FROM mediaitems WHERE user_id=$1 AND is_hidden=false AND is_deleted=false` +
		` %s ORDER BY created_at DESC OFFSET $2 LIMIT $3`
	queryGetPlaceNewCoverMediaItem = `SELECT DISTINCT ON (pm.place_id) pm.place_id, pm.mediaitem_id FROM place_mediaitems pm` +
		` JOIN places p ON p.id = pm.place_id WHERE p.user_id=$1 AND p.cover_mediaitem_id=$2 AND pm.mediaitem_id!=$2`
	queryGetThingNewCoverMediaItem = `SELECT DISTINCT ON (tm.thing_id) tm.thing_id, tm.mediaitem_id FROM thing_mediaitems tm` +
		` JOIN things t ON t.id = tm.thing_id WHERE t.user_id=$1 AND t.cover_mediaitem_id=$2 AND tm.mediaitem_id!=$2`
	queryGetPeopleNewCoverMediaItem = `SELECT DISTINCT ON (pm.people_id) pm.people_id, pm.mediaitem_id FROM people_mediaitems pm` +
		` JOIN people p ON p.id = pm.people_id WHERE p.user_id=$1 AND p.cover_mediaitem_id=$2 AND pm.mediaitem_id!=$2`
	queryGetAlbumNewCoverMediaItem = `SELECT DISTINCT ON (am.album_id) am.album_id, am.mediaitem_id FROM album_mediaitems am` +
		` JOIN albums a ON a.id = am.album_id WHERE a.user_id=$1 AND a.cover_mediaitem_id=$2 AND am.mediaitem_id!=$2`
	queryUpdatePlaceCoverMediaItem  = `UPDATE places SET cover_mediaitem_id = $3 WHERE user_id = $1 AND id = $2`
	queryUpdateThingCoverMediaItem  = `UPDATE things SET cover_mediaitem_id = $3 WHERE user_id = $1 AND id = $2`
	queryUpdateAlbumCoverMediaItem  = `UPDATE albums SET cover_mediaitem_id = $3 WHERE user_id = $1 AND id = $2`
	queryUpdatePeopleCoverMediaItem = `UPDATE people SET cover_mediaitem_id = $3 WHERE user_id = $1 AND id = $2`
	queryInsertMediaItem            = `INSERT INTO mediaitems (id, user_id, filename, mediaitem_type, mediaitem_category,` +
		` status, source_url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	queryQueueMediaItem      = `INSERT INTO queue(id, user_id, mediaitem_id, components, status) VALUES (gen_random_uuid(), $1, $2, $3, $4)`
	queryUpdateMediaItemHash = `UPDATE mediaitems SET hash=$3 WHERE user_id=$1 AND id=$2`
)

type ( // MediaItemRequest ...
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
		place, err := models.ScanRowsToPlace(rows)
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
		thing, err := models.ScanRowsToThing(rows)
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
		person, err := models.ScanRowsToPerson(rows)
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
		album, err := models.ScanRowsToAlbum(rows)
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
		&mediaItem.UserID, &mediaItem.Filename, &mediaItem.Hash, &mediaItem.Description, &mediaItem.MimeType,
		&mediaItem.SourceURL, &mediaItem.PreviewURL, &mediaItem.ThumbnailURL, &mediaItem.Placeholder,
		&mediaItem.IsFavourite, &mediaItem.IsHidden, &mediaItem.IsDeleted, &mediaItem.Status, &mediaItem.MediaItemType,
		&mediaItem.MediaItemCategory, &mediaItem.Width, &mediaItem.Height, &mediaItem.CreationTime,
		&mediaItem.CameraMake, &mediaItem.CameraModel, &mediaItem.FocalLength, &mediaItem.ApertureFnumber,
		&mediaItem.IsoEquivalent, &mediaItem.ExposureTime, &mediaItem.Megapixels, &mediaItem.Latitude,
		&mediaItem.Longitude, &mediaItem.FPS, &mediaItem.EXIFData, &mediaItem.Keywords,
		&mediaItem.CreatedAt, &mediaItem.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "mediaitem not found")
		}
		slog.Error("error getting mediaitem", "error", err)

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
	_, err = h.DB.Exec(ctx.Request().Context(), queryUpdateMediaItem, mediaItem.UserID,
		mediaItem.ID, mediaItem.Description, mediaItem.IsFavourite, mediaItem.IsHidden, mediaItem.UpdatedAt)
	if err != nil {
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
	mediaItem := models.MediaItem{
		ID: uid, UserID: userID, UpdatedAt: time.Now(),
	}
	_, err = h.DB.Exec(ctx.Request().Context(), queryDeleteMediaItem, mediaItem.UserID,
		mediaItem.ID, mediaItem.UpdatedAt)
	if err != nil {
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

func (h *Handler) updateCoverMediaItems(ctx context.Context, userID, mediaItemID uuid.UUID) error {
	var (
		entities         = []string{"album", "place", "thing", "people"}
		entityGetQueries = []string{
			queryGetAlbumNewCoverMediaItem, queryGetPlaceNewCoverMediaItem, queryGetThingNewCoverMediaItem,
			queryGetPeopleNewCoverMediaItem,
		}
		entityUpdateQueries = []string{
			queryUpdateAlbumCoverMediaItem, queryUpdatePlaceCoverMediaItem, queryUpdateThingCoverMediaItem,
			queryUpdatePeopleCoverMediaItem,
		}
		err error
	)
	entityItemsToUpdate := make([][][]uuid.UUID, len(entities))
	mtx, err := h.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction for updating cover mediaitem: %w", err)
	}
	defer func() {
		if err != nil {
			_ = mtx.Rollback(ctx)
		}
	}()
	for idx, entity := range entities {
		entityItemsToUpdate[idx], err = h.getEntityToUpdate(ctx, mtx, entity, entityGetQueries[idx], userID, mediaItemID)
		if err != nil {
			return err
		}
	}
	for idx, entityItemToUpdate := range entityItemsToUpdate {
		for _, itemToUpdate := range entityItemToUpdate {
			_, err = mtx.Exec(ctx, entityUpdateQueries[idx], userID, itemToUpdate[0], itemToUpdate[1])
			if err != nil {
				return fmt.Errorf("error updating %s cover mediaitem: %w", entities[idx], err)
			}
		}
	}
	if err = mtx.Commit(ctx); err != nil {
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
		entityToUpdate := make([]uuid.UUID, 2) //nolint: mnd
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
		mediaItem, err := models.ScanRowsToMediaItem(rows)
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

	features, _ := ctx.Get("features").(models.Features)

	if strings.Contains(command, "start") {
		mediaItem := createNewMediaItem(userID, file.Filename)
		mediaItem.SourceURL = fmt.Sprintf("%s/%s", h.Config.DiskRoot, mediaItem.ID)
		_, err = h.DB.Exec(ctx.Request().Context(), queryInsertMediaItem, mediaItem.ID, mediaItem.UserID,
			mediaItem.Filename, mediaItem.MediaItemType, mediaItem.MediaItemCategory,
			mediaItem.Status, mediaItem.SourceURL, mediaItem.CreatedAt, mediaItem.UpdatedAt)
		if err != nil {
			slog.Error("error inserting mediaitem", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = h.saveToDisk(ctx.Request().Context(), userID.String(), mediaItem.ID.String(), features, openedFile, strings.Contains(command, "finish"))
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusCreated, &MediaItemResponse{
			ID: mediaItem.ID.String(),
		})
	}

	err = h.saveToDisk(ctx.Request().Context(), userID.String(), session, features, openedFile, strings.Contains(command, "finish"))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

func (h *Handler) saveToDisk(ctx context.Context, userID, mediaItemID string, features models.Features, openedFile multipart.File, finish bool) error {
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

		err = h.queueMediaItemForProcessing(ctx, userID, mediaItemID, features)
		if err != nil {
			slog.Error("error queuing mediaitem for processing", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return nil
}

func (h *Handler) queueMediaItemForProcessing(ctx context.Context, userID, mediaItemID string, features models.Features) error { //nolint: cyclop
	components := fmt.Sprintf("%s,%s", api.MediaItemComponent_METADATA.String(), api.MediaItemComponent_PREVIEW_THUMBNAIL.String())
	if h.Config.ML.Places && features.Places {
		components += ("," + api.MediaItemComponent_PLACES.String())
	}
	if h.Config.Classification && features.Things {
		components += ("," + api.MediaItemComponent_CLASSIFICATION.String())
	}
	if h.Config.Faces && features.People {
		components += ("," + api.MediaItemComponent_FACES.String())
	}
	if h.Config.OCR && features.Explore {
		components += ("," + api.MediaItemComponent_OCR.String())
	}
	if h.Config.Search && features.Explore {
		components += ("," + api.MediaItemComponent_SEARCH.String())
	}
	_, err := h.DB.Exec(ctx, queryQueueMediaItem, userID, mediaItemID, components, api.MediaItemStatus_UNSPECIFIED)
	if err != nil {
		return err
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
	_, err = h.DB.Exec(ctx, queryUpdateMediaItemHash, mediaItem.UserID, mediaItem.ID, mediaItem.Hash)
	if err != nil {
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
	mediaItem.MediaItemType = api.MediaItemType_UNKNOWN.String()
	mediaItem.MediaItemCategory = api.MediaItemCategory_DEFAULT.String()
	mediaItem.Status = api.MediaItemStatus_UNSPECIFIED.String()
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
		Description: mediaItemRequest.Description, IsFavourite: mediaItemRequest.IsFavourite,
		IsHidden: mediaItemRequest.IsHidden, IsDeleted: mediaItemRequest.IsDeleted,
	}
	if reflect.DeepEqual(models.MediaItem{}, mediaItem) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem")
	}

	return &mediaItem, nil
}

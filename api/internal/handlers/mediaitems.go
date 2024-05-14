package handlers

import (
	"api/internal/models"
	"api/pkg/services/worker"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

const (
	HeaderUploadType         = "X-Smriti-Upload-Type"    // for resumable
	HeaderUploadCommand      = "X-Smriti-Upload-Command" // start, continue, finish
	HeaderUploadChunkOffset  = "X-Smriti-Upload-Chunk-Offset"
	HeaderUploadChunkSession = "X-Smriti-Upload-Chunk-Session"

	fileFlag       = os.O_WRONLY | os.O_APPEND | os.O_CREATE
	filePermission = 0o644
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
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	mediaItem.UserID = userID
	places := []models.Place{}
	err = h.DB.Model(&mediaItem).Preload("CoverMediaItem").Association("Places").Find(&places)
	if err != nil {
		slog.Error("error getting mediaitem places", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	mediaItem.UserID = userID
	things := []models.Thing{}
	err = h.DB.Model(&mediaItem).Preload("CoverMediaItem").Association("Things").Find(&things)
	if err != nil {
		slog.Error("error getting mediaitem things", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	mediaItem.UserID = userID
	people := []models.People{}
	err = h.DB.Model(&mediaItem).Preload("CoverMediaItemFace").Association("People").Find(&people)
	if err != nil {
		slog.Error("error getting mediaitem people", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	mediaItem.UserID = userID
	albums := []models.Album{}
	err = h.DB.Model(&mediaItem).Preload("CoverMediaItem").Association("Albums").Find(&albums)
	if err != nil {
		slog.Error("error getting mediaitem albums", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	result := h.DB.Where("id=? AND user_id=?", uid, userID).First(&mediaItem)
	if result.Error != nil {
		slog.Error("error getting mediaitem", "error", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "mediaitem not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
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
	result := h.DB.Model(&mediaItem).Updates(mediaItem)
	if result.Error != nil {
		slog.Error("error updating mediaItem", "error", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
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
	deleted := true
	mediaItem := models.MediaItem{ID: uid, UserID: userID, IsDeleted: &deleted}
	result := h.DB.Model(&mediaItem).Updates(mediaItem)
	if result.Error != nil {
		slog.Error("error updating mediaItem", "error", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	// album
	err = h.updateCoverMediaItems(uid)
	if err != nil {
		slog.Error("error updating associated cover mediaitems", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

func (h *Handler) updateCoverMediaItems(mediaItemID uuid.UUID) error { //nolint: funlen,cyclop
	var (
		albumsToUpdate []models.Album
		placesToUpdate []models.Place
		thingsToUpdate []models.Thing
		peopleToUpdate []models.People
	)
	result := h.DB.Model(&models.Album{}).Preload("MediaItems").Where("cover_mediaitem_id = ?", mediaItemID).Find(&albumsToUpdate)
	if result.Error != nil {
		return fmt.Errorf("error getting albums: %w", result.Error)
	}
	for _, album := range albumsToUpdate {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(album.MediaItems))))
		var newCoverMediaItemID *uuid.UUID
		newCoverMediaItemID = &album.MediaItems[randomIndex.Int64()].ID
		if len(album.MediaItems) == 1 {
			newCoverMediaItemID = nil
		}
		result := h.DB.Model(&models.Album{UserID: album.UserID, ID: album.ID}).Omit("MediaItems").Updates(map[string]interface{}{
			"CoverMediaItemID": newCoverMediaItemID,
		})
		if result.Error != nil {
			return fmt.Errorf("error updating album cover mediaitem: %w", result.Error)
		}
	}
	result = h.DB.Model(&models.Place{}).Preload("MediaItems").Where("cover_mediaitem_id = ?", mediaItemID).Find(&placesToUpdate)
	if result.Error != nil {
		return fmt.Errorf("error getting places: %w", result.Error)
	}
	for _, place := range placesToUpdate {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(place.MediaItems))))
		var newCoverMediaItemID *uuid.UUID
		newCoverMediaItemID = &place.MediaItems[randomIndex.Int64()].ID
		if len(place.MediaItems) == 1 {
			newCoverMediaItemID = nil
		}
		result := h.DB.Model(&models.Place{UserID: place.UserID, ID: place.ID}).Omit("MediaItems").Updates(map[string]interface{}{
			"CoverMediaItemID": newCoverMediaItemID,
		})
		if result.Error != nil {
			return fmt.Errorf("error updating place cover mediaitem: %w", result.Error)
		}
	}
	result = h.DB.Model(&models.Thing{}).Preload("MediaItems").Where("cover_mediaitem_id = ?", mediaItemID).Find(&thingsToUpdate)
	if result.Error != nil {
		return fmt.Errorf("error getting things: %w", result.Error)
	}
	for _, thing := range thingsToUpdate {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(thing.MediaItems))))
		var newCoverMediaItemID *uuid.UUID
		newCoverMediaItemID = &thing.MediaItems[randomIndex.Int64()].ID
		if len(thing.MediaItems) == 1 {
			newCoverMediaItemID = nil
		}
		result := h.DB.Model(&models.Thing{UserID: thing.UserID, ID: thing.ID}).Omit("MediaItems").Updates(map[string]interface{}{
			"CoverMediaItemID": newCoverMediaItemID,
		})
		if result.Error != nil {
			return fmt.Errorf("error updating thing cover mediaitem: %w", result.Error)
		}
	}
	result = h.DB.Model(&models.People{}).Preload("MediaItems").Where("cover_mediaitem_id = ?", mediaItemID).Find(&peopleToUpdate)
	if result.Error != nil {
		return fmt.Errorf("error getting people: %w", result.Error)
	}
	for _, person := range peopleToUpdate {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(person.MediaItems))))
		var newCoverMediaItemID *uuid.UUID
		newCoverMediaItemID = &person.MediaItems[randomIndex.Int64()].ID
		if len(person.MediaItems) == 1 {
			newCoverMediaItemID = nil
		}
		result := h.DB.Model(&models.People{UserID: person.UserID, ID: person.ID}).Omit("MediaItems").Updates(map[string]interface{}{
			"CoverMediaItemID": newCoverMediaItemID,
		})
		if result.Error != nil {
			return fmt.Errorf("error updating people cover mediaitem: %w", result.Error)
		}
	}
	return nil
}

// GetMediaItems ...
func (h *Handler) GetMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	filters := getMediaItemFilters(ctx)
	mediaItems := []models.MediaItem{}
	result := h.DB.Where("user_id=? AND is_hidden=false AND is_deleted=false"+filters, userID).
		Find(&mediaItems).
		Offset(offset).
		Limit(limit)
	if result.Error != nil {
		slog.Error("error getting mediaitems", "error", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// UploadMediaItems ...
func (h *Handler) UploadMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	features, _ := ctx.Get("features").(models.Features)
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

	components := []string{}
	if strings.Contains(command, "finish") {
		components = h.getComponents(features)
	}

	if strings.Contains(command, "start") {
		mediaItem := createNewMediaItem(userID, file.Filename)
		result := h.DB.Create(&mediaItem)
		if result.Error != nil {
			slog.Error("error inserting mediaitem", "error", result.Error)
			return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
		}

		err = h.saveToDiskAndSendToWorker(userID.String(), mediaItem.ID.String(),
			openedFile, components)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusCreated, &MediaItemResponse{
			ID: mediaItem.ID.String(),
		})
	}

	err = h.saveToDiskAndSendToWorker(userID.String(), session,
		openedFile, components)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

func (h *Handler) saveToDiskAndSendToWorker(userID, mediaItemID string, openedFile multipart.File, components []string) error {
	dstFile, err := os.OpenFile(fmt.Sprintf("%s/%s", h.Config.Storage.DiskRoot, mediaItemID), fileFlag, filePermission)
	if err != nil {
		slog.Error("error opening file", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	_, err = io.Copy(dstFile, openedFile)
	if err != nil {
		slog.Error("error copying file", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(components) != 0 {
		err = h.generateHashForDuplicates(userID, mediaItemID, dstFile.Name())
		if err != nil {
			if strings.Contains(err.Error(), "violates unique constraint") {
				slog.Error("error due to duplicate mediaitem", "error", err)
				return echo.NewHTTPError(http.StatusConflict, "mediaitem already exists")
			}
			slog.Error("error while generating hash for mediaitem", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		_, err = h.Worker.MediaItemProcess(context.Background(), &worker.MediaItemProcessRequest{
			UserId:     userID,
			Id:         mediaItemID,
			FilePath:   h.Config.Storage.DiskRoot,
			Components: components,
		})
		if err != nil {
			slog.Error("error sending mediaitem for processing", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return nil
}

func (h *Handler) generateHashForDuplicates(userID, mediaItemID, filePath string) error {
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
	result := h.DB.Model(&mediaItem).Updates(mediaItem)
	if result.Error != nil {
		slog.Error("error updating mediaitem hash", "error", result.Error)
		return result.Error
	}

	return nil
}

//nolint:cyclop
func (h *Handler) getComponents(features models.Features) []string {
	components := []string{"metadata"}
	if h.Config.ML.Places && features.Places {
		components = append(components, "places")
	}
	if h.Config.ML.Classification && features.Things {
		components = append(components, "classification")
	}
	if h.Config.ML.OCR && features.Explore {
		components = append(components, "ocr")
	}
	if h.Config.ML.Search && features.Explore {
		components = append(components, "search")
	}
	if h.Config.ML.Faces && features.People {
		components = append(components, "faces")
	}
	components = append(components, "finalize")
	return components
}

func createNewMediaItem(userID uuid.UUID, fileName string) *models.MediaItem {
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uuid.NewV4()
	mediaItem.UserID = userID
	mediaItem.Filename = fileName
	mediaItem.MediaItemType = models.Unknown
	mediaItem.MediaItemCategory = models.Default
	mediaItem.Status = models.Processing
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

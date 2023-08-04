package handlers

import (
	"api/internal/models"
	"api/pkg/services/worker"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
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

	fileFlag       = os.O_WRONLY | os.O_APPEND | os.O_CREATE //nolint: nosnakecase
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
		slog.Error("error getting mediaitem id", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	mediaItem.UserID = userID
	places := []models.Place{}
	err = h.DB.Model(&mediaItem).Preload("CoverMediaItem").Association("Places").Find(&places)
	if err != nil {
		slog.Error("error getting mediaitem places", slog.Any("error", err))
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
		slog.Error("error getting mediaitem id", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	mediaItem.UserID = userID
	things := []models.Thing{}
	err = h.DB.Model(&mediaItem).Preload("CoverMediaItem").Association("Things").Find(&things)
	if err != nil {
		slog.Error("error getting mediaitem things", slog.Any("error", err))
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
		slog.Error("error getting mediaitem id", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	mediaItem.UserID = userID
	people := []models.People{}
	err = h.DB.Model(&mediaItem).Preload("CoverMediaItem").Association("People").Find(&people)
	if err != nil {
		slog.Error("error getting mediaitem people", slog.Any("error", err))
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
		slog.Error("error getting mediaitem id", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	mediaItem.UserID = userID
	albums := []models.Album{}
	err = h.DB.Model(&mediaItem).Preload("CoverMediaItem").Association("Albums").Find(&albums)
	if err != nil {
		slog.Error("error getting mediaitem albums", slog.Any("error", err))
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
		slog.Error("error getting mediaitem id", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := models.MediaItem{}
	result := h.DB.Where("id=? AND user_id=?", uid, userID).First(&mediaItem)
	if result.Error != nil {
		slog.Error("error getting mediaitem", slog.Any("error", result.Error))
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
	if result.Error != nil || result.RowsAffected != 1 {
		slog.Error("error updating mediaItem", slog.Any("error", result.Error))
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
	if result.Error != nil || result.RowsAffected != 1 {
		slog.Error("error updating mediaItem", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	// work(omkar): update albums, things, people and places to choose another cover mediaitem
	// where this mediaitem is a cover mediaitem
	return ctx.JSON(http.StatusNoContent, nil)
}

// GetMediaItems ...
func (h *Handler) GetMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	filters := getMediaItemFilters(ctx)
	mediaItems := []models.MediaItem{}
	result := h.DB.Where(fmt.Sprintf("user_id=? AND is_hidden=false AND is_deleted=false%s", filters), userID).
		Find(&mediaItems).
		Offset(offset).
		Limit(limit)
	if result.Error != nil {
		slog.Error("error getting mediaitems", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
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
		slog.Error("error uploading mediaitem", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	openedFile, err := file.Open()
	if err != nil {
		slog.Error("error reading uploaded mediaitem", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer openedFile.Close()

	if strings.Contains(command, "start") {
		mediaItem := createNewMediaItem(userID, file.Filename)
		result := h.DB.Create(&mediaItem)
		if result.Error != nil {
			slog.Error("error inserting mediaitem", slog.Any("error", result.Error))
			return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
		}

		err = h.saveToDiskAndSendToWorker(userID.String(), mediaItem.ID.String(),
			openedFile, strings.Contains(command, "finish"))
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusCreated, &MediaItemResponse{
			ID: mediaItem.ID.String(),
		})
	}

	err = h.saveToDiskAndSendToWorker(userID.String(), session,
		openedFile, strings.Contains(command, "finish"))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

func (h *Handler) saveToDiskAndSendToWorker(userID, mediaItemID string, openedFile multipart.File, sendToWorker bool) error {
	dstFile, err := os.OpenFile(fmt.Sprintf("%s/%s", h.Config.Storage.DiskRoot, mediaItemID), fileFlag, filePermission)
	if err != nil {
		slog.Error("error opening file", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	_, err = io.Copy(dstFile, openedFile)
	if err != nil {
		slog.Error("error copying file", slog.Any("error", err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if sendToWorker {
		err = h.generateHashForDuplicates(userID, mediaItemID, dstFile.Name())
		if err != nil {
			if strings.Contains(err.Error(), "violates unique constraint") {
				slog.Error("error due to duplicate mediaitem", slog.Any("error", err))
				return echo.NewHTTPError(http.StatusConflict, "mediaitem already exists")
			}
			slog.Error("error while generating hash for mediaitem", slog.Any("error", err))
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		_, err = h.Worker.MediaItemProcess(context.Background(), &worker.MediaItemProcessRequest{
			UserId:   userID,
			Id:       mediaItemID,
			FilePath: h.Config.Storage.DiskRoot,
		})
		if err != nil {
			slog.Error("error sending mediaitem for processing", slog.Any("error", err))
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return nil
}

func (h *Handler) generateHashForDuplicates(userID, mediaItemID, filePath string) error {
	openedFile, err := os.Open(filePath)
	if err != nil {
		slog.Error("error opening file for generating hash", slog.Any("error", err))
		return err
	}
	defer openedFile.Close()

	fileHash := sha256.New()
	if _, err := io.Copy(fileHash, openedFile); err != nil {
		slog.Error("error copying file for generating hash", slog.Any("error", err))
		return err
	}

	mediaItemHash := fmt.Sprintf("%x", fileHash.Sum(nil))

	mediaItem := new(models.MediaItem)
	mediaItem.ID = uuid.FromStringOrNil(mediaItemID)
	mediaItem.UserID = uuid.FromStringOrNil(userID)
	mediaItem.Hash = &mediaItemHash
	result := h.DB.Model(&mediaItem).Updates(mediaItem)
	if result.Error != nil {
		slog.Error("error updating mediaitem hash", slog.Any("error", result.Error))
		return result.Error
	}

	return nil
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
		slog.Error("error getting mediaitem id", slog.Any("error", err))
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	return uid, err
}

func getMediaItem(ctx echo.Context) (*models.MediaItem, error) {
	mediaItemRequest := new(MediaItemRequest)
	err := ctx.Bind(mediaItemRequest)
	if err != nil {
		slog.Error("error getting mediaitem", slog.Any("error", err))
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

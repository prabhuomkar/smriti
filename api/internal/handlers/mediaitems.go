package handlers

import (
	"api/internal/models"
	"api/pkg/services/worker"
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const (
	HeaderUploadType         = "X-Smriti-Upload-Type"    // for resumable
	HeaderUploadCommand      = "X-Smriti-Upload-Command" // start, continue, finish
	HeaderUploadChunkOffset  = "X-Smriti-Upload-Chunk-Offset"
	HeaderUploadChunkSession = "X-Smriti-Upload-Chunk-Session"

	streamChunkByteSize = 1024
)

type (
	// MediaItemRequest ...
	MediaItemRequest struct {
		Description *string `json:"description"`
		IsFavourite *bool   `json:"favourite"`
		IsHidden    *bool   `json:"hidden"`
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
		log.Printf("error getting mediaitem id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	mediaItem.UserID = userID
	places := []models.Place{}
	err = h.DB.Model(&mediaItem).Preload("CoverMediaItem").Association("Places").Find(&places)
	if err != nil {
		log.Printf("error getting mediaitem places: %+v", err)
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
		log.Printf("error getting mediaitem id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	mediaItem.UserID = userID
	things := []models.Thing{}
	err = h.DB.Model(&mediaItem).Preload("CoverMediaItem").Association("Things").Find(&things)
	if err != nil {
		log.Printf("error getting mediaitem things: %+v", err)
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
		log.Printf("error getting mediaitem id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	mediaItem.UserID = userID
	people := []models.People{}
	err = h.DB.Model(&mediaItem).Preload("CoverMediaItem").Association("People").Find(&people)
	if err != nil {
		log.Printf("error getting mediaitem people: %+v", err)
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
		log.Printf("error getting mediaitem id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := new(models.MediaItem)
	mediaItem.ID = uid
	mediaItem.UserID = userID
	albums := []models.Album{}
	err = h.DB.Model(&mediaItem).Preload("CoverMediaItem").Association("Albums").Find(&albums)
	if err != nil {
		log.Printf("error getting mediaitem albums: %+v", err)
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
		log.Printf("error getting mediaitem id: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	mediaItem := models.MediaItem{}
	result := h.DB.Where("id=? AND user_id=?", uid, userID).First(&mediaItem)
	if result.Error != nil {
		log.Printf("error getting mediaitem: %+v", result.Error)
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
		log.Printf("error updating mediaItem: %+v", result.Error)
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
		log.Printf("error updating mediaItem: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// GetMediaItems ...
func (h *Handler) GetMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	offset, limit := getOffsetAndLimit(ctx)
	mediaItems := []models.MediaItem{}
	result := h.DB.Where("user_id=? AND is_hidden=false OR is_deleted=false", userID).
		Find(&mediaItems).
		Offset(offset).
		Limit(limit)
	if result.Error != nil {
		log.Printf("error getting mediaitems: %+v", result.Error)
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

// UploadMediaItems ...
func (h *Handler) UploadMediaItems(ctx echo.Context) error {
	userID := getRequestingUserID(ctx)
	command := "start, finish"
	offset := 0
	session := ""
	var err error
	uploadType := ctx.Request().Header.Get(HeaderUploadType)
	if uploadType == "resumable" {
		command, offset, session, err = validateChunk(ctx)
		if err != nil {
			return err
		}
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		log.Printf("error uploading mediaitem: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	openedFile, err := file.Open()
	if err != nil {
		log.Printf("error reading uploaded mediaitem: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer openedFile.Close()

	if strings.Contains(command, "start") {
		mediaItem := createNewMediaItem(userID, file.Filename)
		result := h.DB.Create(&mediaItem)
		if result.Error != nil {
			log.Printf("error inserting mediaitem: %+v", result.Error)
			return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
		}

		err = sendFileToWorker(h.Worker, userID.String(), mediaItem.ID.String(), command, offset, openedFile)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusCreated, &MediaItemResponse{
			ID: mediaItem.ID.String(),
		})
	}

	err = sendFileToWorker(h.Worker, userID.String(), session, command, offset, openedFile)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

func sendFileToWorker(workerClient worker.WorkerClient, userID, fileID, command string, offset int, file multipart.File) error { //nolint: lll
	stream, err := workerClient.MediaItemProcess(context.Background())
	if err != nil {
		log.Printf("error creating stream for sending mediaitem to worker: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, streamChunkByteSize)
	for {
		numBytes, err := reader.Read(buffer)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Printf("error reading uploaded file: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = stream.Send(&worker.MediaItemProcessRequest{
			UserId:  userID,
			Id:      fileID,
			Offset:  int64(offset),
			Command: command,
			Content: buffer[:numBytes],
		})
		if err != nil {
			log.Printf("error sending mediaitem to worker: %+v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("error receiving response from worker: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !res.Ok {
		log.Println("error due to bad response from worker")
		return echo.NewHTTPError(http.StatusInternalServerError, "error uploading mediaitem")
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

func validateChunk(ctx echo.Context) (string, int, string, error) {
	command := ctx.Request().Header.Get(HeaderUploadCommand)
	offset, _ := strconv.Atoi(ctx.Request().Header.Get(HeaderUploadChunkOffset))

	if len(command) == 0 {
		log.Printf("error getting command for resumable upload")
		return "", 0, "", echo.NewHTTPError(http.StatusBadRequest, "invalid command for resumable upload")
	}
	if command != "start" && offset == 0 {
		log.Printf("error getting chunk offset for resumable upload")
		return "", 0, "", echo.NewHTTPError(http.StatusBadRequest, "invalid chunk offset for resumable upload")
	}
	session := ctx.Request().Header.Get(HeaderUploadChunkSession)
	if command != "start" && len(session) == 0 {
		log.Printf("error getting chunk session for resumable upload")
		return "", 0, "", echo.NewHTTPError(http.StatusBadRequest, "invalid chunk session for resumable upload")
	}

	return command, offset, session, nil
}

func getMediaItemID(ctx echo.Context) (uuid.UUID, error) {
	id := ctx.Param("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Printf("error getting mediaitem id: %+v", err)
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem id")
	}
	return uid, err
}

func getMediaItem(ctx echo.Context) (*models.MediaItem, error) {
	mediaItemRequest := new(MediaItemRequest)
	err := ctx.Bind(mediaItemRequest)
	if err != nil {
		log.Printf("error getting mediaitem: %+v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem")
	}
	mediaItem := models.MediaItem{
		Description: mediaItemRequest.Description,
		IsFavourite: mediaItemRequest.IsFavourite,
		IsHidden:    mediaItemRequest.IsHidden,
	}
	if reflect.DeepEqual(models.MediaItem{}, mediaItem) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid mediaitem")
	}
	return &mediaItem, nil
}

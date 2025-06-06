package models

import (
	"api/pkg/cache"
	"api/pkg/storage"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const MediaItemsTable = "mediaitems"

type (
	// MediaItemStatus ...
	MediaItemStatus string

	// MediaItemType ...
	MediaItemType string

	// MediaItemCategory ...
	MediaItemCategory string

	// MediaitemEmbedding ...
	MediaitemEmbedding struct {
		MediaitemID uuid.UUID        `json:"-"`
		Embedding   *pgvector.Vector `json:"-"`
	}

	// MediaitemFace ...
	MediaitemFace struct {
		ID          uuid.UUID        `json:"-"`
		MediaitemID uuid.UUID        `json:"-"`
		PeopleID    *uuid.UUID       `json:"-"`
		Embedding   *pgvector.Vector `json:"-"`
		Thumbnail   string           `json:"thumbnail"`
	}

	// MediaItem ...
	MediaItem struct {
		ID                uuid.UUID         `json:"id"`
		UserID            uuid.UUID         `json:"userId"`
		Filename          string            `json:"filename"`
		Hash              *string           `json:"hash,omitempty"`
		Description       *string           `json:"description,omitempty"`
		MimeType          string            `json:"mimeType"`
		SourceURL         string            `json:"sourceUrl"`
		PreviewURL        string            `json:"previewUrl"`
		ThumbnailURL      string            `json:"thumbnailUrl"`
		Placeholder       string            `json:"placeholder"`
		IsFavourite       *bool             `json:"favourite"`
		IsHidden          *bool             `json:"hidden"`
		IsDeleted         *bool             `json:"deleted"`
		Status            MediaItemStatus   `json:"status"`
		MediaItemType     MediaItemType     `json:"mediaItemType"`
		MediaItemCategory MediaItemCategory `json:"mediaItemCategory"`
		Width             int               `json:"width"`
		Height            int               `json:"height"`
		CreationTime      time.Time         `json:"creationTime"`
		CameraMake        *string           `json:"cameraMake,omitempty"`
		CameraModel       *string           `json:"cameraModel,omitempty"`
		FocalLength       *string           `json:"focalLength,omitempty"`
		ApertureFnumber   *string           `json:"apertureFNumber,omitempty"`
		IsoEquivalent     *string           `json:"isoEquivalent,omitempty"`
		ExposureTime      *string           `json:"exposureTime,omitempty"`
		Latitude          *float64          `json:"latitude,omitempty"`
		Longitude         *float64          `json:"longitude,omitempty"`
		FPS               *string           `json:"fps,omitempty"`
		EXIFData          *string           `json:"-"`
		Keywords          *string           `json:"-"`
		CreatedAt         time.Time         `json:"createdAt"`
		UpdatedAt         time.Time         `json:"updatedAt"`
	}
)

const (
	StatusUnspecified MediaItemStatus = "UNSPECIFIED"
	StatusProcessing  MediaItemStatus = "PROCESSING"
	StatusReady       MediaItemStatus = "READY"
	StatusFailed      MediaItemStatus = "FAILED"

	TypeUnknown MediaItemType = "unknown"
	TypePhoto   MediaItemType = "photo"
	TypeVideo   MediaItemType = "video"

	CategoryDefault    MediaItemCategory = "default"
	CategoryScreenshot MediaItemCategory = "screenshot"
	CategoryPanorama   MediaItemCategory = "panorama"
	CategorySlow       MediaItemCategory = "slow"
	CategoryMotion     MediaItemCategory = "motion"
	CategoryLive       MediaItemCategory = "live"
	CategoryTimelapse  MediaItemCategory = "timelapse"

	preFetchTime = 24
)

// TableName ...
func (MediaItem) TableName() string {
	return MediaItemsTable
}

func ScanRowsToMediaItem(rows pgx.Rows) (MediaItem, error) {
	mediaItem := MediaItem{}
	err := rows.Scan(
		&mediaItem.ID,
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
	return mediaItem, err
}

// MediaItemURLPlugin ...
type MediaItemURLPlugin struct {
	Storage storage.Provider
	Cache   cache.Provider
}

// TransformMediaItemURL ...
func (m *MediaItemURLPlugin) TransformMediaItemURL(gormDB *gorm.DB) {
	if m.Storage.Type() == "disk" {
		return
	}
	if gormDB.Statement.Schema != nil {
		var mswg sync.WaitGroup
		mediaItemTypes := []string{"SourceURL", "PreviewURL", "ThumbnailURL"}
		mswg.Add(len(mediaItemTypes))
		for _, fieldName := range mediaItemTypes {
			go m.transformMediaItemURL(&mswg, gormDB, fieldName)
		}
		mswg.Wait()
	}
}

//nolint:gocognit,cyclop
func (m *MediaItemURLPlugin) transformMediaItemURL(
	wg *sync.WaitGroup,
	gormDB *gorm.DB,
	fieldName string,
) {
	defer wg.Done()
	field := gormDB.Statement.Schema.LookUpField(fieldName)
	if field != nil { //nolint: nestif
		switch gormDB.Statement.ReflectValue.Kind() { //nolint: exhaustive
		case reflect.Slice, reflect.Array:
			for i := range gormDB.Statement.ReflectValue.Len() {
				if fieldValue, isZero := field.ValueOf(gormDB.Statement.Context, gormDB.Statement.ReflectValue.Index(i)); !isZero {
					if val, ok := fieldValue.(string); ok {
						err := field.Set(
							gormDB.Statement.Context,
							gormDB.Statement.ReflectValue.Index(i),
							m.getMediaItemURL(fieldName, val),
						)
						if err != nil {
							slog.Error(
								"error setting field value",
								"field",
								fieldName,
								"value",
								val,
								"error",
								err,
							)
						}
					}
				}
			}
		case reflect.Struct:
			if fieldValue, isZero := field.ValueOf(gormDB.Statement.Context, gormDB.Statement.ReflectValue); !isZero {
				if val, ok := fieldValue.(string); ok {
					err := field.Set(
						gormDB.Statement.Context,
						gormDB.Statement.ReflectValue,
						m.getMediaItemURL(fieldName, val),
					)
					if err != nil {
						slog.Error(
							"error setting value for field",
							"field",
							fieldName,
							"value",
							val,
							"error",
							err,
						)
					}
				}
			}
		}
	}
}

func (m *MediaItemURLPlugin) getMediaItemURL(
	fieldName, filePath string,
) string {
	// get from cache if exists
	preFetchedVal, err := m.Cache.Get(filePath)
	if err == nil {
		if preFetchedURL, ok := preFetchedVal.(string); ok {
			return preFetchedURL
		}
	}

	slog.Error("error getting mediaitem url from cache", slog.Any("error", err))

	// generate from storage provider and add to cache
	fileType := getFileType(fieldName)
	fileID := strings.ReplaceAll(filePath, fmt.Sprintf("/%s/", fileType), "")

	fetchedURL, err := m.Storage.Get(fileType, fileID)
	if err != nil {
		slog.Error(
			"error getting mediaitem url from storage",
			slog.Any("error", err),
		)
		return ""
	}

	err = m.Cache.SetWithExpire(filePath, fetchedURL, preFetchTime*time.Hour)
	if err != nil {
		slog.Error(
			"error caching mediaitem url from storage",
			slog.Any("error", err),
		)
	}

	return fetchedURL
}

func getFileType(fieldName string) string {
	switch fieldName {
	case "SourceURL":
		return "originals"
	case "PreviewURL":
		return "previews"
	case "ThumbnailURL":
		return "thumbnails"
	}
	return "unknown"
}

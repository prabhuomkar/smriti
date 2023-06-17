package models

import (
	"api/pkg/cache"
	"api/pkg/storage"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

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

	// MediaItem ...
	//nolint: lll
	MediaItem struct {
		ID                uuid.UUID         `json:"id" gorm:"primaryKey;index:,unique;type:uuid"`
		UserID            uuid.UUID         `json:"userId" gorm:"column:user_id"`
		Filename          string            `json:"filename"`
		Description       *string           `json:"description,omitempty"`
		MimeType          string            `json:"mimeType"`
		SourceURL         string            `json:"sourceUrl" gorm:"column:source_url"`
		PreviewURL        string            `json:"previewUrl" gorm:"column:preview_url"`
		ThumbnailURL      string            `json:"thumbnailUrl" gorm:"column:thumbnail_url"`
		IsFavourite       *bool             `json:"favourite" gorm:"column:is_favourite;default:false"`
		IsHidden          *bool             `json:"hidden" gorm:"column:is_hidden;default:false"`
		IsDeleted         *bool             `json:"deleted" gorm:"column:is_deleted;default:false"`
		Status            MediaItemStatus   `json:"status"`
		MediaItemType     MediaItemType     `json:"mediaItemType" gorm:"column:mediaitem_type"`
		MediaItemCategory MediaItemCategory `json:"mediaItemCategory" gorm:"column:mediaitem_category"`
		Width             int               `json:"width"`
		Height            int               `json:"height"`
		CreationTime      time.Time         `json:"creationTime"`
		CameraMake        *string           `json:"cameraMake,omitempty"`
		CameraModel       *string           `json:"cameraModel,omitempty"`
		FocalLength       *string           `json:"focalLength,omitempty"`
		ApertureFnumber   *string           `json:"apertureFNumber,omitempty" gorm:"column:aperture_fnumber"`
		IsoEquivalent     *string           `json:"isoEquivalent,omitempty"`
		ExposureTime      *string           `json:"exposureTime,omitempty"`
		Latitude          *float64          `json:"latitude,omitempty"`
		Longitude         *float64          `json:"longitude,omitempty"`
		FPS               *string           `json:"fps,omitempty"`
		CreatedAt         time.Time         `json:"createdAt"`
		UpdatedAt         time.Time         `json:"updatedAt"`
		Albums            []*Album          `json:"-" gorm:"many2many:album_mediaitems;foreignKey:ID;joinForeignKey:MediaitemID;references:ID;joinReferences:AlbumID"`
		Places            []*Place          `json:"-" gorm:"many2many:place_mediaitems;foreignKey:ID;joinForeignKey:MediaitemID;references:ID;joinReferences:PlaceID"`
		Things            []*Thing          `json:"-" gorm:"many2many:thing_mediaitems;foreignKey:ID;joinForeignKey:MediaitemID;references:ID;joinReferences:ThingID"`
		People            []*People         `json:"-" gorm:"many2many:people_mediaitems;foreignKey:ID;joinForeignKey:MediaitemID;references:ID;joinReferences:PeopleID"`
	}
)

const (
	Unspecified MediaItemStatus = "UNSPECIFIED"
	Processing  MediaItemStatus = "PROCESSING"
	Ready       MediaItemStatus = "READY"
	Failed      MediaItemStatus = "FAILED"

	Unknown MediaItemType = "unknown"
	Photo   MediaItemType = "photo"
	Video   MediaItemType = "video"

	Default    MediaItemCategory = "default"
	Screenshot MediaItemCategory = "screenshot"
	Panorama   MediaItemCategory = "panorama"
	Slow       MediaItemCategory = "slow"
	Motion     MediaItemCategory = "motion"
	Live       MediaItemCategory = "live"

	preFetchTime = 24
)

// TableName ...
func (MediaItem) TableName() string {
	return MediaItemsTable
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

func (m *MediaItemURLPlugin) transformMediaItemURL(wg *sync.WaitGroup, gormDB *gorm.DB, fieldName string) { //nolint: gocognit, cyclop, lll
	defer wg.Done()
	field := gormDB.Statement.Schema.LookUpField(fieldName)
	if field != nil { //nolint: nestif
		switch gormDB.Statement.ReflectValue.Kind() { //nolint: exhaustive
		case reflect.Slice, reflect.Array:
			for i := 0; i < gormDB.Statement.ReflectValue.Len(); i++ {
				if fieldValue, isZero := field.ValueOf(gormDB.Statement.Context, gormDB.Statement.ReflectValue.Index(i)); !isZero {
					if val, ok := fieldValue.(string); ok {
						err := field.Set(gormDB.Statement.Context, gormDB.Statement.ReflectValue.Index(i),
							m.getMediaItemURL(fieldName, val))
						if err != nil {
							log.Printf("error setting %s value for %s: %+v", fieldName, val, err)
						}
					}
				}
			}
		case reflect.Struct:
			if fieldValue, isZero := field.ValueOf(gormDB.Statement.Context, gormDB.Statement.ReflectValue); !isZero {
				if val, ok := fieldValue.(string); ok {
					err := field.Set(gormDB.Statement.Context, gormDB.Statement.ReflectValue,
						m.getMediaItemURL(fieldName, val))
					if err != nil {
						log.Printf("error setting %s value for %s: %+v", fieldName, val, err)
					}
				}
			}
		}
	}
}

func (m *MediaItemURLPlugin) getMediaItemURL(fieldName, filePath string) string {
	// get from cache if exists
	preFetchedVal, err := m.Cache.Get(filePath)
	if err == nil {
		if preFetchedURL, ok := preFetchedVal.(string); ok {
			return preFetchedURL
		}
	}

	log.Printf("error getting mediaitem url from cache: %+v", err)

	// generate from storage provider and add to cache
	fileType := getFileType(fieldName)
	fileID := strings.ReplaceAll(filePath, fmt.Sprintf("/%s/", fileType), "")

	fetchedURL, err := m.Storage.Get(fileType, fileID)
	if err != nil {
		log.Printf("error getting mediaitem url from storage: %+v", err)
		return ""
	}

	err = m.Cache.SetWithExpire(filePath, fetchedURL, preFetchTime*time.Hour)
	if err != nil {
		log.Printf("error caching mediaitem url from storage: %+v", err)
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

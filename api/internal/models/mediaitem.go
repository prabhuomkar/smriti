package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
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
		SourceURL         string            `json:"sourceUrl"`
		PreviewURL        string            `json:"previewUrl"`
		ThumbnailURL      string            `json:"thumbnailUrl"`
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
)

// TableName ...
func (MediaItem) TableName() string {
	return MediaItemsTable
}

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

	// MediaItem ...
	MediaItem struct {
		ID              uuid.UUID       `json:"id" gorm:"primaryKey"`
		Filename        string          `json:"filename"`
		Description     string          `json:"description"`
		MimeType        string          `json:"mimeType"`
		SourceURL       string          `json:"sourceUrl"`
		PreviewURL      string          `json:"previewUrl"`
		ThumbnailURL    string          `json:"thumbnailUrl"`
		IsFavourite     bool            `json:"favourite"`
		IsHidden        bool            `json:"hidden"`
		IsDeleted       bool            `json:"deleted"`
		Status          MediaItemStatus `json:"status"`
		MediaItemType   MediaItemType   `json:"mediaItemType" gorm:"column:mediaitem_type"`
		Width           int             `json:"width"`
		Height          int             `json:"height"`
		CreationTime    time.Time       `json:"creationTime"`
		CameraMake      string          `json:"cameraMake,omitempty"`
		CameraModel     string          `json:"cameraModel,omitempty"`
		FocalLength     string          `json:"focalLength,omitempty"`
		ApertureFnumber string          `json:"apertureFnumber,omitempty" gorm:"column:aperture_fnumber"`
		IsoEquivalent   string          `json:"isoEquivalent,omitempty"`
		ExposureTime    string          `json:"exposureTime,omitempty"`
		Location        []byte          `json:"location,omitempty"`
		FPS             string          `json:"fps,omitempty"`
		CreatedAt       time.Time       `json:"createdAt"`
		UpdatedAt       time.Time       `json:"updatedAt"`
		Albums          []*Album        `json:"-" gorm:"many2many:album_mediaitems;foreignKey:ID;joinForeignKey:MediaitemID"`
		Places          []*Place        `json:"-" gorm:"many2many:place_mediaitems;foreignKey:ID;joinForeignKey:MediaitemID"`
		Things          []*Thing        `json:"-" gorm:"many2many:thing_mediaitems;foreignKey:ID;joinForeignKey:MediaitemID"`
		People          []*People       `json:"-" gorm:"many2many:people_mediaitems;foreignKey:ID;joinForeignKey:MediaitemID"`
	}
)

const (
	Unspecified MediaItemStatus = "UNSPECIFIED"
	Processing  MediaItemStatus = "PROCESSING"
	Ready       MediaItemStatus = "READY"
	Failed      MediaItemStatus = "FAILED"

	Photo MediaItemType = "photo"
	Video MediaItemType = "video"
)

// TableName ...
func (MediaItem) TableName() string {
	return MediaItemsTable
}

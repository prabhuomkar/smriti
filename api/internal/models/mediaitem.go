package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

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
		SourceUrl       string          `json:"sourceUrl"`
		PreviewUrl      string          `json:"previewUrl"`
		ThumbnailUrl    string          `json:"thumbnailUrl"`
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
	return "mediaitems"
}

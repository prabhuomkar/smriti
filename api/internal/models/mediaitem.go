package models

import "time"

type (
	// MediaItemStatus ...
	MediaItemStatus string

	// MediaItemType ...
	MediaItemType string

	// MediaItem ...
	MediaItem struct {
		ID              string          `json:"id" db:"id"`
		Filename        string          `json:"filename" db:"filename"`
		Description     string          `json:"description" db:"description"`
		MimeType        string          `json:"mimeType" db:"mime_type"`
		SourceUrl       string          `json:"sourceUrl" db:"source_url"`
		PreviewUrl      string          `json:"previewUrl" db:"preview_url"`
		ThumbnailUrl    string          `json:"thumbnailUrl" db:"thumbnail_url"`
		IsFavourite     bool            `json:"-" db:"is_favourite"`
		IsHidden        bool            `json:"-" db:"is_hidden"`
		IsDeleted       bool            `json:"-" db:"is_deleted"`
		Status          MediaItemStatus `json:"status" db:"status"`
		MediaItemType   MediaItemType   `json:"mediaItemType" db:"mediaitem_type"`
		Width           int             `json:"width" db:"width"`
		Height          int             `json:"height" db:"height"`
		CreationTime    time.Time       `json:"creationTime" db:"creation_time"`
		CameraMake      string          `json:"cameraMake" db:"camera_make"`
		CameraModel     string          `json:"cameraModel" db:"camera_model"`
		FocalLength     string          `json:"focalLength" db:"focal_length"`
		ApertureFnumber string          `json:"apertureFnumber" db:"aperture_fnumber"`
		IsoEquivalent   string          `json:"isoEquivalent" db:"iso_equivalent"`
		ExposureTime    string          `json:"exposureTime" db:"exposure_time"`
		Location        string          `json:"location" db:"location"`
		FPS             string          `json:"fps" db:"fps"`
		CreatedAt       time.Time       `json:"createdAt" db:"created_at"`
		UpdatedAt       time.Time       `json:"updatedAt" db:"updated_at"`
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

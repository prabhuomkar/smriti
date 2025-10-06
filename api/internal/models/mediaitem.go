package models

import (
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	uuid "github.com/satori/go.uuid"
)

const MediaItemsTable = "mediaitems"

type ( // MediaItemStatus ...
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
		ID                uuid.UUID `json:"id"`
		UserID            uuid.UUID `json:"userId"`
		Filename          string    `json:"filename"`
		Hash              *string   `json:"hash,omitempty"`
		Description       *string   `json:"description,omitempty"`
		MimeType          string    `json:"mimeType"`
		SourceURL         string    `json:"sourceUrl"`
		PreviewURL        string    `json:"previewUrl"`
		ThumbnailURL      string    `json:"thumbnailUrl"`
		Placeholder       string    `json:"placeholder"`
		IsFavourite       *bool     `json:"favourite"`
		IsHidden          *bool     `json:"hidden"`
		IsDeleted         *bool     `json:"deleted"`
		Status            string    `json:"status"`
		MediaItemType     string    `json:"mediaItemType"`
		MediaItemCategory string    `json:"mediaItemCategory"`
		Width             int       `json:"width"`
		Height            int       `json:"height"`
		CreationTime      time.Time `json:"creationTime"`
		CameraMake        *string   `json:"cameraMake,omitempty"`
		CameraModel       *string   `json:"cameraModel,omitempty"`
		FocalLength       *string   `json:"focalLength,omitempty"`
		ApertureFnumber   *string   `json:"apertureFNumber,omitempty"`
		IsoEquivalent     *string   `json:"isoEquivalent,omitempty"`
		ExposureTime      *string   `json:"exposureTime,omitempty"`
		Megapixels        *string   `json:"megapixels,omitempty"`
		Latitude          *float64  `json:"latitude,omitempty"`
		Longitude         *float64  `json:"longitude,omitempty"`
		FPS               *string   `json:"fps,omitempty"`
		EXIFData          *string   `json:"-"`
		DetectedText      *string   `json:"-"`
		Caption           *string   `json:"-"`
		CreatedAt         time.Time `json:"createdAt"`
		UpdatedAt         time.Time `json:"updatedAt"`
	}

	// CoverMediaItem ...
	CoverMediaItem struct {
		ID                *uuid.UUID `json:"id"`
		UserID            *uuid.UUID `json:"userId"`
		SourceURL         *string    `json:"sourceUrl"`
		PreviewURL        *string    `json:"previewUrl"`
		ThumbnailURL      *string    `json:"thumbnailUrl"`
		Placeholder       *string    `json:"placeholder"`
		MediaItemType     *string    `json:"mediaItemType"`
		MediaItemCategory *string    `json:"mediaItemCategory"`
		Width             *int       `json:"width"`
		Height            *int       `json:"height"`
	}
)

// TableName ...
func (MediaItem) TableName() string {
	return MediaItemsTable
}

func ScanRowsToMediaItem(rows pgx.Rows) (MediaItem, error) {
	mediaItem := MediaItem{}
	err := rows.Scan(&mediaItem.ID, &mediaItem.UserID, &mediaItem.Filename, &mediaItem.Hash, &mediaItem.Description,
		&mediaItem.MimeType, &mediaItem.SourceURL, &mediaItem.PreviewURL, &mediaItem.ThumbnailURL,
		&mediaItem.Placeholder, &mediaItem.IsFavourite, &mediaItem.IsHidden, &mediaItem.IsDeleted, &mediaItem.Status,
		&mediaItem.MediaItemType, &mediaItem.MediaItemCategory, &mediaItem.Width, &mediaItem.Height,
		&mediaItem.CreationTime, &mediaItem.CameraMake, &mediaItem.CameraModel, &mediaItem.FocalLength,
		&mediaItem.ApertureFnumber, &mediaItem.IsoEquivalent, &mediaItem.ExposureTime, &mediaItem.Megapixels,
		&mediaItem.Latitude, &mediaItem.Longitude, &mediaItem.FPS, &mediaItem.EXIFData, &mediaItem.DetectedText,
		&mediaItem.Caption, &mediaItem.CreatedAt, &mediaItem.UpdatedAt)

	return mediaItem, err
}

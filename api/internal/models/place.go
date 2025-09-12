package models

import (
	"time"

	"github.com/jackc/pgx/v5"
	uuid "github.com/satori/go.uuid"
)

const PlaceTable = "places"

// Place ...
type Place struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uuid.UUID  `json:"userId"`
	Name             string     `json:"name"`
	Postcode         *string    `json:"postcode"`
	Country          *string    `json:"country"`
	Locality         *string    `json:"locality"`
	Area             *string    `json:"area"`
	IsHidden         *bool      `json:"hidden"`
	CoverMediaItemID *uuid.UUID `json:"coverMediaItemId,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	CoverMediaItem   *MediaItem `json:"coverMediaItem,omitempty"`
}

// TableName ...
func (Place) TableName() string {
	return PlaceTable
}

func ScanRowsToPlace(rows pgx.Rows) (Place, error) {
	place := Place{CoverMediaItem: &MediaItem{}}
	err := rows.Scan(
		&place.ID,
		&place.UserID,
		&place.Name,
		&place.Postcode,
		&place.Country,
		&place.Locality,
		&place.Area,
		&place.IsHidden,
		&place.CoverMediaItemID,
		&place.CreatedAt,
		&place.UpdatedAt,
		&place.CoverMediaItem.ID,
		&place.CoverMediaItem.UserID,
		&place.CoverMediaItem.Filename,
		&place.CoverMediaItem.Hash,
		&place.CoverMediaItem.Description,
		&place.CoverMediaItem.MimeType,
		&place.CoverMediaItem.SourceURL,
		&place.CoverMediaItem.PreviewURL,
		&place.CoverMediaItem.ThumbnailURL,
		&place.CoverMediaItem.Placeholder,
		&place.CoverMediaItem.IsFavourite,
		&place.CoverMediaItem.IsHidden,
		&place.CoverMediaItem.IsDeleted,
		&place.CoverMediaItem.Status,
		&place.CoverMediaItem.MediaItemType,
		&place.CoverMediaItem.MediaItemCategory,
		&place.CoverMediaItem.Width,
		&place.CoverMediaItem.Height,
		&place.CoverMediaItem.CreationTime,
		&place.CoverMediaItem.CameraMake,
		&place.CoverMediaItem.CameraModel,
		&place.CoverMediaItem.FocalLength,
		&place.CoverMediaItem.ApertureFnumber,
		&place.CoverMediaItem.IsoEquivalent,
		&place.CoverMediaItem.ExposureTime,
		&place.CoverMediaItem.Latitude,
		&place.CoverMediaItem.Longitude,
		&place.CoverMediaItem.FPS,
		&place.CoverMediaItem.EXIFData,
		&place.CoverMediaItem.Keywords,
		&place.CoverMediaItem.CreatedAt,
		&place.CoverMediaItem.UpdatedAt)

	return place, err
}

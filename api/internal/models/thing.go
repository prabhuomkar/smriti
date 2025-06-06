package models

import (
	"time"

	"github.com/jackc/pgx/v5"
	uuid "github.com/satori/go.uuid"
)

const ThingTable = "things"

// Thing ...
type Thing struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uuid.UUID  `json:"userId"`
	Name             string     `json:"name"`
	IsHidden         *bool      `json:"hidden"`
	CoverMediaItemID *uuid.UUID `json:"coverMediaItemId,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	CoverMediaItem   *MediaItem `json:"coverMediaItem,omitempty"`
}

// TableName ...
func (Thing) TableName() string {
	return ThingTable
}

func ScanRowsToThing(rows pgx.Rows) (Thing, error) {
	thing := Thing{CoverMediaItem: &MediaItem{}}
	err := rows.Scan(
		&thing.ID,
		&thing.UserID,
		&thing.Name,
		&thing.IsHidden,
		&thing.CoverMediaItemID,
		&thing.CreatedAt,
		&thing.UpdatedAt,
		&thing.CoverMediaItem.ID,
		&thing.CoverMediaItem.UserID,
		&thing.CoverMediaItem.Filename,
		&thing.CoverMediaItem.Hash,
		&thing.CoverMediaItem.Description,
		&thing.CoverMediaItem.MimeType,
		&thing.CoverMediaItem.SourceURL,
		&thing.CoverMediaItem.PreviewURL,
		&thing.CoverMediaItem.ThumbnailURL,
		&thing.CoverMediaItem.Placeholder,
		&thing.CoverMediaItem.IsFavourite,
		&thing.CoverMediaItem.IsHidden,
		&thing.CoverMediaItem.IsDeleted,
		&thing.CoverMediaItem.Status,
		&thing.CoverMediaItem.MediaItemType,
		&thing.CoverMediaItem.MediaItemCategory,
		&thing.CoverMediaItem.Width,
		&thing.CoverMediaItem.Height,
		&thing.CoverMediaItem.CreationTime,
		&thing.CoverMediaItem.CameraMake,
		&thing.CoverMediaItem.CameraModel,
		&thing.CoverMediaItem.FocalLength,
		&thing.CoverMediaItem.ApertureFnumber,
		&thing.CoverMediaItem.IsoEquivalent,
		&thing.CoverMediaItem.ExposureTime,
		&thing.CoverMediaItem.Latitude,
		&thing.CoverMediaItem.Longitude,
		&thing.CoverMediaItem.FPS,
		&thing.CoverMediaItem.EXIFData,
		&thing.CoverMediaItem.Keywords,
		&thing.CoverMediaItem.CreatedAt,
		&thing.CoverMediaItem.UpdatedAt)
	return thing, err
}

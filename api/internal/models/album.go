package models

import (
	"time"

	"github.com/jackc/pgx/v5"
	uuid "github.com/satori/go.uuid"
)

const AlbumsTable = "albums"

// Album ...
type Album struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uuid.UUID  `json:"userId"`
	Name             string     `json:"name"`
	Description      *string    `json:"description"`
	IsShared         *bool      `json:"shared,omitempty"`
	IsHidden         *bool      `json:"hidden,omitempty"`
	MediaItemsCount  *int       `json:"mediaItemsCount,omitempty"`
	CoverMediaItemID *uuid.UUID `json:"coverMediaItemId,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	CoverMediaItem   *MediaItem `json:"coverMediaItem"`
}

// TableName ...
func (Album) TableName() string {
	return AlbumsTable
}

func ScanRowsToAlbum(rows pgx.Rows) (Album, error) {
	album := Album{CoverMediaItem: &MediaItem{}}
	err := rows.Scan(&album.ID, &album.UserID, &album.Name, &album.Description, &album.IsShared, &album.IsHidden,
		&album.MediaItemsCount, &album.CoverMediaItemID, &album.CreatedAt, &album.UpdatedAt, &album.CoverMediaItem.ID,
		&album.CoverMediaItem.UserID, &album.CoverMediaItem.Filename, &album.CoverMediaItem.Hash,
		&album.CoverMediaItem.Description, &album.CoverMediaItem.MimeType, &album.CoverMediaItem.SourceURL,
		&album.CoverMediaItem.PreviewURL, &album.CoverMediaItem.ThumbnailURL, &album.CoverMediaItem.Placeholder,
		&album.CoverMediaItem.IsFavourite, &album.CoverMediaItem.IsHidden, &album.CoverMediaItem.IsDeleted,
		&album.CoverMediaItem.Status, &album.CoverMediaItem.MediaItemType, &album.CoverMediaItem.MediaItemCategory,
		&album.CoverMediaItem.Width, &album.CoverMediaItem.Height, &album.CoverMediaItem.CreationTime,
		&album.CoverMediaItem.CameraMake, &album.CoverMediaItem.CameraModel, &album.CoverMediaItem.FocalLength,
		&album.CoverMediaItem.ApertureFnumber, &album.CoverMediaItem.IsoEquivalent, &album.CoverMediaItem.ExposureTime,
		&album.CoverMediaItem.Megapixels, &album.CoverMediaItem.Latitude, &album.CoverMediaItem.Longitude,
		&album.CoverMediaItem.FPS, &album.CoverMediaItem.EXIFData, &album.CoverMediaItem.Keywords,
		&album.CoverMediaItem.CreatedAt, &album.CoverMediaItem.UpdatedAt)

	return album, err
}

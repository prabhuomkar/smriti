package models

import (
	"time"

	"github.com/jackc/pgx/v5"
	uuid "github.com/satori/go.uuid"
)

const AlbumsTable = "albums"

// Album ...
type Album struct {
	ID               uuid.UUID       `json:"id"`
	UserID           uuid.UUID       `json:"userId"`
	Name             string          `json:"name"`
	Description      *string         `json:"description"`
	IsShared         *bool           `json:"shared,omitempty"`
	IsHidden         *bool           `json:"hidden,omitempty"`
	MediaItemsCount  *int            `json:"mediaItemsCount,omitempty"`
	CoverMediaItemID *uuid.UUID      `json:"coverMediaItemId,omitempty"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
	CoverMediaItem   *CoverMediaItem `json:"coverMediaItem"`
}

// TableName ...
func (Album) TableName() string {
	return AlbumsTable
}

func ScanRowsToAlbum(rows pgx.Rows) (Album, error) {
	album := Album{}
	coverMediaItem := &CoverMediaItem{}

	err := rows.Scan(&album.ID, &album.UserID, &album.Name, &album.Description, &album.IsShared, &album.IsHidden,
		&album.MediaItemsCount, &album.CoverMediaItemID, &album.CreatedAt, &album.UpdatedAt, &coverMediaItem.ID,
		&coverMediaItem.UserID, &coverMediaItem.SourceURL, &coverMediaItem.PreviewURL,
		&coverMediaItem.ThumbnailURL, &coverMediaItem.Placeholder, &coverMediaItem.MediaItemType,
		&coverMediaItem.MediaItemCategory, &coverMediaItem.Width, &coverMediaItem.Height)
	if err == nil && album.CoverMediaItemID != nil {
		album.CoverMediaItem = coverMediaItem
	}

	return album, err
}

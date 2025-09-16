package models

import (
	"time"

	"github.com/jackc/pgx/v5"
	uuid "github.com/satori/go.uuid"
)

const ThingTable = "things"

// Thing ...
type Thing struct {
	ID               uuid.UUID       `json:"id"`
	UserID           uuid.UUID       `json:"userId"`
	Name             string          `json:"name"`
	IsHidden         *bool           `json:"hidden"`
	CoverMediaItemID *uuid.UUID      `json:"coverMediaItemId,omitempty"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
	CoverMediaItem   *CoverMediaItem `json:"coverMediaItem,omitempty"`
}

// TableName ...
func (Thing) TableName() string {
	return ThingTable
}

func ScanRowsToThing(rows pgx.Rows) (Thing, error) {
	thing := Thing{}
	coverMediaItem := &CoverMediaItem{}

	err := rows.Scan(&thing.ID, &thing.UserID, &thing.Name, &thing.IsHidden, &thing.CoverMediaItemID, &thing.CreatedAt,
		&thing.UpdatedAt, &coverMediaItem.ID,
		&coverMediaItem.UserID, &coverMediaItem.SourceURL, &coverMediaItem.PreviewURL,
		&coverMediaItem.ThumbnailURL, &coverMediaItem.Placeholder, &coverMediaItem.MediaItemType,
		&coverMediaItem.MediaItemCategory, &coverMediaItem.Width, &coverMediaItem.Height)
	if err == nil && thing.CoverMediaItemID != nil {
		thing.CoverMediaItem = coverMediaItem
	}

	return thing, err
}

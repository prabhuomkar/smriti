package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const PlaceTable = "places"

// Place ...
type Place struct {
	ID               uuid.UUID       `json:"id"`
	UserID           uuid.UUID       `json:"userId"`
	Name             string          `json:"name"`
	Postcode         *string         `json:"postcode"`
	Country          *string         `json:"country"`
	Locality         *string         `json:"locality"`
	Area             *string         `json:"area"`
	IsHidden         *bool           `json:"hidden"`
	CoverMediaItemID *uuid.UUID      `json:"coverMediaItemId,omitempty"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
	CoverMediaItem   *CoverMediaItem `json:"coverMediaItem,omitempty"`
}

// TableName ...
func (Place) TableName() string {
	return PlaceTable
}

func ScanRowsToPlace(rows pgx.Rows) (Place, error) {
	place := Place{}
	coverMediaItem := &CoverMediaItem{}

	err := rows.Scan(&place.ID, &place.UserID, &place.Name, &place.Postcode, &place.Country,
		&place.Locality, &place.Area, &place.IsHidden, &place.CoverMediaItemID, &place.CreatedAt,
		&place.UpdatedAt, &coverMediaItem.ID,
		&coverMediaItem.UserID, &coverMediaItem.SourceURL, &coverMediaItem.PreviewURL,
		&coverMediaItem.ThumbnailURL, &coverMediaItem.Placeholder, &coverMediaItem.MediaItemType,
		&coverMediaItem.MediaItemCategory, &coverMediaItem.Width, &coverMediaItem.Height)
	if err == nil && place.CoverMediaItemID != nil {
		place.CoverMediaItem = coverMediaItem
	}

	return place, err
}

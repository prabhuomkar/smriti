package models

import (
	"time"

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

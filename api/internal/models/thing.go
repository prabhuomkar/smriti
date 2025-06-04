package models

import (
	"time"

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

package models

import (
	"time"

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

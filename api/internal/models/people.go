package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const PeopleTable = "people"

// People ...
type People struct {
	ID                   uuid.UUID      `json:"id"`
	UserID               uuid.UUID      `json:"userId"`
	Name                 string         `json:"name"`
	IsHidden             *bool          `json:"hidden"`
	CoverMediaItemID     *uuid.UUID     `json:"coverMediaItemId,omitempty"`
	CoverMediaItemFaceID *uuid.UUID     `json:"coverMediaItemFaceId,omitempty"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	CoverMediaItem       *MediaItem     `json:"coverMediaItem,omitempty"`
	CoverMediaItemFace   *MediaitemFace `json:"coverMediaItemFace,omitempty"`
}

// TableName ...
func (People) TableName() string {
	return PeopleTable
}

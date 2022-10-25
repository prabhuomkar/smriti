package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const AlbumsTable = "albums"

// Album ...
type Album struct {
	ID               uuid.UUID    `json:"id" gorm:"primaryKey"`
	UserID           uuid.UUID    `json:"userId" gorm:"column:user_id"`
	Name             string       `json:"name"`
	Description      *string      `json:"description"`
	IsShared         *bool        `json:"shared,omitempty" gorm:"column:is_shared;default:false"`
	IsHidden         *bool        `json:"hidden,omitempty" gorm:"column:is_hidden;default:false"`
	MediaItemsCount  *int         `json:"mediaItemsCount,omitempty" gorm:"column:mediaitems_count;default:0"`
	CoverMediaItemID *uuid.UUID   `json:"coverMediaItemId,omitempty" gorm:"column:cover_mediaitem_id"`
	CreatedAt        time.Time    `json:"createdAt"`
	UpdatedAt        time.Time    `json:"updatedAt"`
	CoverMediaItem   *MediaItem   `json:"coverMediaItem" gorm:"foreignkey:ID;references:CoverMediaItemID"`
	MediaItems       []*MediaItem `json:"-" gorm:"many2many:album_mediaitems;References:ID;joinReferences:MediaitemID"`
}

// TableName ...
func (Album) TableName() string {
	return AlbumsTable
}

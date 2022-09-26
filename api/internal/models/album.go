package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const AlbumsTable = "albums"

// Album ...
type Album struct {
	ID                         uuid.UUID    `json:"id" gorm:"primaryKey"`
	Name                       string       `json:"name"`
	Description                *string      `json:"description"`
	IsShared                   *bool        `json:"shared,omitempty" gorm:"column:is_shared;default:false"`
	IsHidden                   *bool        `json:"hidden,omitempty" gorm:"column:is_hidden;default:false"`
	MediaItemsCount            *int         `json:"mediaItemsCount,omitempty" gorm:"column:mediaitems_count;default:0"`
	CoverMediaItemID           *uuid.UUID   `json:"coverMediaItemId,omitempty" gorm:"column:cover_mediaitem_id"`
	CoverMediaItemThumbnailUrl *string      `json:"coverMediaItemThumbnailUrl,omitempty" gorm:"column:cover_mediaitem_thumbnail_url"`
	CreatedAt                  time.Time    `json:"createdAt"`
	UpdatedAt                  time.Time    `json:"updatedAt"`
	MediaItems                 []*MediaItem `json:"-" gorm:"many2many:album_mediaitems;References:ID;joinReferences:MediaitemID"`
}

// TableName ...
func (Album) TableName() string {
	return AlbumsTable
}

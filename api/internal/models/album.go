package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Album ...
type Album struct {
	ID                         uuid.UUID `json:"id" gorm:"primaryKey"`
	Name                       string    `json:"name"`
	Description                string    `json:"description"`
	IsShared                   bool      `json:"shared" gorm:"column:is_shared"`
	IsHidden                   bool      `json:"hidden" gorm:"column:is_hidden"`
	MediaItemsCount            int       `json:"mediaItemsCount" gorm:"column:mediaitems_count"`
	CoverMediaItemID           uuid.UUID `json:"coverMediaItemId" gorm:"column:cover_mediaitem_id"`
	CoverMediaItemThumbnailUrl string    `json:"coverMediaItemThumbnailUrl" gorm:"column:cover_mediaitem_thumbnail_url"`
	CreatedAt                  time.Time `json:"createdAt"`
	UpdatedAt                  time.Time `json:"updatedAt"`
}

// TableName ...
func (Album) TableName() string {
	return "albums"
}

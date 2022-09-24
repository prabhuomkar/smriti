package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Album ...
type Album struct {
	ID                         string    `json:"id" db:"id"`
	Name                       string    `json:"name" db:"name"`
	Description                string    `json:"description" db:"description"`
	IsShared                   bool      `json:"-" db:"is_shared"`
	IsHidden                   bool      `json:"-" db:"is_hidden"`
	MediaItemsCount            int       `json:"mediaitemsCount" db:"mediaitems_count"`
	CoverMediaItemID           string    `json:"coverMediaItemId" db:"cover_mediaitem_id"`
	CoverMediaItemThumbnailUrl string    `json:"coverMediaItemThumbnailUrl" db:"cover_mediaitem_thumbnail_url"`
	CreatedAt                  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt                  time.Time `json:"updatedAt" db:"updated_at"`
}

func (a *Album) NewID() {
	a.ID = uuid.NewV4().String()
}

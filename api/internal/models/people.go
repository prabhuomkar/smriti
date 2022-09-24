package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// People ...
type People struct {
	ID                         string    `json:"id" db:"id"`
	Name                       string    `json:"name" db:"name"`
	IsHidden                   bool      `json:"-" db:"is_hidden"`
	CoverMediaItemID           string    `json:"coverMediaItemId" db:"cover_mediaitem_id"`
	CoverMediaItemThumbnailUrl string    `json:"coverMediaItemThumbnailUrl" db:"cover_mediaitem_thumbnail_url"`
	CreatedAt                  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt                  time.Time `json:"updatedAt" db:"updated_at"`
}

func (p *People) NewID() {
	p.ID = uuid.NewV4().String()
}

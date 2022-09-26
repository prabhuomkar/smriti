package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const PeopleTable = "people"

// People ...
type People struct {
	ID                         uuid.UUID    `json:"id" gorm:"primaryKey"`
	Name                       string       `json:"name"`
	IsHidden                   bool         `json:"hidden"`
	CoverMediaItemID           uuid.UUID    `json:"coverMediaItemId" gorm:"column:cover_mediaitem_id"`
	CoverMediaItemThumbnailURL string       `json:"coverMediaItemThumbnailUrl" gorm:"column:cover_mediaitem_thumbnail_url"`
	CreatedAt                  time.Time    `json:"createdAt"`
	UpdatedAt                  time.Time    `json:"updatedAt"`
	MediaItems                 []*MediaItem `json:"-" gorm:"many2many:people_mediaitems;References:ID;joinReferences:MediaitemID"`
}

// TableName ...
func (People) TableName() string {
	return PeopleTable
}

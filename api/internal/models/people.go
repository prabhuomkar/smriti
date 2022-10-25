package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const PeopleTable = "people"

// People ...
type People struct {
	ID               uuid.UUID    `json:"id" gorm:"primaryKey"`
	UserID           uuid.UUID    `json:"userId" gorm:"column:user_id"`
	Name             string       `json:"name"`
	IsHidden         *bool        `json:"hidden"`
	CoverMediaItemID uuid.UUID    `json:"coverMediaItemId" gorm:"column:cover_mediaitem_id"`
	CreatedAt        time.Time    `json:"createdAt"`
	UpdatedAt        time.Time    `json:"updatedAt"`
	CoverMediaItem   *MediaItem   `json:"coverMediaItem" gorm:"foreignkey:ID;references:CoverMediaItemID"`
	MediaItems       []*MediaItem `json:"-" gorm:"many2many:people_mediaitems;References:ID;joinReferences:MediaitemID"`
}

// TableName ...
func (People) TableName() string {
	return PeopleTable
}

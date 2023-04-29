package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const PlaceTable = "places"

// Place ...
type Place struct {
	ID               uuid.UUID    `json:"id" gorm:"primaryKey;index:,unique;type:uuid"`
	UserID           uuid.UUID    `json:"userId" gorm:"column:user_id"`
	Name             string       `json:"name"`
	Postcode         *string      `json:"postcode"`
	Town             *string      `json:"town"`
	City             *string      `json:"city"`
	State            *string      `json:"state"`
	Country          *string      `json:"country"`
	IsHidden         *bool        `json:"hidden" gorm:"column:is_hidden;default:false"`
	CoverMediaItemID *uuid.UUID   `json:"coverMediaItemId" gorm:"column:cover_mediaitem_id;type:uuid"`
	CreatedAt        time.Time    `json:"createdAt"`
	UpdatedAt        time.Time    `json:"updatedAt"`
	CoverMediaItem   *MediaItem   `json:"coverMediaItem" gorm:"references:ID"`
	MediaItems       []*MediaItem `json:"-" gorm:"many2many:place_mediaitems;references:ID;joinReferences:MediaitemID"`
}

// TableName ...
func (Place) TableName() string {
	return PlaceTable
}

package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Place ...
type Place struct {
	ID                         uuid.UUID `json:"id" gorm:"primaryKey"`
	Name                       string    `json:"name"`
	Postcode                   string    `json:"postcode"`
	Suburb                     string    `json:"suburb"`
	Road                       string    `json:"road"`
	Town                       string    `json:"town"`
	City                       string    `json:"city"`
	County                     string    `json:"county"`
	District                   string    `json:"district"`
	State                      string    `json:"state"`
	Country                    string    `json:"country"`
	IsHidden                   bool      `json:"hidden"`
	CoverMediaItemID           uuid.UUID `json:"coverMediaItemId" gorm:"column:cover_mediaitem_id"`
	CoverMediaItemThumbnailUrl string    `json:"coverMediaItemThumbnailUrl" gorm:"column:cover_mediaitem_thumbnail_url"`
	CreatedAt                  time.Time `json:"createdAt"`
	UpdatedAt                  time.Time `json:"updatedAt"`
}

// TableName ...
func (Place) TableName() string {
	return "places"
}

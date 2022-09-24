package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Place ...
type Place struct {
	ID                         string    `json:"id" db:"id"`
	Name                       string    `json:"name" db:"name"`
	Postcode                   string    `json:"postcode" db:"postcode"`
	Suburb                     string    `json:"suburb" db:"suburb"`
	Road                       string    `json:"road" db:"road"`
	Town                       string    `json:"town" db:"town"`
	City                       string    `json:"city" db:"city"`
	County                     string    `json:"county" db:"county"`
	District                   string    `json:"district" db:"district"`
	State                      string    `json:"state" db:"state"`
	Country                    string    `json:"country" db:"country"`
	IsHidden                   bool      `json:"-" db:"is_hidden"`
	CoverMediaItemID           string    `json:"coverMediaItemId" db:"cover_mediaitem_id"`
	CoverMediaItemThumbnailUrl string    `json:"coverMediaItemThumbnailUrl" db:"cover_mediaitem_thumbnail_url"`
	CreatedAt                  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt                  time.Time `json:"updatedAt" db:"updated_at"`
}

func (p *Place) NewID() {
	p.ID = uuid.NewV4().String()
}

package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const PeopleTable = "people"

// People ...
type People struct {
	ID               uuid.UUID        `json:"id" gorm:"primaryKey;index:,unique;type:uuid"`
	UserID           uuid.UUID        `json:"userId" gorm:"column:user_id"`
	Name             string           `json:"name"`
	IsHidden         *bool            `json:"hidden" gorm:"column:is_hidden;default:false"`
	CoverMediaItemID *uuid.UUID       `json:"coverMediaItemId" gorm:"column:cover_mediaitem_id;type:uuid"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
	CoverMediaItem   *MediaItem       `json:"coverMediaItem" gorm:"references:ID"`
	MediaItems       []*MediaItem     `json:"-" gorm:"many2many:people_mediaitems;references:ID;joinReferences:MediaitemID"`
	Faces            []*MediaitemFace `json:"-" gorm:"many2many:mediaitem_faces;references:ID;joinReferences:MediaitemID"`
}

// TableName ...
func (People) TableName() string {
	return PeopleTable
}

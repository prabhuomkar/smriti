package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const PeopleTable = "people"

// People ...
type People struct {
	ID                   uuid.UUID      `json:"id"`
	UserID               uuid.UUID      `json:"userId"`
	Name                 string         `json:"name"`
	IsHidden             *bool          `json:"hidden"`
	CoverMediaItemID     *uuid.UUID     `json:"coverMediaItemId,omitempty"`
	CoverMediaItemFaceID *uuid.UUID     `json:"coverMediaItemFaceId,omitempty"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	CoverMediaItem       *MediaItem     `json:"coverMediaItem,omitempty"`
	CoverMediaItemFace   *MediaitemFace `json:"coverMediaItemFace,omitempty"`
}

// TableName ...
func (People) TableName() string {
	return PeopleTable
}

func ScanRowsToPerson(rows pgx.Rows) (People, error) {
	person := People{CoverMediaItemFace: &MediaitemFace{}}
	err := rows.Scan(&person.ID, &person.UserID, &person.Name, &person.IsHidden, &person.CoverMediaItemID,
		&person.CoverMediaItemFaceID, &person.CreatedAt, &person.UpdatedAt, &person.CoverMediaItemFace.ID,
		&person.CoverMediaItemFace.MediaitemID, &person.CoverMediaItemFace.PeopleID,
		&person.CoverMediaItemFace.Embedding, &person.CoverMediaItemFace.Thumbnail)

	return person, err
}

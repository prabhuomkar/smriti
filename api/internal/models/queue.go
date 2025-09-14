package models

import (
	uuid "github.com/satori/go.uuid"
)

const QueueTable = "queue"

// Queue ...
type Queue struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	MediaItemID uuid.UUID
	Type        string
	Components  string
	Status      MediaItemStatus
}

// TableName ...
func (Queue) TableName() string {
	return QueueTable
}

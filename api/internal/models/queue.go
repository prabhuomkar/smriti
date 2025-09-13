package models

import (
	uuid "github.com/satori/go.uuid"
)

const QueueTable = "queue"

// Queue ...
type Queue struct {
	ID         uuid.UUID       `json:"id"`
	Components string          `json:"components"`
	Status     MediaItemStatus `json:"status"`
}

// TableName ...
func (Queue) TableName() string {
	return QueueTable
}

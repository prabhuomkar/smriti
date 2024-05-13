package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const JobsTable = "jobs"

type (
	// JobStatus ...
	JobStatus string
)

// Job ...
type Job struct {
	ID             uuid.UUID  `json:"id" gorm:"primaryKey;index:,unique;type:uuid"`
	UserID         uuid.UUID  `json:"userId" gorm:"column:user_id"`
	Status         JobStatus  `json:"status"`
	Components     string     `json:"components"`
	LastMediItemID *uuid.UUID `json:"lastMediaItemId,omitempty" gorm:"column:last_mediaitem_id"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

const (
	JobScheduled JobStatus = "SCHEDULED"
	JobRunning   JobStatus = "RUNNING"
	JobPaused    JobStatus = "PAUSED"
	JobCompleted JobStatus = "COMPLETED"
	JobStopped   JobStatus = "STOPPED"
)

// TableName ...
func (Job) TableName() string {
	return JobsTable
}

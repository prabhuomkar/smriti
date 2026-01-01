package models

import (
	"time"

	"github.com/google/uuid"
)

const JobsTable = "jobs"

type ( // JobStatus ...
	JobStatus string
)

// Job ...
type Job struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"userId"`
	Status     JobStatus `json:"status"`
	Components []string  `json:"components"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
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

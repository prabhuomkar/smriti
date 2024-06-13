package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobsTableName(t *testing.T) {
	job := Job{}
	assert.Equal(t, JobsTable, job.TableName())
}

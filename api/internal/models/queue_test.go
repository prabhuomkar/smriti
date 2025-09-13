package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueTableName(t *testing.T) {
	queue := Queue{}
	assert.Equal(t, QueueTable, queue.TableName())
}

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThingsTableName(t *testing.T) {
	thing := Thing{}
	assert.Equal(t, ThingTable, thing.TableName())
}

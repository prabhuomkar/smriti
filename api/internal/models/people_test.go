package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeopleTableName(t *testing.T) {
	people := People{}
	assert.Equal(t, PeopleTable, people.TableName())
}

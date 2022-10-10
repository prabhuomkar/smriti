package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersTableName(t *testing.T) {
	user := User{}
	assert.Equal(t, UsersTable, user.TableName())
}

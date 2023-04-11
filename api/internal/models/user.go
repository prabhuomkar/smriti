package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const UsersTable = "users"

// User ...
type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Features  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName ...
func (User) TableName() string {
	return UsersTable
}

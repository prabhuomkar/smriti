package handlers

import (
	"api/config"

	"gorm.io/gorm"
)

// Handler ...
type Handler struct {
	Config *config.Config
	DB     *gorm.DB
	// cache
	// grpc client
}

package handlers

import (
	"api/config"

	"github.com/jmoiron/sqlx"
)

// Handler ...
type Handler struct {
	Config *config.Config
	DB     *sqlx.DB
	// cache
	// grpc client
}

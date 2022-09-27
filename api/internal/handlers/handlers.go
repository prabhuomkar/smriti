package handlers

import (
	"api/config"
	"strconv"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

// Handler ...
type Handler struct {
	Config *config.Config
	DB     *gorm.DB
	// cache
	// grpc client
}

const (
	base         = 10
	bitSize      = 32
	defaultPage  = 1
	defaultLimit = 15
)

func getOffsetAndLimit(ctx echo.Context) (int, int) {
	// page
	qpage := ctx.QueryParam("page")
	page, err := strconv.ParseUint(qpage, base, bitSize)
	if err != nil {
		page = defaultPage
	}
	// limit
	qlimit := ctx.QueryParam("limit")
	limit, err := strconv.ParseUint(qlimit, base, bitSize)
	if err != nil {
		limit = defaultLimit
	}
	return int((page - 1) * limit), int(limit)
}

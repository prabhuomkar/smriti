package handlers

import (
	"api/config"
	"api/pkg/services/worker"
	"strconv"

	"github.com/bluele/gcache"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Handler ...
type Handler struct {
	Config *config.Config
	DB     *gorm.DB
	Worker *worker.WorkerClient
	Cache  gcache.Cache
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

func getRequestingUserID(ctx echo.Context) uuid.UUID {
	userID, _ := ctx.Get("userID").(string)
	return uuid.FromStringOrNil(userID)
}

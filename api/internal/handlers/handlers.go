package handlers

import (
	"api/config"
	"api/internal/models"
	"api/pkg/cache"
	"api/pkg/services/worker"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Handler ...
type Handler struct {
	Config *config.Config
	DB     *gorm.DB
	Worker worker.WorkerClient
	Cache  cache.Provider
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

func getMonthAndDate(ctx echo.Context) (string, string, error) {
	monthDate := ctx.Param("monthDate")
	//nolint: gomnd, mnd
	if len(monthDate) == 4 { // MMDD
		return monthDate[:2], monthDate[2:], nil
	}
	return "", "", fmt.Errorf("invalid monthDate: %s", monthDate)
}

func getMediaItemFilters(ctx echo.Context) string {
	filterQuery := ""
	mediaItemType := ctx.QueryParam("type")
	if mediaItemType != "" {
		filterQuery += fmt.Sprintf(" AND mediaitem_type = '%s'", mediaItemType)
	}
	mediaItemCategory := ctx.QueryParam("category")
	if mediaItemCategory != "" {
		filterQuery += fmt.Sprintf(" AND mediaitem_category = '%s'", mediaItemCategory)
	}
	mediaItemStatus := ctx.QueryParam("status")
	if mediaItemStatus != "" && (mediaItemStatus == string(models.Unspecified) ||
		mediaItemStatus == string(models.Ready) || mediaItemStatus == string(models.Processing) ||
		mediaItemStatus == string(models.Failed)) {
		filterQuery += fmt.Sprintf(" AND status = '%s'", mediaItemStatus)
	} else {
		filterQuery += fmt.Sprintf(" AND status = '%s'", string(models.Ready))
	}
	return filterQuery
}

func getAlbumSortOrder(ctx echo.Context) string {
	if ctx.QueryParam("sort") == "name" {
		return "name asc"
	}
	return "updated_at desc"
}

func getAlbumShared(ctx echo.Context) bool {
	queryParam := ctx.QueryParam("shared")
	if queryParam == "" || strings.ToLower(queryParam) == "false" || strings.ToLower(queryParam) != "true" {
		return false
	}
	return true
}

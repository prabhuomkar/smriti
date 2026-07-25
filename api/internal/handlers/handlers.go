package handlers

import (
	"api/config"
	"api/pkg/cache"
	"api/pkg/database"
	"api/pkg/services/api"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler ...
type Handler struct {
	Config *config.Config
	DB     database.DBInterface
	Cache  cache.Provider
}

var errInvalidMonthDate = errors.New("invalid monthDate")

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

	return int((page - 1) * limit), int(limit) //nolint: gosec
}

func getRequestingUserID(ctx echo.Context) uuid.UUID {
	userIDStr, _ := ctx.Get("userID").(string)

	userID, _ := uuid.Parse(userIDStr)

	return userID
}

func getMonthAndDate(ctx echo.Context) (string, string, error) {
	monthDate := ctx.Param("monthDate")
	//nolint: mnd
	if len(monthDate) == 4 { // MMDD
		return monthDate[:2], monthDate[2:], nil
	}

	return "", "", fmt.Errorf("%w: %s", errInvalidMonthDate, monthDate)
}

func getMediaItemFilters(ctx echo.Context) string {
	filterQuery := ""
	mediaItemType := ctx.QueryParam("type")
	if mediaItemType != "" {
		if _, ok := api.MediaItemType_value[mediaItemType]; ok {
			filterQuery += fmt.Sprintf(" AND mediaitem_type = '%s'", mediaItemType)
		}
	}
	mediaItemCategory := ctx.QueryParam("category")
	if mediaItemCategory != "" {
		if _, ok := api.MediaItemCategory_value[mediaItemCategory]; ok {
			filterQuery += fmt.Sprintf(" AND mediaitem_category = '%s'", mediaItemCategory)
		}
	}
	mediaItemStatus := ctx.QueryParam("status")
	if mediaItemStatus != "" {
		if _, ok := api.MediaItemStatus_value[mediaItemStatus]; ok {
			filterQuery += fmt.Sprintf(" AND status = '%s'", mediaItemStatus)
		}
	} else {
		filterQuery += fmt.Sprintf(" AND status = '%s'", api.MediaItemStatus_READY.String())
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
	if queryParam == "" || strings.ToLower(queryParam) == "false" ||
		strings.ToLower(queryParam) != "true" {
		return false
	}

	return true
}

package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"api/config"
	"api/internal/handlers"
	"api/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	userCols  = []string{"id", "name", "username", "password", "features", "created_at", "updated_at"}
	albumCols = []string{
		"id", "user_id", "name", "description", "is_shared", "is_hidden",
		"mediaitems_count", "cover_mediaitem_id", "created_at", "updated_at",
	}
	mediaitemCols = []string{
		"id", "user_id", "filename", "hash", "description", "mime_type", "source_url", "preview_url", "thumbnail_url",
		"placeholder", "is_favourite", "is_hidden", "is_deleted", "status", "mediaitem_type", "mediaitem_category",
		"width", "height", "creation_time", "camera_make", "camera_model", "focal_length", "aperture_fnumber",
		"iso_equivalent", "exposure_time", "latitude", "longitude", "fps", "exif_data", "keywords", "created_at", "updated_at",
	}
)

func TestFeatureCheckForbidden(t *testing.T) {
	// handler
	cfg := &config.Config{Feature: config.Feature{
		Albums:     false,
		Favourites: false,
		Hidden:     false,
		Trash:      false,
		Explore:    false,
		Places:     false,
		Things:     false,
		People:     false,
		Sharing:    false,
		Jobs:       false,
	}}
	handler := &handlers.Handler{
		Config: cfg,
		DB:     nil,
	}
	featureHandlerMap := map[string]interface{}{
		"albums":     handler.GetAlbums,
		"favourites": handler.GetFavouriteMediaItems,
		"hidden":     handler.GetHiddenMediaItems,
		"trash":      handler.GetDeletedMediaItems,
		"explore":    handler.GetPlaces,
		"places":     handler.GetPlaces,
		"things":     handler.GetThings,
		"people":     handler.GetPeople,
		"jobs":       handler.GetJobs,
	}
	for feature, handler := range featureHandlerMap {
		// test
		server := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/v1/route", nil)
		rec := httptest.NewRecorder()
		checkFeature := FeatureCheck(cfg, feature)
		server.GET("/v1/route", checkFeature(handler.(func(ctx echo.Context) error)))
		server.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	}
}

func TestFeatureCheckOK(t *testing.T) {
	// handler
	cfg := &config.Config{Feature: config.Feature{
		Albums: true,
	}}
	// mock db
	// database
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()
	// handler
	handler := &handlers.Handler{
		Config: cfg,
		DB:     mockDB,
	}
	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.* FROM albums`)).
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows(append(albumCols, mediaitemCols...)))
	featureHandlerMap := map[string]interface{}{
		"albums": handler.GetAlbums,
	}
	for feature, handler := range featureHandlerMap {
		// test
		server := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/v1/route", nil)
		rec := httptest.NewRecorder()
		// context
		ctx := server.NewContext(req, rec)
		var features models.Features
		_ = json.Unmarshal([]byte(`{"albums":true}`), &features)
		ctx.Set("features", features)
		checkFeature := FeatureCheck(cfg, feature)
		err := checkFeature(handler.(func(ctx echo.Context) error))(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

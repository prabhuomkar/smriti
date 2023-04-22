package middlewares

import (
	"api/config"
	"api/internal/handlers"
	"api/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	mockGDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock",
		DriverName:           "postgres",
		Conn:                 mockDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	assert.NoError(t, err)
	// handler
	handler := &handlers.Handler{
		Config: cfg,
		DB:     mockGDB,
	}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "is_shared", "is_hidden", "cover_mediaitem_id",
			"mediaitems_count", "created_at", "updated_at"}))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "filename", "description", "mime_type", "source_url", "preview_url",
			"thumbnail_url", "is_favourite", "is_hidden", "is_deleted", "status", "mediaitem_type", "width",
			"height", "creation_time", "camera_make", "camera_model", "focal_length", "aperture_fnumber",
			"iso_equivalent", "exposure_time", "latitude", "longitude", "fps", "created_at", "updated_at"}))
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

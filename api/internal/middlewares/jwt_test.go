package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"api/config"
	"api/internal/auth"
	"api/internal/handlers"
	"api/internal/models"
	"api/pkg/cache"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bluele/gcache"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestJWTCheckUnauthorizedWithNoToken(t *testing.T) {
	// handler
	cfg := &config.Config{Feature: config.Feature{
		Albums: true,
	}}
	// mock cache
	cache := &cache.InMemoryCache{Connection: gcache.New(1024).LRU().Build()}
	handler := &handlers.Handler{
		Config: cfg,
		Cache:  cache,
	}
	checkJWT := JWTCheck(cfg, cache)

	// test
	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/route", nil)
	rec := httptest.NewRecorder()
	server.GET("/v1/route", checkJWT(handler.GetAlbums))
	server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTCheckUnauthorizedWithBadToken(t *testing.T) {
	// handler
	cfg := &config.Config{Feature: config.Feature{
		Albums: true,
	}}
	// mock cache
	cache := &cache.InMemoryCache{Connection: gcache.New(1024).LRU().Build()}
	handler := &handlers.Handler{
		Config: cfg,
		Cache:  cache,
	}
	checkJWT := JWTCheck(cfg, cache)

	// test
	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/route", nil)
	req.Header.Set("Authorization", "Bearer incorrect.token")
	rec := httptest.NewRecorder()
	server.GET("/v1/route", checkJWT(handler.GetAlbums))
	server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTCheckOK(t *testing.T) {
	// handler
	cfg := &config.Config{
		Feature: config.Feature{
			Albums: true,
		},
		Auth: config.Auth{
			AccessTTL: 60,
		},
	}
	accessToken, _ := auth.GetAccessAndRefreshTokens(cfg, models.User{ID: uuid.NewV4(), Username: "username"})
	// mock cache
	cache := &cache.InMemoryCache{Connection: gcache.New(1024).LRU().Build()}
	_ = cache.SetWithExpire(accessToken, nil, 1*time.Minute)
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
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "description", "is_shared", "is_hidden", "cover_mediaitem_id",
			"mediaitems_count", "created_at", "updated_at",
		}))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "filename", "description", "mime_type", "source_url", "preview_url",
			"thumbnail_url", "is_favourite", "is_hidden", "is_deleted", "status", "mediaitem_type", "mediaitem_category",
			"width",
			"height", "creation_time", "camera_make", "camera_model", "focal_length", "aperture_fnumber",
			"iso_equivalent", "exposure_time", "latitude", "longitude", "fps", "created_at", "updated_at",
		}))
	checkJWT := JWTCheck(cfg, cache)

	// test
	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/route", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	rec := httptest.NewRecorder()
	server.GET("/v1/route", checkJWT(handler.GetAlbums))
	server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

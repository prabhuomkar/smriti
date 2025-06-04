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

	"github.com/bluele/gcache"
	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

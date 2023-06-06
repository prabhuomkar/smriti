package middlewares

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"api/config"
	"api/internal/handlers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var sampleTime, _ = time.Parse("2006-01-02 15:04:05 -0700", "2022-09-22 11:22:33 +0530")

func TestBasicAuthCheckUnauthorizedWithNoAuth(t *testing.T) {
	// handler
	cfg := &config.Config{Admin: config.Admin{
		Username: "test",
		Password: "testT3st!",
	}}
	handler := &handlers.Handler{
		Config: cfg,
		DB:     nil,
	}
	checkBasicAuth := BasicAuthCheck(cfg)

	// test
	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/route", nil)
	rec := httptest.NewRecorder()
	server.GET("/v1/route", checkBasicAuth(handler.GetUsers))
	server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestBasicAuthCheckUnauthorizedWithBadAuth(t *testing.T) {
	// handler
	cfg := &config.Config{Admin: config.Admin{
		Username: "test",
		Password: "testT3st!",
	}}
	handler := &handlers.Handler{
		Config: cfg,
		DB:     nil,
	}
	checkBasicAuth := BasicAuthCheck(cfg)

	// test
	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/route", nil)
	req.SetBasicAuth("incorrect", "incorrect")
	rec := httptest.NewRecorder()
	server.GET("/v1/route", checkBasicAuth(handler.GetUsers))
	server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestBasicAuthCheckOK(t *testing.T) {
	// handler
	cfg := &config.Config{Admin: config.Admin{
		Username: "test",
		Password: "testT3st!",
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
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "username", "password", "created_at", "updated_at"}).
			AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "username", "password", sampleTime, sampleTime).
			AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "name", "username", "password", sampleTime, sampleTime))
	checkBasicAuth := BasicAuthCheck(cfg)

	// test
	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/route", nil)
	req.SetBasicAuth("test", "testT3st!")
	rec := httptest.NewRecorder()
	server.GET("/v1/route", checkBasicAuth(handler.GetUsers))
	server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

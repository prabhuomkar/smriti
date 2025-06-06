package middlewares

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"api/config"
	"api/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sampleTime, _ = time.Parse(
	"2006-01-02 15:04:05 -0700",
	"2022-09-22 11:22:33 +0530",
)

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
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()
	// handler
	handler := &handlers.Handler{
		Config: cfg,
		DB:     mockDB,
	}
	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users`)).
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows(userCols).
			AddRow(
				"4d05b5f6-17c2-475e-87fe-3fc8b9567179",
				"name",
				"username",
				"password",
				"",
				sampleTime,
				sampleTime,
			).
			AddRow(
				"4d05b5f6-17c2-475e-87fe-3fc8b9567180",
				"name",
				"username",
				"password",
				"",
				sampleTime,
				sampleTime,
			),
		)
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

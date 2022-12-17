package handlers

import (
	"api/config"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bluele/gcache"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Test struct {
	Name            string
	Method          string
	Route           string
	Path            string
	Header          map[string]string
	Body            io.Reader
	MockDB          func(mock sqlmock.Sqlmock)
	mockCache       []func(interface{}, interface{}) (interface{}, error)
	Handler         func(handler *Handler) func(ctx echo.Context) error
	ExpectedResCode int
	ExpectedResBody string
}

func executeTests(t *testing.T, tests []Test) {
	t.Helper()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// server
			server := echo.New()
			req := httptest.NewRequest(test.Method, test.Path, test.Body)
			for key, val := range test.Header {
				if key == echo.HeaderAuthorization {
					req.Header.Set(key, fmt.Sprintf("Bearer %s", val))
				} else {
					req.Header.Set(key, val)
				}
			}
			rec := httptest.NewRecorder()
			// context
			ctx := server.NewContext(req, rec)
			ctx.Set("userID", "4d05b5f6-17c2-475e-87fe-3fc8b9567179")
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
			if test.MockDB != nil {
				test.MockDB(mock)
			}
			mockCache := gcache.New(1024).LRU().Build()
			if test.mockCache != nil {
				mockCache = gcache.New(1024).
					LRU().
					SerializeFunc(test.mockCache[0]).
					DeserializeFunc(test.mockCache[1]).
					Build()
			}
			for key, val := range test.Header {
				if key == echo.HeaderAuthorization {
					mockCache.Set(val, true)
				}
			}
			// handler
			handler := &Handler{
				Config: &config.Config{Auth: config.Auth{RefreshTTL: 60}},
				DB:     mockGDB,
				Cache:  mockCache,
			}
			server.Match([]string{test.Method}, test.Route, test.Handler(handler))
			server.ServeHTTP(rec, req)
			// assert
			assert.Equal(t, test.ExpectedResCode, rec.Code)
			assert.Contains(t, strings.TrimSpace(rec.Body.String()), test.ExpectedResBody)
		})
	}
}

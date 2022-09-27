package handlers

import (
	"api/config"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
	Body            string
	MockDB          func(mock sqlmock.Sqlmock)
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
			req := httptest.NewRequest(test.Method, test.Path, strings.NewReader(test.Body))
			if len(test.Body) > 0 {
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}
			rec := httptest.NewRecorder()
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
				Logger: logger.Default.LogMode(logger.Info),
			})
			assert.NoError(t, err)
			if test.MockDB != nil {
				test.MockDB(mock)
			}
			// handler
			handler := &Handler{
				Config: &config.Config{},
				DB:     mockGDB,
			}
			server.Match([]string{test.Method}, test.Route, test.Handler(handler))
			server.ServeHTTP(rec, req)
			// assert
			assert.Equal(t, test.ExpectedResCode, rec.Code)
			assert.Contains(t, strings.TrimSpace(rec.Body.String()), test.ExpectedResBody)
		})
	}
}

package handlers

import (
	"api/config"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	Name            string
	Route           string
	Path            string
	Body            io.Reader
	MockDB          func(mock sqlmock.Sqlmock)
	Handler         func(handler *Handler) func(ctx echo.Context) error
	ExpectedResCode int
	ExpectedResBody string
}

func executeTests(t *testing.T, tests []Test) {
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// server
			server := echo.New()
			req := httptest.NewRequest(http.MethodGet, test.Path, test.Body)
			rec := httptest.NewRecorder()
			// database
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()
			mockDBx := sqlx.NewDb(mockDB, "sqlmock")
			if test.MockDB != nil {
				test.MockDB(mock)
			}
			// handler
			handler := &Handler{
				Config: &config.Config{},
				DB:     mockDBx,
			}
			server.GET(test.Route, test.Handler(handler))
			server.ServeHTTP(rec, req)
			// assert
			assert.Equal(t, test.ExpectedResCode, rec.Code)
			assert.Equal(t, test.ExpectedResBody, strings.TrimSpace(rec.Body.String()))
		})
	}
}

package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
)

func TestGetFeatures(t *testing.T) {
	tests := []Test{
		{
			"get features with error",
			http.MethodGet,
			"/v1/features",
			"/v1/features",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetFeatures
			},
			http.StatusOK,
			"{}",
		},
		{
			"get features successfully",
			http.MethodGet,
			"/v1/features",
			"/v1/features",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderAuthorization: "atoken",
			},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetFeatures
			},
			http.StatusOK,
			`{"albums":true,"explore":true,"places":true,"things":true,"people":true}`,
		},
	}
	executeTests(t, tests)
}

func TestGetVersion(t *testing.T) {
	tests := []Test{
		{
			"get version successfully",
			http.MethodGet,
			"/v1/version",
			"/v1/version",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetVersion
			},
			http.StatusOK,
			``,
		},
	}
	executeTests(t, tests)
}

func TestGetDisk(t *testing.T) {
	tests := []Test{
		{
			"get disk successfully",
			http.MethodGet,
			"/v1/disk",
			"/v1/disk",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetDisk
			},
			http.StatusOK,
			``,
		},
	}
	executeTests(t, tests)
}

func TestSearch(t *testing.T) {
	tests := []Test{
		{
			"search mediaitems with bad request",
			http.MethodGet,
			"/v1/search",
			"/v1/search",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Search
			},
			http.StatusBadRequest,
			"invalid search query",
		},
		{
			"search mediaitems with no results",
			http.MethodGet,
			"/v1/search",
			"/v1/search?q=keyword",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WillReturnRows(sqlmock.NewRows(mediaitemCols))
			},
			nil,
			&mockWorkerGRPCClient{wantOk: true},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Search
			},
			http.StatusOK,
			"[]",
		},
		{
			"search mediaitems with 2 rows",
			http.MethodGet,
			"/v1/search",
			"/v1/search?q=keyword",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WillReturnRows(getMockedMediaItemRows())
			},
			nil,
			&mockWorkerGRPCClient{wantOk: true},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Search
			},
			http.StatusOK,
			mediaitemsResponseBody,
		},
		{
			"search mediaitems with error",
			http.MethodGet,
			"/v1/search",
			"/v1/search?q=keyword",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			&mockWorkerGRPCClient{wantOk: true},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Search
			},
			http.StatusInternalServerError,
			"some db error",
		},
		{
			"search mediaitems with error getting embedding",
			http.MethodGet,
			"/v1/search",
			"/v1/search?q=keyword",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			nil,
			nil,
			&mockWorkerGRPCClient{wantErr: true},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Search
			},
			http.StatusInternalServerError,
			"some grpc error",
		},
	}
	executeTests(t, tests)
}

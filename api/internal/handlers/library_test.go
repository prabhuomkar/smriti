package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"
)

func TestGetFavouriteMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get favourite mediaitems with empty table",
			"/v1/favourites",
			"/v1/favourites",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems WHERE is_favourite=true`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetFavouriteMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get favourite mediaitems with 2 rows",
			"/v1/favourites",
			"/v1/favourites",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems WHERE is_favourite=true`))
				expectedQuery.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetFavouriteMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get favourite mediaitems with error",
			"/v1/favourites",
			"/v1/favourites",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems WHERE is_favourite=true`))
				expectedQuery.WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetFavouriteMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetHiddenMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get hidden mediaitems with empty table",
			"/v1/hidden",
			"/v1/hidden",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems WHERE is_hidden=true`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetHiddenMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get hidden mediaitems with 2 rows",
			"/v1/hidden",
			"/v1/hidden",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems WHERE is_hidden=true`))
				expectedQuery.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetHiddenMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get hidden mediaitems with error",
			"/v1/hidden",
			"/v1/hidden",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems WHERE is_hidden=true`))
				expectedQuery.WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetHiddenMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetDeletedMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get deleted mediaitems with empty table",
			"/v1/trash",
			"/v1/trash",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems WHERE is_deleted=true`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetDeletedMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get deleted mediaitems with 2 rows",
			"/v1/trash",
			"/v1/trash",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems WHERE is_deleted=true`))
				expectedQuery.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetDeletedMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get deleted mediaitems with error",
			"/v1/trash",
			"/v1/trash",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems WHERE is_deleted=true`))
				expectedQuery.WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetDeletedMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

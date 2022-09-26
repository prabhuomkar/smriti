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
			http.MethodGet,
			"/v1/favourites",
			"/v1/favourites",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetFavouriteMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get favourite mediaitems with 2 rows",
			http.MethodGet,
			"/v1/favourites",
			"/v1/favourites",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`))
				expectedMock.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetFavouriteMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get favourite mediaitems with error",
			http.MethodGet,
			"/v1/favourites",
			"/v1/favourites",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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
			http.MethodGet,
			"/v1/hidden",
			"/v1/hidden",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetHiddenMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get hidden mediaitems with 2 rows",
			http.MethodGet,
			"/v1/hidden",
			"/v1/hidden",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`))
				expectedMock.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetHiddenMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get hidden mediaitems with error",
			http.MethodGet,
			"/v1/hidden",
			"/v1/hidden",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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
			http.MethodGet,
			"/v1/trash",
			"/v1/trash",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetDeletedMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get deleted mediaitems with 2 rows",
			http.MethodGet,
			"/v1/trash",
			"/v1/trash",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`))
				expectedMock.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetDeletedMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get deleted mediaitems with error",
			http.MethodGet,
			"/v1/trash",
			"/v1/trash",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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

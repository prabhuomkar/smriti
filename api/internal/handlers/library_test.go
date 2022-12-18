package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
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
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(mediaitemCols))
			},
			nil,
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
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRows())
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetFavouriteMediaItems
			},
			http.StatusOK,
			mediaitemsResponseBody,
		},
		{
			"get favourite mediaitems with error",
			http.MethodGet,
			"/v1/favourites",
			"/v1/favourites",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetFavouriteMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestAddFavouriteMediaItems(t *testing.T) {
	tests := []Test{
		{
			"add favourite mediaitems with bad payload",
			http.MethodPost,
			"/v1/favourites",
			"/v1/favourites",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddFavouriteMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitems"}`,
		},
		{
			"add favourite mediaitems with bad mediaitem",
			http.MethodPost,
			"/v1/favourites",
			"/v1/favourites",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["bad-mediaitem-id"]}`),
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddFavouriteMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"add favourite mediaitems with success",
			http.MethodPost,
			"/v1/favourites",
			"/v1/favourites",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddFavouriteMediaItems
			},
			http.StatusNoContent,
			"",
		},
		{
			"add favourite mediaitems with error",
			http.MethodPost,
			"/v1/favourites",
			"/v1/favourites",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddFavouriteMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestRemoveFavouriteMediaItems(t *testing.T) {
	tests := []Test{
		{
			"remove favourite mediaitems with bad payload",
			http.MethodDelete,
			"/v1/favourites",
			"/v1/favourites",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveFavouriteMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitems"}`,
		},
		{
			"remove favourite mediaitems with bad mediaitem",
			http.MethodDelete,
			"/v1/favourites",
			"/v1/favourites",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["bad-mediaitem-id"]}`),
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveFavouriteMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"remove favourite mediaitems with success",
			http.MethodDelete,
			"/v1/favourites",
			"/v1/favourites",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveFavouriteMediaItems
			},
			http.StatusNoContent,
			"",
		},
		{
			"remove favourite mediaitems with error",
			http.MethodDelete,
			"/v1/favourites",
			"/v1/favourites",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveFavouriteMediaItems
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
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(mediaitemCols))
			},
			nil,
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
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRows())
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetHiddenMediaItems
			},
			http.StatusOK,
			mediaitemsResponseBody,
		},
		{
			"get hidden mediaitems with error",
			http.MethodGet,
			"/v1/hidden",
			"/v1/hidden",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetHiddenMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestAddHiddenMediaItems(t *testing.T) {
	tests := []Test{
		{
			"add hidden mediaitems with bad payload",
			http.MethodPost,
			"/v1/hidden",
			"/v1/hidden",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddHiddenMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitems"}`,
		},
		{
			"add hidden mediaitems with bad mediaitem",
			http.MethodPost,
			"/v1/hidden",
			"/v1/hidden",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["bad-mediaitem-id"]}`),
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddHiddenMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"add hidden mediaitems with success",
			http.MethodPost,
			"/v1/hidden",
			"/v1/hidden",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddHiddenMediaItems
			},
			http.StatusNoContent,
			"",
		},
		{
			"add hidden mediaitems with error",
			http.MethodPost,
			"/v1/hidden",
			"/v1/hidden",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddHiddenMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestRemoveHiddenMediaItems(t *testing.T) {
	tests := []Test{
		{
			"remove hidden mediaitems with bad payload",
			http.MethodDelete,
			"/v1/hidden",
			"/v1/hidden",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveHiddenMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitems"}`,
		},
		{
			"remove hidden mediaitems with bad mediaitem",
			http.MethodDelete,
			"/v1/hidden",
			"/v1/hidden",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["bad-mediaitem-id"]}`),
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveHiddenMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"remove hidden mediaitems with success",
			http.MethodDelete,
			"/v1/hidden",
			"/v1/hidden",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveHiddenMediaItems
			},
			http.StatusNoContent,
			"",
		},
		{
			"remove hidden mediaitems with error",
			http.MethodDelete,
			"/v1/hidden",
			"/v1/hidden",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveHiddenMediaItems
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
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(mediaitemCols))
			},
			nil,
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
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRows())
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetDeletedMediaItems
			},
			http.StatusOK,
			mediaitemsResponseBody,
		},
		{
			"get deleted mediaitems with error",
			http.MethodGet,
			"/v1/trash",
			"/v1/trash",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetDeletedMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestAddDeletedMediaItems(t *testing.T) {
	tests := []Test{
		{
			"add deleted mediaitems with bad payload",
			http.MethodPost,
			"/v1/trash",
			"/v1/trash",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddDeletedMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitems"}`,
		},
		{
			"add deleted mediaitems with bad mediaitem",
			http.MethodPost,
			"/v1/trash",
			"/v1/trash",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["bad-mediaitem-id"]}`),
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddDeletedMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"add deleted mediaitems with success",
			http.MethodPost,
			"/v1/trash",
			"/v1/trash",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddDeletedMediaItems
			},
			http.StatusNoContent,
			"",
		},
		{
			"add deleted mediaitems with error",
			http.MethodPost,
			"/v1/trash",
			"/v1/trash",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddDeletedMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestRemoveDeletedMediaItems(t *testing.T) {
	tests := []Test{
		{
			"remove deleted mediaitems with bad payload",
			http.MethodDelete,
			"/v1/trash",
			"/v1/trash",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveDeletedMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitems"}`,
		},
		{
			"remove deleted mediaitems with bad mediaitem",
			http.MethodDelete,
			"/v1/trash",
			"/v1/trash",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["bad-mediaitem-id"]}`),
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveDeletedMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"remove deleted mediaitems with success",
			http.MethodDelete,
			"/v1/trash",
			"/v1/trash",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveDeletedMediaItems
			},
			http.StatusNoContent,
			"",
		},
		{
			"remove deleted mediaitems with error",
			http.MethodDelete,
			"/v1/trash",
			"/v1/trash",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveDeletedMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

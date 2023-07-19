package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
)

func TestGetSharedAlbumMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get shared album mediaitems bad request",
			http.MethodGet,
			"/v1/sharing/:id/mediaItems",
			"/v1/sharing/bad-uuid/mediaItems",
			[]string{"id"},
			[]string{"bad-uuid"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetSharedAlbumMediaItems
			},
			http.StatusBadRequest,
			"invalid shared link",
		},
		{
			"get shared album mediaitems not found",
			http.MethodGet,
			"/v1/sharing/:id/mediaItems",
			"/v1/sharing/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "album_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(mediaitemCols))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetSharedAlbumMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get shared album mediaitems with 2 rows",
			http.MethodGet,
			"/v1/sharing/:id/mediaItems",
			"/v1/sharing/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "album_mediaitems"`)).
					WillReturnRows(getMockedMediaItemRows())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetSharedAlbumMediaItems
			},
			http.StatusOK,
			mediaitemsResponseBody,
		},
		{
			"get shared album mediaitems with error",
			http.MethodGet,
			"/v1/sharing/:id/mediaItems",
			"/v1/sharing/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "album_mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetSharedAlbumMediaItems
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

func TestGetSharedAlbum(t *testing.T) {
	tests := []Test{
		{
			"get shared album bad request",
			http.MethodGet,
			"/v1/sharing/:id",
			"/v1/sharing/bad-uuid",
			[]string{"id"},
			[]string{"bad-uuid"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetSharedAlbum
			},
			http.StatusBadRequest,
			"invalid shared link",
		},
		{
			"get shared album not found",
			http.MethodGet,
			"/v1/sharing/:id",
			"/v1/sharing/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`)).
					WillReturnRows(sqlmock.NewRows(albumCols))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetSharedAlbum
			},
			http.StatusNotFound,
			"shared link not found",
		},
		{
			"get shared album",
			http.MethodGet,
			"/v1/sharing/:id",
			"/v1/sharing/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`)).
					WillReturnRows(getMockedAlbumRow())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetSharedAlbum
			},
			http.StatusOK,
			albumResponseBody,
		},
		{
			"get shared album with error",
			http.MethodGet,
			"/v1/sharing/:id",
			"/v1/sharing/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetSharedAlbum
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

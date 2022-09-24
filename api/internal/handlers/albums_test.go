package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"
)

var (
	album_cols = []string{"id", "name", "description", "is_shared", "is_hidden", "cover_mediaitem_id",
		"cover_mediaitem_thumbnail_url", "mediaitems_count", "created_at", "updated_at"}
)

func TestGetAlbumMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get album mediaitems bad request",
			"/v1/albums/:id/mediaItems",
			"/v1/albums/bad-uuid/mediaItems",
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbumMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid album id"}`,
		},
		{
			"get album mediaitems not found",
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM album_mediaitems`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbumMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get album mediaitems with 2 rows",
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM album_mediaitems`))
				expectedQuery.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbumMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get album mediaitems with error",
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM album_mediaitems`))
				expectedQuery.WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbumMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestAddAlbumMediaItems(t *testing.T) {

}

func TestRemoveAlbumMediaItems(t *testing.T) {

}

func TestGetAlbum(t *testing.T) {
	tests := []Test{
		{
			"get album bad request",
			"/v1/albums/:id",
			"/v1/albums/bad-uuid",
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbum
			},
			http.StatusBadRequest,
			`{"message":"invalid album id"}`,
		},
		{
			"get album not found",
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM albums WHERE id=`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(album_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbum
			},
			http.StatusNotFound,
			`{"message":"album not found"}`,
		},
		{
			"get album",
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM albums WHERE id=`))
				expectedQuery.WillReturnRows(getMockedAlbumRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbum
			},
			http.StatusOK,
			`{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name","description":"description",` +
				`"mediaitemsCount":12,"coverMediaItemId":"cover_mediaitem_id",` +
				`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url","createdAt":"2022-09-22T11:22:33+05:30",` +
				`"updatedAt":"2022-09-22T11:22:33+05:30"}`,
		},
		{
			"get album with error",
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM albums WHERE id=`))
				expectedQuery.WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbum
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestUpdateAlbum(t *testing.T) {

}

func TestDeleteAlbum(t *testing.T) {

}

func TestGetAlbums(t *testing.T) {
	tests := []Test{
		{
			"get albums with empty table",
			"/v1/albums",
			"/v1/albums",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM albums`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(album_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbums
			},
			http.StatusOK,
			"[]",
		},
		{
			"get albums with 2 rows",
			"/v1/albums",
			"/v1/albums",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM albums`))
				expectedQuery.WillReturnRows(getMockedAlbumRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbums
			},
			http.StatusOK,
			`[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name","description":"description",` +
				`"mediaitemsCount":12,"coverMediaItemId":"cover_mediaitem_id",` +
				`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url","createdAt":"2022-09-22T11:22:33+05:30",` +
				`"updatedAt":"2022-09-22T11:22:33+05:30"},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180","name":"name",` +
				`"description":"description","mediaitemsCount":24,"coverMediaItemId":"cover_mediaitem_id",` +
				`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url","createdAt":"2022-09-22T11:22:33+05:30",` +
				`"updatedAt":"2022-09-22T11:22:33+05:30"}]`,
		},
		{
			"get albums with error",
			"/v1/albums",
			"/v1/albums",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM albums`))
				expectedQuery.WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbums
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestCreateAlbum(t *testing.T) {

}

func getMockedAlbumRows() *sqlmock.Rows {
	return sqlmock.NewRows(album_cols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "description", "true", "false", "cover_mediaitem_id",
			"cover_mediaitem_thumbnail_url", "12", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "name", "description", "false", "true", "cover_mediaitem_id",
			"cover_mediaitem_thumbnail_url", "24", sampleTime, sampleTime)
}

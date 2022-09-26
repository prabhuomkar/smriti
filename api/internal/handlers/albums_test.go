package handlers

import (
	"database/sql/driver"
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
			http.MethodGet,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/bad-uuid/mediaItems",
			"",
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbumMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid album id"}`,
		},
		{
			"get album mediaitems not found",
			http.MethodGet,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`JOIN "album_mediaitems"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbumMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get album mediaitems with 2 rows",
			http.MethodGet,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`JOIN "album_mediaitems"`))
				expectedMock.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbumMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get album mediaitems with error",
			http.MethodGet,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`JOIN "album_mediaitems"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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
	tests := []Test{
		{
			"add album mediaitems bad request",
			http.MethodPost,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/bad-uuid/mediaItems",
			"",
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid album id"}`,
		},
		{
			"add album mediaitems with bad payload",
			http.MethodPost,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			`{"bad":"request"}`,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid album mediaitems"}`,
		},
		{
			"add album mediaitems with bad mediaitem",
			http.MethodPost,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			`{"mediaItems":["bad-mediaitem-id"]}`,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"add album mediaitems with success",
			http.MethodPost,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "albums"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "album_mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			},
			http.StatusNoContent,
			"",
		},
		{
			"add album mediaitems with error",
			http.MethodPost,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "albums"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "album_mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestRemoveAlbumMediaItems(t *testing.T) {
	tests := []Test{
		{
			"remove album mediaitems bad request",
			http.MethodDelete,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/bad-uuid/mediaItems",
			"",
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid album id"}`,
		},
		{
			"remove album mediaitems with bad payload",
			http.MethodDelete,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			`{"bad":"request"}`,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid album mediaitems"}`,
		},
		{
			"remove album mediaitems with bad mediaitem",
			http.MethodDelete,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			`{"mediaItems":["bad-mediaitem-id"]}`,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"remove album mediaitems with success",
			http.MethodDelete,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "album_mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			},
			http.StatusNoContent,
			"",
		},
		{
			"remove album mediaitems with error",
			http.MethodDelete,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "album_mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetAlbum(t *testing.T) {
	tests := []Test{
		{
			"get album bad request",
			http.MethodGet,
			"/v1/albums/:id",
			"/v1/albums/bad-uuid",
			"",
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbum
			},
			http.StatusBadRequest,
			`{"message":"invalid album id"}`,
		},
		{
			"get album not found",
			http.MethodGet,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(album_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbum
			},
			http.StatusNotFound,
			`{"message":"album not found"}`,
		},
		{
			"get album",
			http.MethodGet,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`))
				expectedMock.WillReturnRows(getMockedAlbumRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbum
			},
			http.StatusOK,
			`{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name","description":"description","shared":true,` +
				`"hidden":false,"mediaItemsCount":12,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
				`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url",` +
				`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30"}`,
		},
		{
			"get album with error",
			http.MethodGet,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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
	tests := []Test{
		{
			"update album bad request",
			http.MethodPut,
			"/v1/albums/:id",
			"/v1/albums/bad-uuid",
			"",
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateAlbum
			},
			http.StatusBadRequest,
			`{"message":"invalid album id"}`,
		},
		{
			"update album with bad payload",
			http.MethodPut,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			`{"bad":"request"}`,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateAlbum
			},
			http.StatusBadRequest,
			`{"message":"invalid album"}`,
		},
		{
			"update album with success",
			http.MethodPut,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			`{"name":"name","description":"description","shared":true,"hidden":true,` +
				`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "albums"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "description", true, true,
						"4d05b5f6-17c2-475e-87fe-3fc8b9567179", sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateAlbum
			},
			http.StatusNoContent,
			"",
		},
		{
			"update album with error",
			http.MethodPut,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			`{"name":"name","description":"description","shared":true,"hidden":true,` +
				`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "albums"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "description", true, true,
						"4d05b5f6-17c2-475e-87fe-3fc8b9567179", sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateAlbum
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestDeleteAlbum(t *testing.T) {
	tests := []Test{
		{
			"delete album bad request",
			http.MethodDelete,
			"/v1/albums/:id",
			"/v1/albums/bad-uuid",
			"",
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteAlbum
			},
			http.StatusBadRequest,
			`{"message":"invalid album id"}`,
		},
		{
			"delete album with success",
			http.MethodDelete,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			`{"name":"name","description":"description","coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "albums"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteAlbum
			},
			http.StatusNoContent,
			"",
		},
		{
			"delete album with error",
			http.MethodDelete,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			`{"name":"name","description":"description","coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "albums"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteAlbum
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetAlbums(t *testing.T) {
	tests := []Test{
		{
			"get albums with empty table",
			http.MethodGet,
			"/v1/albums",
			"/v1/albums",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(album_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbums
			},
			http.StatusOK,
			"[]",
		},
		{
			"get albums with 2 rows",
			http.MethodGet,
			"/v1/albums",
			"/v1/albums",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`))
				expectedMock.WillReturnRows(getMockedAlbumRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbums
			},
			http.StatusOK,
			`[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name","description":"description","shared":true,` +
				`"hidden":false,"mediaItemsCount":12,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
				`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url","createdAt":"2022-09-22T11:22:33+05:30",` +
				`"updatedAt":"2022-09-22T11:22:33+05:30"},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180","name":"name",` +
				`"description":"description","shared":false,"hidden":true,"mediaItemsCount":24,` +
				`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
				`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url","createdAt":"2022-09-22T11:22:33+05:30",` +
				`"updatedAt":"2022-09-22T11:22:33+05:30"}]`,
		},
		{
			"get albums with error",
			http.MethodGet,
			"/v1/albums",
			"/v1/albums",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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
	tests := []Test{
		{
			"create album with bad payload",
			http.MethodPost,
			"/v1/albums",
			"/v1/albums",
			`{"bad":"request"}`,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateAlbum
			},
			http.StatusBadRequest,
			`{"message":"invalid album"}`,
		},
		{
			"create album with success",
			http.MethodPost,
			"/v1/albums",
			"/v1/albums",
			`{"name":"name","description":"description","coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "albums"`)).
					WithArgs(sqlmock.AnyArg(), "name", "description", false, false, 0,
						"4d05b5f6-17c2-475e-87fe-3fc8b9567179", nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateAlbum
			},
			http.StatusCreated,
			`"name":"name","description":"description","shared":false,` +
				`"hidden":false,"mediaItemsCount":0,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",`,
		},
		{
			"create album with error",
			http.MethodPost,
			"/v1/albums",
			"/v1/albums",
			`{"name":"name","description":"description","coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "albums"`)).
					WithArgs(sqlmock.AnyArg(), "name", "description", false, false, 0,
						"4d05b5f6-17c2-475e-87fe-3fc8b9567179", nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateAlbum
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func getMockedAlbumRows() *sqlmock.Rows {
	return sqlmock.NewRows(album_cols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "description", "true", "false", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"cover_mediaitem_thumbnail_url", "12", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "name", "description", "false", "true", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"cover_mediaitem_thumbnail_url", "24", sampleTime, sampleTime)
}

type AnyID struct{}

func (a AnyID) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

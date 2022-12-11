package handlers

import (
	"database/sql/driver"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"
)

var (
	albumCols = []string{"id", "user_id", "name", "description", "is_shared", "is_hidden", "cover_mediaitem_id",
		"mediaitems_count", "created_at", "updated_at"}
	albumResponseBody = `{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"name":"name","description":"description",` +
		`"shared":true,"hidden":false,"mediaItemsCount":12,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` +
		`"coverMediaItem":{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","favourite":true,"hidden":false,"deleted":false,"status":"status",` +
		`"mediaItemType":"mediaitem_type","width":720,"height":480,"creationTime":"2022-09-22T11:22:33+05:30",` +
		`"cameraMake":"camera_make","cameraModel":"camera_model","focalLength":"focal_length",` +
		`"apertureFnumber":"aperture_fnumber","isoEquivalent":"iso_equivalent","exposureTime":"exposure_time",` +
		`"latitude":17.580249,"longitude":-70.278493,"fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}}`
	albumsResponseBody = `[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"name":"name","description":"description",` +
		`"shared":true,"hidden":false,"mediaItemsCount":12,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` +
		`"coverMediaItem":{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","favourite":true,"hidden":false,"deleted":false,"status":"status",` +
		`"mediaItemType":"mediaitem_type","width":720,"height":480,"creationTime":"2022-09-22T11:22:33+05:30",` +
		`"cameraMake":"camera_make","cameraModel":"camera_model","focalLength":"focal_length",` +
		`"apertureFnumber":"aperture_fnumber","isoEquivalent":"iso_equivalent","exposureTime":"exposure_time",` +
		`"latitude":17.580249,"longitude":-70.278493,"fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name",` +
		`"description":"description","shared":false,"hidden":true,"mediaItemsCount":24,` +
		`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30","coverMediaItem":{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"filename":"filename","description":"description","mimeType":"mime_type","sourceUrl":"source_url",` +
		`"previewUrl":"preview_url","thumbnailUrl":"thumbnail_url","favourite":true,"hidden":false,"deleted":false,` +
		`"status":"status","mediaItemType":"mediaitem_type","width":720,"height":480,` +
		`"creationTime":"2022-09-22T11:22:33+05:30","cameraMake":"camera_make","cameraModel":"camera_model",` +
		`"focalLength":"focal_length","apertureFnumber":"aperture_fnumber","isoEquivalent":"iso_equivalent",` +
		`"exposureTime":"exposure_time","latitude":17.580249,"longitude":-70.278493,"fps":"fps",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30"}}]`
)

func TestGetAlbumMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get album mediaitems bad request",
			http.MethodGet,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/bad-uuid/mediaItems",
			map[string]string{},
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
			http.MethodGet,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "album_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(mediaitemCols))
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
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "album_mediaitems"`)).
					WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbumMediaItems
			},
			http.StatusOK,
			mediaitemsResponseBody,
		},
		{
			"get album mediaitems with error",
			http.MethodGet,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "album_mediaitems"`)).
					WillReturnError(errors.New("some db error"))
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
			map[string]string{},
			nil,
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
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitems"}`,
		},
		{
			"add album mediaitems with bad mediaitem",
			http.MethodPost,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["bad-mediaitem-id"]}`),
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
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "albums"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "album_mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "mediaitems" JOIN "album_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "albums"`)).
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
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
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
		{
			"add album mediaitems with error updating album",
			http.MethodPost,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "albums"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "album_mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "mediaitems" JOIN "album_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "albums"`)).
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
			map[string]string{},
			nil,
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
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitems"}`,
		},
		{
			"remove album mediaitems with bad mediaitem",
			http.MethodDelete,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["bad-mediaitem-id"]}`),
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
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "album_mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "mediaitems" JOIN "album_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "albums"`)).
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
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
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
		{
			"remove album mediaitems with error getting cover mediaitem",
			http.MethodDelete,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "album_mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
		{
			"remove album mediaitems with error updating album",
			http.MethodDelete,
			"/v1/albums/:id/mediaItems",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "album_mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "mediaitems" JOIN "album_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "albums"`)).
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
			map[string]string{},
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
			http.MethodGet,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`)).
					WillReturnRows(sqlmock.NewRows(albumCols))
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
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`)).
					WillReturnRows(getMockedAlbumRow())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbum
			},
			http.StatusOK,
			albumResponseBody,
		},
		{
			"get album with error",
			http.MethodGet,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`)).
					WillReturnError(errors.New("some db error"))
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
			map[string]string{},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateAlbum
			},
			http.StatusBadRequest,
			`{"message":"invalid album id"}`,
		},
		{
			"update album with no payload",
			http.MethodPut,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateAlbum
			},
			http.StatusBadRequest,
			`{"message":"invalid album"}`,
		},
		{
			"update album with bad payload",
			http.MethodPut,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateAlbum
			},
			http.StatusBadRequest,
			`{"message":"invalid album"}`,
		},
		{
			"update album with bad cover mediaitem id",
			http.MethodPut,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"name":"name","description":"description","coverMediaItemId":"bad-mediaitem-id"}`),
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateAlbum
			},
			http.StatusBadRequest,
			`{"message":"invalid album cover mediaitem id"}`,
		},
		{
			"update album with success",
			http.MethodPut,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"name":"name","description":"description","shared":true,"hidden":true,` +
				`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`),
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
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"name":"name","description":"description","shared":true,"hidden":true,` +
				`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`),
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
			map[string]string{},
			nil,
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
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "album_mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
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
			"delete album with error while clearing linked mediaitems",
			http.MethodDelete,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "album_mediaitems"`)).
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
		{
			"delete album with error",
			http.MethodDelete,
			"/v1/albums/:id",
			"/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "album_mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
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
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`)).
					WillReturnRows(sqlmock.NewRows(albumCols))
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
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`)).
					WillReturnRows(getMockedAlbumRows())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbums
			},
			http.StatusOK,
			albumsResponseBody,
		},
		{
			"get albums with error",
			http.MethodGet,
			"/v1/albums",
			"/v1/albums",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "albums"`)).
					WillReturnError(errors.New("some db error"))
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
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateAlbum
			},
			http.StatusBadRequest,
			`{"message":"invalid album"}`,
		},
		{
			"create album with no payload",
			http.MethodPost,
			"/v1/albums",
			"/v1/albums",
			map[string]string{},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateAlbum
			},
			http.StatusBadRequest,
			`{"message":"invalid album"}`,
		},
		{
			"create album with bad cover mediaitem id",
			http.MethodPost,
			"/v1/albums",
			"/v1/albums",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"name":"name","description":"description","coverMediaItemId":"bad-mediaitem-id"}`),
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateAlbum
			},
			http.StatusBadRequest,
			`{"message":"invalid album cover mediaitem id"}`,
		},
		{
			"create album with success",
			http.MethodPost,
			"/v1/albums",
			"/v1/albums",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"name":"name","description":"description","coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "albums"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "name", "description", false, false, 0,
						"4d05b5f6-17c2-475e-87fe-3fc8b9567179", sqlmock.AnyArg(), sqlmock.AnyArg()).
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
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"name":"name","description":"description","coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "albums"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "name", "description", false, false, 0,
						"4d05b5f6-17c2-475e-87fe-3fc8b9567179", sqlmock.AnyArg(), sqlmock.AnyArg()).
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

func getMockedAlbumRow() *sqlmock.Rows {
	return sqlmock.NewRows(albumCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "description",
			"true", "false", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "12", sampleTime, sampleTime)
}

func getMockedAlbumRows() *sqlmock.Rows {
	return sqlmock.NewRows(albumCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "description",
			"true", "false", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "12", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "description",
			"false", "true", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "24", sampleTime, sampleTime)
}

type AnyID struct{}

func (a AnyID) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v4"
	uuid "github.com/satori/go.uuid"
)

var (
	sampleName             = "name"
	sampleCoverMediaItemID = uuid.FromStringOrNil("4d05b5f6-17c2-475e-87fe-3fc8b9567179")
	sampleMediaItemsCount  = 12

	albumCols = []string{
		"id", "user_id", "name", "description", "is_shared", "is_hidden", "mediaitems_count", "cover_mediaitem_id", "created_at", "updated_at",
	}
	albumResponseBody = `{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"name":"name","description":"description",` +
		`"shared":true,"hidden":false,"mediaItemsCount":12,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` +
		`"coverMediaItem":{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","placeholder":"placeholder",` +
		`"mediaItemType":"mediaitem_type","mediaItemCategory":"mediaitem_category","width":720,"height":480}}`
	albumsResponseBody = `[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"name":"name","description":"description",` +
		`"shared":true,"hidden":false,"mediaItemsCount":12,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` +
		`"coverMediaItem":{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","placeholder":"placeholder","mediaItemType":"mediaitem_type",` +
		`"mediaItemCategory":"mediaitem_category","width":720,"height":480}},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name",` +
		`"description":"description","shared":false,"hidden":true,"mediaItemsCount":12,` +
		`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30","coverMediaItem":{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","placeholder":"placeholder","mediaItemType":"mediaitem_type",` +
		`"mediaItemCategory":"mediaitem_category","width":720,"height":480}}]`
)

func TestGetAlbumMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get album mediaitems bad request", http.MethodGet, "/v1/albums/:id/mediaItems", "/v1/albums/bad-uuid/mediaItems", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbumMediaItems
			}, http.StatusBadRequest, "invalid album id",
		},
		{
			"get album mediaitems not found", http.MethodGet, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(mediaitemCols))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbumMediaItems
			}, http.StatusOK, "[]",
		},
		{
			"get album mediaitems with error", http.MethodGet, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbumMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get album mediaitems with error in scanning", http.MethodGet, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(mediaitemCols).
						AddRow("invalid", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "filename", nil, &sampleDescription, "mime_type", "source_url", "preview_url", "thumbnail_url", "placeholder", &sampleBoolTrue, &sampleBoolFalse, &sampleBoolFalse, "status", "mediaitem_type", "mediaitem_category", 720, 480, sampleTime, &sampleCameraMake, &sampleCameraModel, &sampleFocalLength, &sampleApertureFnumber, &sampleIsoEquivalent, &sampleExposureTime, &sampleMegapixels, &sampleLatitude, &sampleLongitude, &sampleFPS, nil, nil, sampleTime, sampleTime))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbumMediaItems
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get album mediaitems with 2 rows", http.MethodGet, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedMediaItemRows())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbumMediaItems
			}, http.StatusOK, mediaitemsResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestAddAlbumMediaItems(t *testing.T) {
	tests := []Test{
		{
			"add album mediaitems bad request", http.MethodPost, "/v1/albums/:id/mediaItems", "/v1/albums/bad-uuid/mediaItems", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			}, http.StatusBadRequest, "invalid album id",
		},
		{
			"add album mediaitems with bad payload", http.MethodPost, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"bad":"request"}`), nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			}, http.StatusBadRequest, "invalid mediaitems",
		},
		{
			"add album mediaitems with bad mediaitem", http.MethodPost, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"mediaItems":["bad-mediaitem-id"]}`), nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			}, http.StatusBadRequest, "invalid mediaitem id",
		},
		{
			"add album mediaitems with error starting transaction", http.MethodPost, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{}).WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"add album mediaitems with error adding album mediaitems", http.MethodPost, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"add album mediaitems with error getting album mediaitem id and count", http.MethodPost, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(regexp.QuoteMeta(`SELECT mediaitem_id, COUNT(*) OVER() AS mediaitems_count FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"add album mediaitems with error updating album", http.MethodPost, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT mediaitem_id, COUNT(*) OVER() AS mediaitems_count FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"mediaitem_id", "mediaitems_count"}).
						AddRow(&sampleCoverMediaItemID, sampleMediaItemsCount))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"add album mediaitems with error committing transaction", http.MethodPost, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT mediaitem_id, COUNT(*) OVER() AS mediaitems_count FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"mediaitem_id", "mediaitems_count"}).
						AddRow(&sampleCoverMediaItemID, sampleMediaItemsCount))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit().WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"add album mediaitems with success", http.MethodPost, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT mediaitem_id, COUNT(*) OVER() AS mediaitems_count FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"mediaitem_id", "mediaitems_count"}).
						AddRow(&sampleCoverMediaItemID, sampleMediaItemsCount))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit()
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.AddAlbumMediaItems
			}, http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestRemoveAlbumMediaItems(t *testing.T) {
	tests := []Test{
		{
			"remove album mediaitems bad request", http.MethodDelete, "/v1/albums/:id/mediaItems", "/v1/albums/bad-uuid/mediaItems", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			}, http.StatusBadRequest, "invalid album id",
		},
		{
			"remove album mediaitems with bad payload", http.MethodDelete, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"bad":"request"}`), nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			}, http.StatusBadRequest, "invalid mediaitems",
		},
		{
			"remove album mediaitems with bad mediaitem", http.MethodDelete, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"mediaItems":["bad-mediaitem-id"]}`), nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			}, http.StatusBadRequest, "invalid mediaitem id",
		},
		{
			"remove album mediaitems with error starting transaction", http.MethodDelete, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{}).WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"remove album mediaitems with error removing album mediaitems", http.MethodDelete, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"remove album mediaitems with error getting album mediaitem id and count", http.MethodDelete, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT mediaitem_id, COUNT(*) OVER() AS mediaitems_count FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"remove album mediaitems with error updating album", http.MethodDelete, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT mediaitem_id, COUNT(*) OVER() AS mediaitems_count FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"mediaitem_id", "mediaitems_count"}).
						AddRow(&sampleCoverMediaItemID, sampleMediaItemsCount))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"remove album mediaitems with error committing transaction", http.MethodDelete, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT mediaitem_id, COUNT(*) OVER() AS mediaitems_count FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"mediaitem_id", "mediaitems_count"}).
						AddRow(&sampleCoverMediaItemID, sampleMediaItemsCount))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit().WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"remove album mediaitems with success", http.MethodDelete, "/v1/albums/:id/mediaItems", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"mediaItems":["4d05b5f6-17c2-475e-87fe-3fc8b9567179"]}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBeginTx(pgx.TxOptions{})
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT mediaitem_id, COUNT(*) OVER() AS mediaitems_count FROM album_mediaitems`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"mediaitem_id", "mediaitems_count"}).
						AddRow(&sampleCoverMediaItemID, sampleMediaItemsCount))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit()
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.RemoveAlbumMediaItems
			}, http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestGetAlbum(t *testing.T) {
	tests := []Test{
		{
			"get album bad request", http.MethodGet, "/v1/albums/:id", "/v1/albums/bad-uuid", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbum
			}, http.StatusBadRequest, "invalid album id",
		},
		{
			"get album not found", http.MethodGet, "/v1/albums/:id", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(albumCols, coverMediaItemCols...)))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbum
			}, http.StatusNotFound, "album not found",
		},
		{
			"get album with error", http.MethodGet, "/v1/albums/:id", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbum
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get album with error in scanning", http.MethodGet, "/v1/albums/:id", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(albumCols, coverMediaItemCols...)).
						AddRow("invalid", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", &sampleDescription, &sampleBoolTrue, &sampleBoolFalse, &sampleMediaItemsCount, &sampleCoverMediaItemID, sampleTime, sampleTime, "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", &sampleSourceURL, &samplePreviewURL, &sampleThumbnailURL, &samplePlaceholder, &sampleMediaItemType, &sampleMediaItemCategory, &sampleWidth, &sampleHeight))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbum
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get album with success", http.MethodGet, "/v1/albums/:id", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedAlbumRow())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbum
			}, http.StatusOK, albumResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestUpdateAlbum(t *testing.T) {
	tests := []Test{
		{
			"update album bad request", http.MethodPut, "/v1/albums/:id", "/v1/albums/bad-uuid", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateAlbum
			}, http.StatusBadRequest, "invalid album id",
		},
		{
			"update album with no payload", http.MethodPut, "/v1/albums/:id", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateAlbum
			}, http.StatusBadRequest, "invalid album",
		},
		{
			"update album with bad payload", http.MethodPut, "/v1/albums/:id", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"bad":"request"}`), nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateAlbum
			}, http.StatusBadRequest, "invalid album",
		},
		{
			"update album with bad cover mediaitem id", http.MethodPut, "/v1/albums/:id", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"name":"name","description":"description","coverMediaItemId":"bad-mediaitem-id"}`), nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateAlbum
			}, http.StatusBadRequest, "invalid album cover mediaitem id",
		},
		{
			"update album with error", http.MethodPut, "/v1/albums/:id", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"name":"name","description":"description","shared":true,"hidden":true,` +
				`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE albums`)).
					WithArgs(sampleCoverMediaItemID, sampleCoverMediaItemID, sampleName, &sampleDescription, &sampleBoolTrue, &sampleBoolTrue, &sampleCoverMediaItemID, pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateAlbum
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"update album with success", http.MethodPut, "/v1/albums/:id", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"name":"name","description":"description","shared":true,"hidden":true,` +
				`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE albums`)).
					WithArgs(sampleCoverMediaItemID, sampleCoverMediaItemID, sampleName, &sampleDescription, &sampleBoolTrue, &sampleBoolTrue, &sampleCoverMediaItemID, pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateAlbum
			}, http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestDeleteAlbum(t *testing.T) {
	tests := []Test{
		{
			"delete album bad request", http.MethodDelete, "/v1/albums/:id", "/v1/albums/bad-uuid", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteAlbum
			}, http.StatusBadRequest, "invalid album id",
		},
		{
			"delete album with error", http.MethodDelete, "/v1/albums/:id", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteAlbum
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"delete album with success", http.MethodDelete, "/v1/albums/:id", "/v1/albums/4d05b5f6-17c2-475e-87fe-3fc8b9567179", []string{"id"}, []string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteAlbum
			}, http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestGetAlbums(t *testing.T) {
	tests := []Test{
		{
			"get albums with empty table", http.MethodGet, "/v1/albums", "/v1/albums?sort=name&shared=true", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(albumCols, coverMediaItemCols...)))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbums
			}, http.StatusOK, "[]",
		},
		{
			"get albums with error", http.MethodGet, "/v1/albums", "/v1/albums", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbums
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get albums with error in scanning", http.MethodGet, "/v1/albums", "/v1/albums", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(albumCols, coverMediaItemCols...)).
						AddRow("invalid", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", &sampleDescription, &sampleBoolTrue, &sampleBoolFalse, &sampleMediaItemsCount, &sampleCoverMediaItemID, sampleTime, sampleTime, "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", &sampleSourceURL, &samplePreviewURL, &sampleThumbnailURL, &samplePlaceholder, &sampleMediaItemType, &sampleMediaItemCategory, &sampleWidth, &sampleHeight))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbums
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get albums with 2 rows", http.MethodGet, "/v1/albums", "/v1/albums?sort=updatedAt", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.id, m.user_id, m.source_url, m.preview_url, m.thumbnail_url, m.placeholder,`+
					` m.mediaitem_type, m.mediaitem_category, m.width, m.height FROM albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedAlbumRows())
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetAlbums
			}, http.StatusOK, albumsResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestCreateAlbum(t *testing.T) {
	tests := []Test{
		{
			"create album with bad payload", http.MethodPost, "/v1/albums", "/v1/albums", []string{}, []string{}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"bad":"request}`), nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateAlbum
			}, http.StatusBadRequest, "invalid album",
		},
		{
			"create album with no payload", http.MethodPost, "/v1/albums", "/v1/albums", []string{}, []string{}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateAlbum
			}, http.StatusBadRequest, "invalid album",
		},
		{
			"create album with bad cover mediaitem id", http.MethodPost, "/v1/albums", "/v1/albums", []string{}, []string{}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"name":"name","description":"description","coverMediaItemId":"bad-mediaitem-id"}`), nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateAlbum
			}, http.StatusBadRequest, "invalid album cover mediaitem id",
		},
		{
			"create album with error", http.MethodPost, "/v1/albums", "/v1/albums", []string{}, []string{}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"name":"name","description":"description"}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), "name", &sampleDescription, pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateAlbum
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"create album with success", http.MethodPost, "/v1/albums", "/v1/albums", []string{}, []string{}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"name":"name","description":"description","coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO albums`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), "name", &sampleDescription, pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			}, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateAlbum
			}, http.StatusCreated, `"name":"name","description":"description",`,
		},
	}
	executeTests(t, tests)
}

func getMockedAlbumRow() *pgxmock.Rows {
	return pgxmock.NewRows(append(albumCols, coverMediaItemCols...)).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", &sampleDescription, &sampleBoolTrue, &sampleBoolFalse, &sampleMediaItemsCount, &sampleCoverMediaItemID, sampleTime, sampleTime, &sampleCoverMediaItemID, &sampleCoverMediaItemID, &sampleSourceURL, &samplePreviewURL, &sampleThumbnailURL, &samplePlaceholder, &sampleMediaItemType, &sampleMediaItemCategory, &sampleWidth, &sampleHeight)
}

func getMockedAlbumRows() *pgxmock.Rows {
	return pgxmock.NewRows(append(albumCols, coverMediaItemCols...)).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", &sampleDescription, &sampleBoolTrue, &sampleBoolFalse, &sampleMediaItemsCount, &sampleCoverMediaItemID, sampleTime, sampleTime, &sampleCoverMediaItemID, &sampleCoverMediaItemID, &sampleSourceURL, &samplePreviewURL, &sampleThumbnailURL, &samplePlaceholder, &sampleMediaItemType, &sampleMediaItemCategory, &sampleWidth, &sampleHeight).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", &sampleDescription, &sampleBoolFalse, &sampleBoolTrue, &sampleMediaItemsCount, &sampleCoverMediaItemID, sampleTime, sampleTime, &sampleCoverMediaItemID, &sampleCoverMediaItemID, &sampleSourceURL, &samplePreviewURL, &sampleThumbnailURL, &samplePlaceholder, &sampleMediaItemType, &sampleMediaItemCategory, &sampleWidth, &sampleHeight)
}

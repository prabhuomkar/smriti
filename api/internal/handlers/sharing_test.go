package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v4"
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
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(mediaitemCols))
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
			"get shared album mediaitems with error",
			http.MethodGet,
			"/v1/sharing/:id/mediaItems",
			"/v1/sharing/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
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
		{
			"get shared album mediaitems with error in scanning",
			http.MethodGet,
			"/v1/sharing/:id/mediaItems",
			"/v1/sharing/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(mediaitemCols).AddRow("invalid", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
						"filename", nil, &sampleDescription, "mime_type", "source_url", "preview_url",
						"thumbnail_url", "placeholder", &sampleBoolTrue, &sampleBoolFalse, &sampleBoolFalse, "status", "mediaitem_type", "mediaitem_category", 720,
						480, sampleTime, &sampleCameraMake, &sampleCameraModel, &sampleFocalLength, &sampleApertureFnumber,
						&sampleIsoEquivalent, &sampleExposureTime, &sampleLatitude, &sampleLongitude, &sampleFPS, nil, nil, sampleTime, sampleTime))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetSharedAlbumMediaItems
			},
			http.StatusInternalServerError,
			"Scanning value error",
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
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
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
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.* FROM albums`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(albumCols))
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
			"get shared album with error",
			http.MethodGet,
			"/v1/sharing/:id",
			"/v1/sharing/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.* FROM albums`)).
					WithArgs(pgxmock.AnyArg()).
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
		{
			"get shared album with error in scanning",
			http.MethodGet,
			"/v1/sharing/:id",
			"/v1/sharing/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.* FROM albums`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(append(albumCols, mediaitemCols...)).AddRow("invalid", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", &sampleDescription,
						&sampleBoolTrue, &sampleBoolFalse, &sampleMediaItemsCount, &sampleCoverMediaItemID, sampleTime, sampleTime, "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
						"filename", nil, &sampleDescription, "mime_type", "source_url", "preview_url",
						"thumbnail_url", "placeholder", &sampleBoolTrue, &sampleBoolFalse, &sampleBoolFalse, "status", "mediaitem_type", "mediaitem_category", 720,
						480, sampleTime, &sampleCameraMake, &sampleCameraModel, &sampleFocalLength, &sampleApertureFnumber,
						&sampleIsoEquivalent, &sampleExposureTime, &sampleLatitude, &sampleLongitude, &sampleFPS, nil, nil, sampleTime, sampleTime))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetSharedAlbum
			},
			http.StatusInternalServerError,
			"Scanning value error",
		},
		{
			"get shared album with success",
			http.MethodGet,
			"/v1/sharing/:id",
			"/v1/sharing/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.*, m.* FROM albums`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnRows(getMockedAlbumRow())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetSharedAlbum
			},
			http.StatusOK,
			albumResponseBody,
		},
	}
	executeTests(t, tests)
}

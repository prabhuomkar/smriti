package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v4"
)

func TestGetFeatures(t *testing.T) {
	tests := []Test{
		{
			"get features with error", http.MethodGet, "/v1/features", "/v1/features", []string{}, []string{}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetFeatures
			}, http.StatusOK, "{}",
		},
		{
			"get features successfully", http.MethodGet, "/v1/features", "/v1/features", []string{}, []string{}, map[string]string{
				echo.HeaderAuthorization: "atoken",
			}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetFeatures
			}, http.StatusOK, `{"albums":true,"explore":true,"places":true,"things":true,"people":true}`,
		},
	}
	executeTests(t, tests)
}

func TestGetVersion(t *testing.T) {
	tests := []Test{
		{
			"get version successfully", http.MethodGet, "/version", "/version", []string{}, []string{}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetVersion
			}, http.StatusOK, ``,
		},
	}
	executeTests(t, tests)
}

func TestGetDisk(t *testing.T) {
	tests := []Test{
		{
			"get disk successfully", http.MethodGet, "/disk", "/disk", []string{}, []string{}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetDisk
			}, http.StatusOK, ``,
		},
	}
	executeTests(t, tests)
}

func TestSearch(t *testing.T) {
	tests := []Test{
		{
			"search mediaitems with bad request", http.MethodGet, "/v1/search", "/v1/search", []string{}, []string{}, map[string]string{}, nil, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.Search
			}, http.StatusBadRequest, "invalid search query",
		},
		{
			"search mediaitems with no results", http.MethodGet, "/v1/search", "/v1/search?q=keyword", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(mediaitemCols))
			}, nil, &mockWorkerGRPCClient{wantOk: true}, func(handler *Handler) func(ctx echo.Context) error {
				return handler.Search
			}, http.StatusOK, "[]",
		},
		{
			"search mediaitems with error getting embedding", http.MethodGet, "/v1/search", "/v1/search?q=keyword", []string{}, []string{}, map[string]string{}, nil, nil, nil, &mockWorkerGRPCClient{wantErr: true}, func(handler *Handler) func(ctx echo.Context) error {
				return handler.Search
			}, http.StatusInternalServerError, "some grpc error",
		},
		{
			"search mediaitems with error", http.MethodGet, "/v1/search", "/v1/search?q=keyword", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, &mockWorkerGRPCClient{wantOk: true}, func(handler *Handler) func(ctx echo.Context) error {
				return handler.Search
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"search mediaitems with error in scanning", http.MethodGet, "/v1/search", "/v1/search?q=keyword", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(mediaitemCols).AddRow("invalid", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "filename", nil, &sampleDescription, "mime_type", "source_url", "preview_url", "thumbnail_url", "placeholder", &sampleBoolTrue, &sampleBoolFalse, &sampleBoolFalse, "status", "mediaitem_type", "mediaitem_category", 720, 480, sampleTime, &sampleCameraMake, &sampleCameraModel, &sampleFocalLength, &sampleApertureFnumber, &sampleIsoEquivalent, &sampleExposureTime, &sampleMegapixels, &sampleLatitude, &sampleLongitude, &sampleFPS, nil, nil, sampleTime, sampleTime))
			}, nil, &mockWorkerGRPCClient{wantOk: true}, func(handler *Handler) func(ctx echo.Context) error {
				return handler.Search
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"search mediaitems with 2 rows", http.MethodGet, "/v1/search", "/v1/search?q=keyword", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(getMockedMediaItemRows())
			}, nil, &mockWorkerGRPCClient{wantOk: true}, func(handler *Handler) func(ctx echo.Context) error {
				return handler.Search
			}, http.StatusOK, mediaitemsResponseBody,
		},
	}
	executeTests(t, tests)
}

package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"
)

var (
	sampleTime, _  = time.Parse("2006-01-02 15:04:05 -0700", "2022-09-22 11:22:33 +0530")
	mediaitem_cols = []string{"id", "filename", "description", "mime_type", "source_url", "preview_url",
		"thumbnail_url", "is_favourite", "is_hidden", "is_deleted", "status", "mediaitem_type", "width",
		"height", "creation_time", "camera_make", "camera_model", "focal_length", "aperture_fnumber",
		"iso_equivalent", "exposure_time", "location", "fps", "created_at", "updated_at"}
	mediaitem_response_body = `{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","filename":"filename","description":"description",` +
		`"mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url","thumbnailUrl":"thumbnail_url",` +
		`"status":"status","mediaItemType":"mediaitem_type","width":720,"height":480,` +
		`"creationTime":"2022-09-22T11:22:33+05:30","cameraMake":"camera_make","cameraModel":"camera_model",` +
		`"focalLength":"focal_length","apertureFnumber":"aperture_fnumber","isoEquivalent":"iso_equivalent",` +
		`"exposureTime":"exposure_time","location":"bG9jYXRpb24=","fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}`
	mediaitems_response_body = `[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","filename":"filename","description":"description",` +
		`"mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url","thumbnailUrl":"thumbnail_url",` +
		`"status":"status","mediaItemType":"mediaitem_type","width":720,"height":480,` +
		`"creationTime":"2022-09-22T11:22:33+05:30","cameraMake":"camera_make","cameraModel":"camera_model",` +
		`"focalLength":"focal_length","apertureFnumber":"aperture_fnumber","isoEquivalent":"iso_equivalent",` +
		`"exposureTime":"exposure_time","location":"bG9jYXRpb24=","fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","status":"status","mediaItemType":"mediaitem_type","width":720,"height":480,` +
		`"creationTime":"2022-09-22T11:22:33+05:30","cameraMake":"camera_make","cameraModel":"camera_model",` +
		`"focalLength":"focal_length","apertureFnumber":"aperture_fnumber","isoEquivalent":"iso_equivalent",` +
		`"exposureTime":"exposure_time","location":"bG9jYXRpb24=","fps":"fps",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30"}]`
)

func TestGetMediaItemPlaces(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem mediaitem bad request",
			"/v1/mediaItems/:id/places",
			"/v1/mediaItems/bad-uuid/places",
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPlaces
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"get mediaitem places with empty table",
			"/v1/mediaItems/:id/places",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/places",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM places INNER JOIN place_mediaitems`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(place_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPlaces
			},
			http.StatusOK,
			"[]",
		},
		{
			"get mediaitem places with 2 rows",
			"/v1/mediaItems/:id/places",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/places",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM places INNER JOIN place_mediaitems`))
				expectedQuery.WillReturnRows(getMockedPlaceRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPlaces
			},
			http.StatusOK,
			places_response_body,
		},
		{
			"get mediaitem places with error",
			"/v1/mediaItems/:id/places",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/places",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM places INNER JOIN place_mediaitems`))
				expectedQuery.WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPlaces
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetMediaItemThings(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem mediaitem bad request",
			"/v1/mediaItems/:id/things",
			"/v1/mediaItems/bad-uuid/things",
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemThings
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"get mediaitem things with empty table",
			"/v1/mediaItems/:id/things",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/things",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM things INNER JOIN thing_mediaitems`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(thing_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemThings
			},
			http.StatusOK,
			"[]",
		},
		{
			"get mediaitem things with 2 rows",
			"/v1/mediaItems/:id/things",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/things",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM things INNER JOIN thing_mediaitems`))
				expectedQuery.WillReturnRows(getMockedThingRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemThings
			},
			http.StatusOK,
			things_response_body,
		},
		{
			"get mediaitem things with error",
			"/v1/mediaItems/:id/things",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/things",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM things INNER JOIN thing_mediaitems`))
				expectedQuery.WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemThings
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetMediaItemPeople(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem mediaitem bad request",
			"/v1/mediaItems/:id/people",
			"/v1/mediaItems/bad-uuid/people",
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPeople
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"get mediaitem people with empty table",
			"/v1/mediaItems/:id/people",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/people",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM people INNER JOIN people_mediaitems`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(people_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPeople
			},
			http.StatusOK,
			"[]",
		},
		{
			"get mediaitem people with 2 rows",
			"/v1/mediaItems/:id/people",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/people",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM people INNER JOIN people_mediaitems`))
				expectedQuery.WillReturnRows(getMockedPeopleRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPeople
			},
			http.StatusOK,
			people_response_body,
		},
		{
			"get mediaitem people with error",
			"/v1/mediaItems/:id/people",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/people",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM people INNER JOIN people_mediaitems`))
				expectedQuery.WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPeople
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetMediaItem(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem bad request",
			"/v1/mediaItems/:id",
			"/v1/mediaItems/bad-uuid",
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"get mediaitem not found",
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems WHERE id=`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			},
			http.StatusNotFound,
			`{"message":"mediaitem not found"}`,
		},
		{
			"get mediaitem",
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems WHERE id=`))
				expectedQuery.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			},
			http.StatusOK,
			mediaitem_response_body,
		},
		{
			"get mediaitem with error",
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems WHERE id=`))
				expectedQuery.WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestUpdateMediaItem(t *testing.T) {

}

func TestDeleteMediaItem(t *testing.T) {

}

func TestGetMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get mediaitems with empty table",
			"/v1/mediaItems",
			"/v1/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get mediaitems with 2 rows",
			"/v1/mediaItems",
			"/v1/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`))
				expectedQuery.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get mediaitems with error",
			"/v1/mediaItems",
			"/v1/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM mediaitems`))
				expectedQuery.WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestUploadMediaItems(t *testing.T) {

}

func getMockedMediaItemRows() *sqlmock.Rows {
	return sqlmock.NewRows(mediaitem_cols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "filename", "description", "mime_type", "source_url", "preview_url",
			"thumbnail_url", "true", "false", "false", "status", "mediaitem_type", 720,
			480, sampleTime, "camera_make", "camera_model", "focal_length", "aperture_fnumber",
			"iso_equivalent", "exposure_time", "location", "fps", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "filename", "description", "mime_type", "source_url", "preview_url",
			"thumbnail_url", "false", "true", "true", "status", "mediaitem_type", 720,
			480, sampleTime, "camera_make", "camera_model", "focal_length", "aperture_fnumber",
			"iso_equivalent", "exposure_time", "location", "fps", sampleTime, sampleTime)
}
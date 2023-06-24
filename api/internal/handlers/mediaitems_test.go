package handlers

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
)

var (
	sampleTime, _ = time.Parse("2006-01-02 15:04:05 -0700", "2022-09-22 11:22:33 +0530")
	mediaitemCols = []string{
		"id", "user_id", "filename", "description", "mime_type", "source_url", "preview_url",
		"thumbnail_url", "is_favourite", "is_hidden", "is_deleted", "status", "mediaitem_type", "mediaitem_category",
		"width", "height", "creation_time", "camera_make", "camera_model", "focal_length", "aperture_fnumber",
		"iso_equivalent", "exposure_time", "latitude", "longitude", "fps", "created_at", "updated_at",
	}
	mediaitemResponseBody = `{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","favourite":true,"hidden":false,"deleted":false,"status":"status",` +
		`"mediaItemType":"mediaitem_type","mediaItemCategory":"mediaitem_category","width":720,"height":480,"creationTime":"2022-09-22T11:22:33+05:30",` +
		`"cameraMake":"camera_make","cameraModel":"camera_model","focalLength":"focal_length",` +
		`"apertureFNumber":"aperture_fnumber","isoEquivalent":"iso_equivalent","exposureTime":"exposure_time",` +
		`"latitude":17.580249,"longitude":-70.278493,"fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}`
	mediaitemsResponseBody = `[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","favourite":true,"hidden":false,"deleted":false,"status":"status",` +
		`"mediaItemType":"mediaitem_type","mediaItemCategory":"mediaitem_category","width":720,"height":480,"creationTime":"2022-09-22T11:22:33+05:30",` +
		`"cameraMake":"camera_make","cameraModel":"camera_model","focalLength":"focal_length",` +
		`"apertureFNumber":"aperture_fnumber","isoEquivalent":"iso_equivalent","exposureTime":"exposure_time",` +
		`"latitude":17.580249,"longitude":-70.278493,"fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","favourite":false,"hidden":true,"deleted":true,"status":"status",` +
		`"mediaItemType":"mediaitem_type","mediaItemCategory":"mediaitem_category","width":720,"height":480,"creationTime":"2022-09-22T11:22:33+05:30",` +
		`"cameraMake":"camera_make","cameraModel":"camera_model","focalLength":"focal_length",` +
		`"apertureFNumber":"aperture_fnumber","isoEquivalent":"iso_equivalent","exposureTime":"exposure_time",` +
		`"latitude":17.580249,"longitude":-70.278493,"fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}]`
)

func TestGetMediaItemPlaces(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem places bad request",
			http.MethodGet,
			"/v1/mediaItems/:id/places",
			"/v1/mediaItems/bad-uuid/places",
			[]string{"id"},
			[]string{"bad-uuid"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPlaces
			},
			http.StatusBadRequest,
			"invalid mediaitem id",
		},
		{
			"get mediaitem places with empty table",
			http.MethodGet,
			"/v1/mediaItems/:id/places",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/places",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "place_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(placeCols))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPlaces
			},
			http.StatusOK,
			"[]",
		},
		{
			"get mediaitem places with success",
			http.MethodGet,
			"/v1/mediaItems/:id/places",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/places",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "place_mediaitems"`)).
					WillReturnRows(getMockedPlaceRows())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPlaces
			},
			http.StatusOK,
			placesResponseBody,
		},
		{
			"get mediaitem places with error",
			http.MethodGet,
			"/v1/mediaItems/:id/places",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/places",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "place_mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPlaces
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

func TestGetMediaItemThings(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem things bad request",
			http.MethodGet,
			"/v1/mediaItems/:id/things",
			"/v1/mediaItems/bad-uuid/things",
			[]string{"id"},
			[]string{"bad-uuid"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemThings
			},
			http.StatusBadRequest,
			"invalid mediaitem id",
		},
		{
			"get mediaitem things with empty table",
			http.MethodGet,
			"/v1/mediaItems/:id/things",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/things",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "thing_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(thingCols))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemThings
			},
			http.StatusOK,
			"[]",
		},
		{
			"get mediaitem things with success",
			http.MethodGet,
			"/v1/mediaItems/:id/things",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/things",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "thing_mediaitems"`)).
					WillReturnRows(getMockedThingRows())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemThings
			},
			http.StatusOK,
			thingsResponseBody,
		},
		{
			"get mediaitem things with error",
			http.MethodGet,
			"/v1/mediaItems/:id/things",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/things",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "thing_mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemThings
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

func TestGetMediaItemPeople(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem people bad request",
			http.MethodGet,
			"/v1/mediaItems/:id/people",
			"/v1/mediaItems/bad-uuid/people",
			[]string{"id"},
			[]string{"bad-uuid"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPeople
			},
			http.StatusBadRequest,
			"invalid mediaitem id",
		},
		{
			"get mediaitem people with empty table",
			http.MethodGet,
			"/v1/mediaItems/:id/people",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/people",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "people_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(peopleCols))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPeople
			},
			http.StatusOK,
			"[]",
		},
		{
			"get mediaitem people with success",
			http.MethodGet,
			"/v1/mediaItems/:id/people",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/people",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "people_mediaitems"`)).
					WillReturnRows(getMockedPeopleRows())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPeople
			},
			http.StatusOK,
			peopleResponseBody,
		},
		{
			"get mediaitem people with error",
			http.MethodGet,
			"/v1/mediaItems/:id/people",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/people",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "people_mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPeople
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

func TestGetMediaItemAlbums(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem albums bad request",
			http.MethodGet,
			"/v1/mediaItems/:id/albums",
			"/v1/mediaItems/bad-uuid/albums",
			[]string{"id"},
			[]string{"bad-uuid"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemAlbums
			},
			http.StatusBadRequest,
			"invalid mediaitem id",
		},
		{
			"get mediaitem albums with empty table",
			http.MethodGet,
			"/v1/mediaItems/:id/albums",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/albums",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "album_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(albumCols))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemAlbums
			},
			http.StatusOK,
			"[]",
		},
		{
			"get mediaitem albums with success",
			http.MethodGet,
			"/v1/mediaItems/:id/albums",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/albums",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "album_mediaitems"`)).
					WillReturnRows(getMockedAlbumRows())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemAlbums
			},
			http.StatusOK,
			albumsResponseBody,
		},
		{
			"get mediaitem albums with error",
			http.MethodGet,
			"/v1/mediaItems/:id/albums",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/albums",
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
				return handler.GetMediaItemAlbums
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

func TestGetMediaItem(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem bad request",
			http.MethodGet,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/bad-uuid",
			[]string{"id"},
			[]string{"bad-uuid"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			},
			http.StatusBadRequest,
			"invalid mediaitem id",
		},
		{
			"get mediaitem not found",
			http.MethodGet,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(mediaitemCols))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			},
			http.StatusNotFound,
			"mediaitem not found",
		},
		{
			"get mediaitem",
			http.MethodGet,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRows())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			},
			http.StatusOK,
			mediaitemResponseBody,
		},
		{
			"get mediaitem with error",
			http.MethodGet,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

func TestUpdateMediaItem(t *testing.T) {
	tests := []Test{
		{
			"update mediaitem bad request",
			http.MethodPut,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/bad-uuid",
			[]string{"id"},
			[]string{"bad-uuid"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateMediaItem
			},
			http.StatusBadRequest,
			"invalid mediaitem id",
		},
		{
			"update mediaitem with no payload",
			http.MethodPut,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateMediaItem
			},
			http.StatusBadRequest,
			"invalid mediaitem",
		},
		{
			"update mediaitem with bad payload",
			http.MethodPut,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateMediaItem
			},
			http.StatusBadRequest,
			"invalid mediaitem",
		},
		{
			"update mediaitem with success",
			http.MethodPut,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"description":"description","favourite":true,"hidden":true}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "description", true, true,
						sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateMediaItem
			},
			http.StatusNoContent,
			"",
		},
		{
			"update mediaitem with error",
			http.MethodPut,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"description":"description","favourite":true,"hidden":true}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "description", true, true,
						sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateMediaItem
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

func TestDeleteMediaItem(t *testing.T) {
	tests := []Test{
		{
			"delete mediaitem bad request",
			http.MethodDelete,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/bad-uuid",
			[]string{"id"},
			[]string{"bad-uuid"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			},
			http.StatusBadRequest,
			"invalid mediaitem id",
		},
		{
			"delete mediaitem with success",
			http.MethodDelete,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", true,
						sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			},
			http.StatusNoContent,
			"",
		},
		{
			"delete mediaitem with error",
			http.MethodDelete,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", true,
						sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

func TestGetMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get mediaitems with empty table",
			http.MethodGet,
			"/v1/mediaItems",
			"/v1/mediaItems",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(mediaitemCols))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get mediaitems with 2 rows",
			http.MethodGet,
			"/v1/mediaItems",
			"/v1/mediaItems?type=photo&category=panorama",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRows())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItems
			},
			http.StatusOK,
			mediaitemsResponseBody,
		},
		{
			"get mediaitems with error",
			http.MethodGet,
			"/v1/mediaItems",
			"/v1/mediaItems",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItems
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

func TestUploadMediaItems(t *testing.T) {
	sampleFile, contentType := getMockedMediaItemFile(t)
	sampleFile2, contentType2 := getMockedMediaItemFile(t)
	sampleFile3, contentType3 := getMockedMediaItemFile(t)
	sampleFile4, contentType4 := getMockedMediaItemFile(t)
	sampleFile5, contentType5 := getMockedMediaItemFile(t)
	tests := []Test{
		{
			"upload mediaitems with invalid command",
			http.MethodPost,
			"/v1/mediaItems",
			"/v1/mediaItems",
			[]string{},
			[]string{},
			map[string]string{
				HeaderUploadType: "resumable",
			},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			},
			http.StatusBadRequest,
			"invalid command for resumable upload",
		},
		{
			"upload mediaitems with invalid offset",
			http.MethodPost,
			"/v1/mediaItems",
			"/v1/mediaItems",
			[]string{},
			[]string{},
			map[string]string{
				HeaderUploadType:    "resumable",
				HeaderUploadCommand: "finish",
			},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			},
			http.StatusBadRequest,
			"invalid chunk offset for resumable upload",
		},
		{
			"upload mediaitems with invalid session",
			http.MethodPost,
			"/v1/mediaItems",
			"/v1/mediaItems",
			[]string{},
			[]string{},
			map[string]string{
				HeaderUploadType:        "resumable",
				HeaderUploadCommand:     "finish",
				HeaderUploadChunkOffset: "1024",
			},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			},
			http.StatusBadRequest,
			"invalid chunk session for resumable upload",
		},
		{
			"upload mediaitems with error uploading for resumable",
			http.MethodPost,
			"/v1/mediaItems",
			"/v1/mediaItems",
			[]string{},
			[]string{},
			map[string]string{
				HeaderUploadType:        "resumable",
				HeaderUploadCommand:     "start",
				HeaderUploadChunkOffset: "0",
			},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			},
			http.StatusBadRequest,
			"request Content-Type isn't multipart/form-data",
		},
		{
			"upload mediaitems with error uploading",
			http.MethodPost,
			"/v1/mediaItems",
			"/v1/mediaItems",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			},
			http.StatusBadRequest,
			"request Content-Type isn't multipart/form-data",
		},
		{
			"upload mediaitems with error inserting mediaitem",
			http.MethodPost,
			"/v1/mediaItems",
			"/v1/mediaItems",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: contentType,
			},
			sampleFile,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			},
			http.StatusInternalServerError,
			"some db error",
		},
		{
			"upload mediaitems with error sending file to worker due to error in mediaitem process",
			http.MethodPost,
			"/v1/mediaItems",
			"/v1/mediaItems",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: contentType2,
			},
			sampleFile2,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			&mockWorkerGRPCClient{wantErr: true},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			},
			http.StatusInternalServerError,
			"some grpc error",
		},
		{
			"upload mediaitems successfully",
			http.MethodPost,
			"/v1/mediaItems",
			"/v1/mediaItems",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: contentType3,
			},
			sampleFile3,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			&mockWorkerGRPCClient{wantOk: true},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			},
			http.StatusCreated,
			`"id"`,
		},
		{
			"upload mediaitems with error for resumable",
			http.MethodPost,
			"/v1/mediaItems",
			"/v1/mediaItems",
			[]string{},
			[]string{},
			map[string]string{
				HeaderUploadType:         "resumable",
				HeaderUploadCommand:      "finish",
				HeaderUploadChunkOffset:  "100",
				HeaderUploadChunkSession: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
				echo.HeaderContentType:   contentType4,
			},
			sampleFile4,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			&mockWorkerGRPCClient{wantErr: true},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			},
			http.StatusInternalServerError,
			"some grpc error",
		},
		{
			"upload mediaitems successfully for resumable",
			http.MethodPost,
			"/v1/mediaItems",
			"/v1/mediaItems",
			[]string{},
			[]string{},
			map[string]string{
				HeaderUploadType:         "resumable",
				HeaderUploadCommand:      "finish",
				HeaderUploadChunkOffset:  "100",
				HeaderUploadChunkSession: "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
				echo.HeaderContentType:   contentType5,
			},
			sampleFile5,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "mediaitems"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			&mockWorkerGRPCClient{wantOk: true},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UploadMediaItems
			},
			http.StatusNoContent,
			``,
		},
	}
	executeTests(t, tests)
}

func getMockedMediaItemRow() *sqlmock.Rows {
	return sqlmock.NewRows(mediaitemCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"filename", "description", "mime_type", "source_url", "preview_url",
			"thumbnail_url", "true", "false", "false", "status", "mediaitem_type", "mediaitem_category", 720,
			480, sampleTime, "camera_make", "camera_model", "focal_length", "aperture_fnumber",
			"iso_equivalent", "exposure_time", "17.580249", "-70.278493", "fps", sampleTime, sampleTime)
}

func getMockedMediaItemRows() *sqlmock.Rows {
	return sqlmock.NewRows(mediaitemCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"filename", "description", "mime_type", "source_url", "preview_url",
			"thumbnail_url", "true", "false", "false", "status", "mediaitem_type", "mediaitem_category", 720,
			480, sampleTime, "camera_make", "camera_model", "focal_length", "aperture_fnumber",
			"iso_equivalent", "exposure_time", "17.580249", "-70.278493", "fps", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"filename", "description", "mime_type", "source_url", "preview_url",
			"thumbnail_url", "false", "true", "true", "status", "mediaitem_type", "mediaitem_category", 720,
			480, sampleTime, "camera_make", "camera_model", "focal_length", "aperture_fnumber",
			"iso_equivalent", "exposure_time", "17.580249", "-70.278493", "fps", sampleTime, sampleTime)
}

func getMockedMediaItemFile(t *testing.T) (*io.PipeReader, string) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		part, err := writer.CreateFormFile("file", "image.png")
		if err != nil {
			t.Error(err)
		}

		// create sample image
		dim := 10
		upLeft := image.Point{0, 0}
		lowRight := image.Point{dim, dim}
		img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
		for x := 0; x < dim; x++ {
			for y := 0; y < dim; y++ {
				img.Set(x, y, color.White)
			}
		}

		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	return pr, writer.FormDataContentType()
}

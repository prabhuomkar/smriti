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
	sampleTime, _ = time.Parse("2006-01-02 15:04:05 -0700", "2022-09-22 11:22:33 +0530")
	mediaitemCols = []string{"id", "filename", "description", "mime_type", "source_url", "preview_url",
		"thumbnail_url", "is_favourite", "is_hidden", "is_deleted", "status", "mediaitem_type", "width",
		"height", "creation_time", "camera_make", "camera_model", "focal_length", "aperture_fnumber",
		"iso_equivalent", "exposure_time", "latitude", "longitude", "fps", "created_at", "updated_at"}
	mediaitemResponseBody = `{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","favourite":true,"hidden":false,"deleted":false,"status":"status",` +
		`"mediaItemType":"mediaitem_type","width":720,"height":480,"creationTime":"2022-09-22T11:22:33+05:30",` +
		`"cameraMake":"camera_make","cameraModel":"camera_model","focalLength":"focal_length",` +
		`"apertureFnumber":"aperture_fnumber","isoEquivalent":"iso_equivalent","exposureTime":"exposure_time",` +
		`"latitude":17.580249,"longitude":-70.278493,"fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}`
	mediaitemsResponseBody = `[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","favourite":true,"hidden":false,"deleted":false,"status":"status",` +
		`"mediaItemType":"mediaitem_type","width":720,"height":480,"creationTime":"2022-09-22T11:22:33+05:30",` +
		`"cameraMake":"camera_make","cameraModel":"camera_model","focalLength":"focal_length",` +
		`"apertureFnumber":"aperture_fnumber","isoEquivalent":"iso_equivalent","exposureTime":"exposure_time",` +
		`"latitude":17.580249,"longitude":-70.278493,"fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","favourite":false,"hidden":true,"deleted":true,"status":"status",` +
		`"mediaItemType":"mediaitem_type","width":720,"height":480,"creationTime":"2022-09-22T11:22:33+05:30",` +
		`"cameraMake":"camera_make","cameraModel":"camera_model","focalLength":"focal_length",` +
		`"apertureFnumber":"aperture_fnumber","isoEquivalent":"iso_equivalent","exposureTime":"exposure_time",` +
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
			"",
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPlaces
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"get mediaitem places with empty table",
			http.MethodGet,
			"/v1/mediaItems/:id/places",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/places",
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "place_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(placeCols))
			},
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
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "place_mediaitems"`)).
					WillReturnRows(getMockedPlaceRows())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
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
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "place_mediaitems"`)).
					WillReturnError(errors.New("some db error"))
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
			"get mediaitem things bad request",
			http.MethodGet,
			"/v1/mediaItems/:id/things",
			"/v1/mediaItems/bad-uuid/things",
			"",
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemThings
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"get mediaitem things with empty table",
			http.MethodGet,
			"/v1/mediaItems/:id/things",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/things",
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "thing_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(thingCols))
			},
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
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "thing_mediaitems"`)).
					WillReturnRows(getMockedThingRows())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
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
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "thing_mediaitems"`)).
					WillReturnError(errors.New("some db error"))
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
			"get mediaitem people bad request",
			http.MethodGet,
			"/v1/mediaItems/:id/people",
			"/v1/mediaItems/bad-uuid/people",
			"",
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemPeople
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"get mediaitem people with empty table",
			http.MethodGet,
			"/v1/mediaItems/:id/people",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/people",
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "people_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(peopleCols))
			},
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
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "people_mediaitems"`)).
					WillReturnRows(getMockedPeopleRows())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
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
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "people_mediaitems"`)).
					WillReturnError(errors.New("some db error"))
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

func TestGetMediaItemAlbums(t *testing.T) {
	tests := []Test{
		{
			"get mediaitem albums bad request",
			http.MethodGet,
			"/v1/mediaItems/:id/albums",
			"/v1/mediaItems/bad-uuid/albums",
			"",
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemAlbums
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"get mediaitem albums with empty table",
			http.MethodGet,
			"/v1/mediaItems/:id/albums",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179/albums",
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "album_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(albumCols))
			},
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
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "album_mediaitems"`)).
					WillReturnRows(getMockedAlbumRows())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
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
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "album_mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItemAlbums
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
			http.MethodGet,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/bad-uuid",
			"",
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"get mediaitem not found",
			http.MethodGet,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(mediaitemCols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetMediaItem
			},
			http.StatusNotFound,
			`{"message":"mediaitem not found"}`,
		},
		{
			"get mediaitem",
			http.MethodGet,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRows())
			},
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
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
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
	tests := []Test{
		{
			"update mediaitem bad request",
			http.MethodPut,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/bad-uuid",
			"",
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateMediaItem
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"update mediaitem with no payload",
			http.MethodPut,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"",
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateMediaItem
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem"}`,
		},
		{
			"update mediaitem with bad payload",
			http.MethodPut,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			`{"bad":"request"}`,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateMediaItem
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem"}`,
		},
		{
			"update mediaitem with success",
			http.MethodPut,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			`{"description":"description","favourite":true,"hidden":true}`,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "description", true, true,
						sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
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
			`{"description":"description","favourite":true,"hidden":true}`,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "description", true, true,
						sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateMediaItem
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
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
			"",
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			},
			http.StatusBadRequest,
			`{"message":"invalid mediaitem id"}`,
		},
		{
			"delete mediaitem with success",
			http.MethodDelete,
			"/v1/mediaItems/:id",
			"/v1/mediaItems/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", true,
						sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
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
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "mediaitems"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", true,
						sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteMediaItem
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
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
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(mediaitemCols))
			},
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
			"/v1/mediaItems",
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRows())
			},
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
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnError(errors.New("some db error"))
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

func getMockedMediaItemRow() *sqlmock.Rows {
	return sqlmock.NewRows(mediaitemCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "filename", "description", "mime_type", "source_url", "preview_url",
			"thumbnail_url", "true", "false", "false", "status", "mediaitem_type", 720,
			480, sampleTime, "camera_make", "camera_model", "focal_length", "aperture_fnumber",
			"iso_equivalent", "exposure_time", "17.580249", "-70.278493", "fps", sampleTime, sampleTime)
}

func getMockedMediaItemRows() *sqlmock.Rows {
	return sqlmock.NewRows(mediaitemCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "filename", "description", "mime_type", "source_url", "preview_url",
			"thumbnail_url", "true", "false", "false", "status", "mediaitem_type", 720,
			480, sampleTime, "camera_make", "camera_model", "focal_length", "aperture_fnumber",
			"iso_equivalent", "exposure_time", "17.580249", "-70.278493", "fps", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "filename", "description", "mime_type", "source_url", "preview_url",
			"thumbnail_url", "false", "true", "true", "status", "mediaitem_type", 720,
			480, sampleTime, "camera_make", "camera_model", "focal_length", "aperture_fnumber",
			"iso_equivalent", "exposure_time", "17.580249", "-70.278493", "fps", sampleTime, sampleTime)
}

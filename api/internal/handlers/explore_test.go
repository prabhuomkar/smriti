package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"
)

var (
	placeCols = []string{"id", "user_id", "name", "postcode", "town", "city", "state",
		"country", "cover_mediaitem_id", "is_hidden", "created_at", "updated_at"}
	thingCols                  = []string{"id", "user_id", "name", "cover_mediaitem_id", "is_hidden", "created_at", "updated_at"}
	peopleCols                 = []string{"id", "user_id", "name", "cover_mediaitem_id", "is_hidden", "created_at", "updated_at"}
	coverMediaItemResponseBody = `"coverMediaItem":{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","filename":"filename",` +
		`"description":"description","mimeType":"mime_type","sourceUrl":"source_url","previewUrl":"preview_url",` +
		`"thumbnailUrl":"thumbnail_url","favourite":true,"hidden":false,"deleted":false,"status":"status",` +
		`"mediaItemType":"mediaitem_type","width":720,"height":480,"creationTime":"2022-09-22T11:22:33+05:30",` +
		`"cameraMake":"camera_make","cameraModel":"camera_model","focalLength":"focal_length",` +
		`"apertureFnumber":"aperture_fnumber","isoEquivalent":"iso_equivalent","exposureTime":"exposure_time",` +
		`"latitude":17.580249,"longitude":-70.278493,"fps":"fps","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}`
	placeResponseBody = `{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"name":"name","postcode":"postcode",` +
		`"town":"town","city":"city",` +
		`"state":"state","country":"country","hidden":true,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` +
		coverMediaItemResponseBody + `}`
	placesResponseBody = `[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"name":"name","postcode":"postcode",` +
		`"town":"town","city":"city",` +
		`"state":"state","country":"country","hidden":true,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` +
		coverMediaItemResponseBody + `},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name",` +
		`"postcode":"postcode","town":"town","city":"city",` +
		`"state":"state","country":"country","hidden":false,` +
		`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30",` + coverMediaItemResponseBody + `}]`
	thingResponseBody = `{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"name":"name",` +
		`"hidden":true,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` +
		coverMediaItemResponseBody + `}`
	thingsResponseBody = `[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"name":"name",` +
		`"hidden":true,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` +
		coverMediaItemResponseBody + `},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name",` +
		`"hidden":false,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30",` + coverMediaItemResponseBody + `}]`
	personResponseBody = `{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name",` +
		`"hidden":true,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` +
		coverMediaItemResponseBody + `}`
	peopleResponseBody = `[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name",` +
		`"hidden":true,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30",` +
		coverMediaItemResponseBody + `},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180",` +
		`"userId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name",` +
		`"hidden":false,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30",` + coverMediaItemResponseBody + `}]`
)

func TestGetPlaces(t *testing.T) {
	tests := []Test{
		{
			"get places with empty table",
			http.MethodGet,
			"/v1/explore/places",
			"/v1/explore/places",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`)).
					WillReturnRows(sqlmock.NewRows(placeCols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaces
			},
			http.StatusOK,
			"[]",
		},
		{
			"get places with 2 rows",
			http.MethodGet,
			"/v1/explore/places",
			"/v1/explore/places",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`)).
					WillReturnRows(getMockedPlaceRows())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaces
			},
			http.StatusOK,
			placesResponseBody,
		},
		{
			"get places with error",
			http.MethodGet,
			"/v1/explore/places",
			"/v1/explore/places",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`)).
					WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaces
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetPlace(t *testing.T) {
	tests := []Test{
		{
			"get place bad request",
			http.MethodGet,
			"/v1/explore/places/:id",
			"/v1/explore/places/bad-uuid",
			map[string]string{},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlace
			},
			http.StatusBadRequest,
			`{"message":"invalid place id"}`,
		},
		{
			"get place not found",
			http.MethodGet,
			"/v1/explore/places/:id",
			"/v1/explore/places/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`)).
					WillReturnRows(sqlmock.NewRows(placeCols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlace
			},
			http.StatusNotFound,
			`{"message":"place not found"}`,
		},
		{
			"get place",
			http.MethodGet,
			"/v1/explore/places/:id",
			"/v1/explore/places/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`)).
					WillReturnRows(getMockedPlaceRow())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlace
			},
			http.StatusOK,
			placeResponseBody,
		},
		{
			"get place with error",
			http.MethodGet,
			"/v1/explore/places/:id",
			"/v1/explore/places/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`)).
					WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlace
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetPlaceMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get place mediaitems bad request",
			http.MethodGet,
			"/v1/places/:id/mediaItems",
			"/v1/places/bad-uuid/mediaItems",
			map[string]string{},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaceMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid place id"}`,
		},
		{
			"get place mediaitems not found",
			http.MethodGet,
			"/v1/places/:id/mediaItems",
			"/v1/places/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "place_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(mediaitemCols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaceMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get place mediaitems with 2 rows",
			http.MethodGet,
			"/v1/places/:id/mediaItems",
			"/v1/places/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "place_mediaitems"`)).
					WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaceMediaItems
			},
			http.StatusOK,
			mediaitemsResponseBody,
		},
		{
			"get place mediaitems with error",
			http.MethodGet,
			"/v1/places/:id/mediaItems",
			"/v1/places/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "place_mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaceMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetThings(t *testing.T) {
	tests := []Test{
		{
			"get things with empty table",
			http.MethodGet,
			"/v1/explore/things",
			"/v1/explore/things",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`)).
					WillReturnRows(sqlmock.NewRows(thingCols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThings
			},
			http.StatusOK,
			"[]",
		},
		{
			"get things with 2 rows",
			http.MethodGet,
			"/v1/explore/things",
			"/v1/explore/things",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`)).
					WillReturnRows(getMockedThingRows())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThings
			},
			http.StatusOK,
			thingsResponseBody,
		},
		{
			"get things with error",
			http.MethodGet,
			"/v1/explore/things",
			"/v1/explore/things",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`)).
					WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThings
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetThing(t *testing.T) {
	tests := []Test{
		{
			"get thing bad request",
			http.MethodGet,
			"/v1/explore/things/:id",
			"/v1/explore/things/bad-uuid",
			map[string]string{},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThing
			},
			http.StatusBadRequest,
			`{"message":"invalid thing id"}`,
		},
		{
			"get thing not found",
			http.MethodGet,
			"/v1/explore/things/:id",
			"/v1/explore/things/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`)).
					WillReturnRows(sqlmock.NewRows(thingCols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThing
			},
			http.StatusNotFound,
			`{"message":"thing not found"}`,
		},
		{
			"get thing",
			http.MethodGet,
			"/v1/explore/things/:id",
			"/v1/explore/things/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`)).
					WillReturnRows(getMockedThingRow())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThing
			},
			http.StatusOK,
			thingResponseBody,
		},
		{
			"get thing with error",
			http.MethodGet,
			"/v1/explore/things/:id",
			"/v1/explore/things/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`)).
					WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThing
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetThingMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get thing mediaitems bad request",
			http.MethodGet,
			"/v1/things/:id/mediaItems",
			"/v1/things/bad-uuid/mediaItems",
			map[string]string{},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThingMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid thing id"}`,
		},
		{
			"get thing mediaitems not found",
			http.MethodGet,
			"/v1/things/:id/mediaItems",
			"/v1/things/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "thing_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(mediaitemCols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThingMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get thing mediaitems with 2 rows",
			http.MethodGet,
			"/v1/things/:id/mediaItems",
			"/v1/things/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "thing_mediaitems"`)).
					WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThingMediaItems
			},
			http.StatusOK,
			mediaitemsResponseBody,
		},
		{
			"get thing mediaitems with error",
			http.MethodGet,
			"/v1/things/:id/mediaItems",
			"/v1/things/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "thing_mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThingMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestUpdatePeople(t *testing.T) {
	tests := []Test{
		{
			"update people bad request",
			http.MethodPut,
			"/v1/people/:id",
			"/v1/people/bad-uuid",
			map[string]string{},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdatePerson
			},
			http.StatusBadRequest,
			`{"message":"invalid people id"}`,
		},
		{
			"update people with no payload",
			http.MethodPut,
			"/v1/people/:id",
			"/v1/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdatePerson
			},
			http.StatusBadRequest,
			`{"message":"invalid people"}`,
		},
		{
			"update people with bad payload",
			http.MethodPut,
			"/v1/people/:id",
			"/v1/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdatePerson
			},
			http.StatusBadRequest,
			`{"message":"invalid people"}`,
		},
		{
			"update people with bad cover mediaitem id",
			http.MethodPut,
			"/v1/people/:id",
			"/v1/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"name":"name","coverMediaItemId":"bad-mediaitem-id"}`),
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdatePerson
			},
			http.StatusBadRequest,
			`{"message":"invalid people cover mediaitem id"}`,
		},
		{
			"update people with success",
			http.MethodPut,
			"/v1/people/:id",
			"/v1/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"name":"name","hidden":true,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "people"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", true,
						"4d05b5f6-17c2-475e-87fe-3fc8b9567179", sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdatePerson
			},
			http.StatusNoContent,
			"",
		},
		{
			"update people with error",
			http.MethodPut,
			"/v1/people/:id",
			"/v1/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"name":"name","hidden":true,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "people"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", true,
						"4d05b5f6-17c2-475e-87fe-3fc8b9567179", sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdatePerson
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetPeople(t *testing.T) {
	tests := []Test{
		{
			"get people with empty table",
			http.MethodGet,
			"/v1/explore/people",
			"/v1/explore/people",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`)).
					WillReturnRows(sqlmock.NewRows(peopleCols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeople
			},
			http.StatusOK,
			"[]",
		},
		{
			"get people with 2 rows",
			http.MethodGet,
			"/v1/explore/people",
			"/v1/explore/people",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`)).
					WillReturnRows(getMockedPeopleRows())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeople
			},
			http.StatusOK,
			peopleResponseBody,
		},
		{
			"get people with error",
			http.MethodGet,
			"/v1/explore/people",
			"/v1/explore/people",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`)).
					WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeople
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetPerson(t *testing.T) {
	tests := []Test{
		{
			"get person bad request",
			http.MethodGet,
			"/v1/explore/people/:id",
			"/v1/explore/people/bad-uuid",
			map[string]string{},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPerson
			},
			http.StatusBadRequest,
			`{"message":"invalid person id"}`,
		},
		{
			"get person not found",
			http.MethodGet,
			"/v1/explore/people/:id",
			"/v1/explore/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`)).
					WillReturnRows(sqlmock.NewRows(peopleCols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPerson
			},
			http.StatusNotFound,
			`{"message":"person not found"}`,
		},
		{
			"get people",
			http.MethodGet,
			"/v1/explore/people/:id",
			"/v1/explore/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`)).
					WillReturnRows(getMockedPeopleRow())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mediaitems"`)).
					WillReturnRows(getMockedMediaItemRow())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPerson
			},
			http.StatusOK,
			personResponseBody,
		},
		{
			"get person with error",
			http.MethodGet,
			"/v1/explore/people/:id",
			"/v1/explore/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`)).
					WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPerson
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func TestGetPeopleMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get people mediaitems bad request",
			http.MethodGet,
			"/v1/people/:id/mediaItems",
			"/v1/people/bad-uuid/mediaItems",
			map[string]string{},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeopleMediaItems
			},
			http.StatusBadRequest,
			`{"message":"invalid people id"}`,
		},
		{
			"get people mediaitems not found",
			http.MethodGet,
			"/v1/people/:id/mediaItems",
			"/v1/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "people_mediaitems"`)).
					WillReturnRows(sqlmock.NewRows(mediaitemCols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeopleMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get people mediaitems with 2 rows",
			http.MethodGet,
			"/v1/people/:id/mediaItems",
			"/v1/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "people_mediaitems"`)).
					WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeopleMediaItems
			},
			http.StatusOK,
			mediaitemsResponseBody,
		},
		{
			"get people mediaitems with error",
			http.MethodGet,
			"/v1/people/:id/mediaItems",
			"/v1/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`JOIN "people_mediaitems"`)).
					WillReturnError(errors.New("some db error"))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeopleMediaItems
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
	}
	executeTests(t, tests)
}

func getMockedPlaceRow() *sqlmock.Rows {
	return sqlmock.NewRows(placeCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"name", "postcode", "town", "city",
			"state", "country", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "true", sampleTime, sampleTime)
}

func getMockedPlaceRows() *sqlmock.Rows {
	return sqlmock.NewRows(placeCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"name", "postcode", "town", "city",
			"state", "country", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "true", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"name", "postcode", "town", "city",
			"state", "country", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "false", sampleTime, sampleTime)
}

func getMockedThingRow() *sqlmock.Rows {
	return sqlmock.NewRows(thingCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"name", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "true", sampleTime, sampleTime)
}

func getMockedThingRows() *sqlmock.Rows {
	return sqlmock.NewRows(thingCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"name", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "true", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"name", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "false", sampleTime, sampleTime)
}

func getMockedPeopleRow() *sqlmock.Rows {
	return sqlmock.NewRows(peopleCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"name", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "true", sampleTime, sampleTime)
}

func getMockedPeopleRows() *sqlmock.Rows {
	return sqlmock.NewRows(peopleCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"name", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "true", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"name", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "false", sampleTime, sampleTime)
}

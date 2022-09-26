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
	place_cols = []string{"id", "name", "postcode", "suburb", "road", "town", "city", "county", "district", "state",
		"country", "cover_mediaitem_id", "cover_mediaitem_thumbnail_url", "is_hidden", "created_at", "updated_at"}
	thing_cols = []string{"id", "name", "cover_mediaitem_id", "cover_mediaitem_thumbnail_url",
		"is_hidden", "created_at", "updated_at"}
	people_cols = []string{"id", "name", "cover_mediaitem_id", "cover_mediaitem_thumbnail_url",
		"is_hidden", "created_at", "updated_at"}
	place_response_body = `{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name","postcode":"postcode",` +
		`"suburb":"suburb","road":"road","town":"town","city":"city","county":"county","district":"district",` +
		`"state":"state","country":"country","hidden":true,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}`
	places_response_body = `[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name","postcode":"postcode",` +
		`"suburb":"suburb","road":"road","town":"town","city":"city","county":"county","district":"district",` +
		`"state":"state","country":"country","hidden":true,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180","name":"name",` +
		`"postcode":"postcode","suburb":"suburb","road":"road","town":"town","city":"city","county":"county",` +
		`"district":"district","state":"state","country":"country","hidden":false,` +
		`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}]`
	thing_response_body = `{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name","hidden":true,` +
		`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}`
	things_response_body = `[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name","hidden":true,` +
		`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180","name":"name",` +
		`"hidden":false,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}]`
	person_response_body = `{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name","hidden":true,` +
		`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}`
	people_response_body = `[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name","hidden":true,` +
		`"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"},{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180","name":"name",` +
		`"hidden":false,"coverMediaItemId":"4d05b5f6-17c2-475e-87fe-3fc8b9567179",` +
		`"coverMediaItemThumbnailUrl":"cover_mediaitem_thumbnail_url","createdAt":"2022-09-22T11:22:33+05:30",` +
		`"updatedAt":"2022-09-22T11:22:33+05:30"}]`
)

func TestGetPlaces(t *testing.T) {
	tests := []Test{
		{
			"get places with empty table",
			http.MethodGet,
			"/v1/explore/places",
			"/v1/explore/places",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(place_cols))
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`))
				expectedMock.WillReturnRows(getMockedPlaceRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaces
			},
			http.StatusOK,
			places_response_body,
		},
		{
			"get places with error",
			http.MethodGet,
			"/v1/explore/places",
			"/v1/explore/places",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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
			"",
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(place_cols))
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`))
				expectedMock.WillReturnRows(getMockedPlaceRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlace
			},
			http.StatusOK,
			place_response_body,
		},
		{
			"get place with error",
			http.MethodGet,
			"/v1/explore/places/:id",
			"/v1/explore/places/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "places"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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
			"",
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`JOIN "place_mediaitems"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`JOIN "place_mediaitems"`))
				expectedMock.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaceMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get place mediaitems with error",
			http.MethodGet,
			"/v1/places/:id/mediaItems",
			"/v1/places/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`JOIN "place_mediaitems"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(thing_cols))
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`))
				expectedMock.WillReturnRows(getMockedThingRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThings
			},
			http.StatusOK,
			things_response_body,
		},
		{
			"get things with error",
			http.MethodGet,
			"/v1/explore/things",
			"/v1/explore/things",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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
			"",
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(thing_cols))
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`))
				expectedMock.WillReturnRows(getMockedThingRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThing
			},
			http.StatusOK,
			thing_response_body,
		},
		{
			"get thing with error",
			http.MethodGet,
			"/v1/explore/things/:id",
			"/v1/explore/things/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "things"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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
			"",
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`JOIN "thing_mediaitems"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`JOIN "thing_mediaitems"`))
				expectedMock.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThingMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get thing mediaitems with error",
			http.MethodGet,
			"/v1/things/:id/mediaItems",
			"/v1/things/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`JOIN "thing_mediaitems"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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
func TestGetPeople(t *testing.T) {
	tests := []Test{
		{
			"get people with empty table",
			http.MethodGet,
			"/v1/explore/people",
			"/v1/explore/people",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(people_cols))
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`))
				expectedMock.WillReturnRows(getMockedPeopleRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeople
			},
			http.StatusOK,
			people_response_body,
		},
		{
			"get people with error",
			http.MethodGet,
			"/v1/explore/people",
			"/v1/explore/people",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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
			"",
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(people_cols))
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`))
				expectedMock.WillReturnRows(getMockedPeopleRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPerson
			},
			http.StatusOK,
			person_response_body,
		},
		{
			"get person with error",
			http.MethodGet,
			"/v1/explore/people/:id",
			"/v1/explore/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "people"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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
			"",
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`JOIN "people_mediaitems"`))
				expectedMock.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
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
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`JOIN "people_mediaitems"`))
				expectedMock.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeopleMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get people mediaitems with error",
			http.MethodGet,
			"/v1/people/:id/mediaItems",
			"/v1/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			"",
			func(mock sqlmock.Sqlmock) {
				expectedMock := mock.ExpectQuery(regexp.QuoteMeta(`JOIN "people_mediaitems"`))
				expectedMock.WillReturnError(errors.New("some db error"))
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

func getMockedPlaceRows() *sqlmock.Rows {
	return sqlmock.NewRows(place_cols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "postcode", "suburb", "road", "town", "city", "county",
			"district", "state", "country", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "cover_mediaitem_thumbnail_url",
			"true", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "name", "postcode", "suburb", "road", "town", "city", "county",
			"district", "state", "country", "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "cover_mediaitem_thumbnail_url",
			"false", sampleTime, sampleTime)
}

func getMockedThingRows() *sqlmock.Rows {
	return sqlmock.NewRows(thing_cols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"cover_mediaitem_thumbnail_url", "true", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "name", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"cover_mediaitem_thumbnail_url", "false", sampleTime, sampleTime)
}

func getMockedPeopleRows() *sqlmock.Rows {
	return sqlmock.NewRows(people_cols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"cover_mediaitem_thumbnail_url", "true", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "name", "4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			"cover_mediaitem_thumbnail_url", "false", sampleTime, sampleTime)
}

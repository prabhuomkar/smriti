package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"
)

func TestGetPlaces(t *testing.T) {

}
func TestGetPlace(t *testing.T) {

}
func TestGetPlaceMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get place mediaitems bad request",
			"/v1/places/:id/mediaItems",
			"/v1/places/bad-uuid/mediaItems",
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
			"/v1/places/:id/mediaItems",
			"/v1/places/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM place_mediaitems`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaceMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get place mediaitems with 2 rows",
			"/v1/places/:id/mediaItems",
			"/v1/places/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM place_mediaitems`))
				expectedQuery.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPlaceMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get place mediaitems with error",
			"/v1/places/:id/mediaItems",
			"/v1/places/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM place_mediaitems`))
				expectedQuery.WillReturnError(errors.New("some db error"))
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

}
func TestGetThing(t *testing.T) {

}
func TestGetThingMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get thing mediaitems bad request",
			"/v1/things/:id/mediaItems",
			"/v1/things/bad-uuid/mediaItems",
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
			"/v1/things/:id/mediaItems",
			"/v1/things/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM thing_mediaitems`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThingMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get thing mediaitems with 2 rows",
			"/v1/things/:id/mediaItems",
			"/v1/things/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM thing_mediaitems`))
				expectedQuery.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetThingMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get thing mediaitems with error",
			"/v1/things/:id/mediaItems",
			"/v1/things/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM thing_mediaitems`))
				expectedQuery.WillReturnError(errors.New("some db error"))
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

}
func TestGetPerson(t *testing.T) {

}
func TestGetPeopleMediaItems(t *testing.T) {
	tests := []Test{
		{
			"get people mediaitems bad request",
			"/v1/people/:id/mediaItems",
			"/v1/people/bad-uuid/mediaItems",
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
			"/v1/people/:id/mediaItems",
			"/v1/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM people_mediaitems`))
				expectedQuery.WillReturnRows(sqlmock.NewRows(mediaitem_cols))
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeopleMediaItems
			},
			http.StatusOK,
			"[]",
		},
		{
			"get people mediaitems with 2 rows",
			"/v1/people/:id/mediaItems",
			"/v1/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM people_mediaitems`))
				expectedQuery.WillReturnRows(getMockedMediaItemRows())
			},
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetPeopleMediaItems
			},
			http.StatusOK,
			mediaitems_response_body,
		},
		{
			"get people mediaitems with error",
			"/v1/people/:id/mediaItems",
			"/v1/people/4d05b5f6-17c2-475e-87fe-3fc8b9567179/mediaItems",
			nil,
			func(mock sqlmock.Sqlmock) {
				expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM people_mediaitems`))
				expectedQuery.WillReturnError(errors.New("some db error"))
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

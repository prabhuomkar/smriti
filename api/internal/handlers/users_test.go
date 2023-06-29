package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
)

var (
	userCols         = []string{"id", "name", "username", "password", "created_at", "updated_at"}
	userResponseBody = `{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name",` +
		`"username":"username",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30"}`
	usersResponseBody = `[{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567179","name":"name",` +
		`"username":"username",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30"},` +
		`{"id":"4d05b5f6-17c2-475e-87fe-3fc8b9567180","name":"name",` +
		`"username":"username",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30"}]`
)

func TestGetUser(t *testing.T) {
	tests := []Test{
		{
			"get user bad request",
			http.MethodGet,
			"/v1/users/:id",
			"/v1/users/bad-uuid",
			[]string{"id"},
			[]string{"bad-uuid"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUser
			},
			http.StatusBadRequest,
			"invalid user id",
		},
		{
			"get user not found",
			http.MethodGet,
			"/v1/users/:id",
			"/v1/users/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
					WillReturnRows(sqlmock.NewRows(userCols))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUser
			},
			http.StatusNotFound,
			"user not found",
		},
		{
			"get user",
			http.MethodGet,
			"/v1/users/:id",
			"/v1/users/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
					WillReturnRows(getMockedUserRow())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUser
			},
			http.StatusOK,
			userResponseBody,
		},
		{
			"get user with error",
			http.MethodGet,
			"/v1/users/:id",
			"/v1/users/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUser
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

func TestUpdateUser(t *testing.T) {
	tests := []Test{
		{
			"update user bad request",
			http.MethodPut,
			"/v1/users/:id",
			"/v1/users/bad-uuid",
			[]string{"id"},
			[]string{"bad-uuid"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateUser
			},
			http.StatusBadRequest,
			"invalid user id",
		},
		{
			"update user with no payload",
			http.MethodPut,
			"/v1/users/:id",
			"/v1/users/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateUser
			},
			http.StatusBadRequest,
			"invalid user",
		},
		{
			"update user with bad payload",
			http.MethodPut,
			"/v1/users/:id",
			"/v1/users/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
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
				return handler.UpdateUser
			},
			http.StatusBadRequest,
			"invalid user",
		},
		{
			"update user with success",
			http.MethodPut,
			"/v1/users/:id",
			"/v1/users/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"name":"name","username":"username","password":"password"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "username", sqlmock.AnyArg(),
						sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateUser
			},
			http.StatusNoContent,
			"",
		},
		{
			"update user with error",
			http.MethodPut,
			"/v1/users/:id",
			"/v1/users/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"name":"name","username":"username","password":"password"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "username", sqlmock.AnyArg(),
						sqlmock.AnyArg(), "4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateUser
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

func TestDeleteUser(t *testing.T) {
	tests := []Test{
		{
			"delete user bad request",
			http.MethodDelete,
			"/v1/users/:id",
			"/v1/users/bad-uuid",
			[]string{"id"},
			[]string{"bad-uuid"},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteUser
			},
			http.StatusBadRequest,
			"invalid user id",
		},
		{
			"delete user with success",
			http.MethodDelete,
			"/v1/users/:id",
			"/v1/users/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteUser
			},
			http.StatusNoContent,
			"",
		},
		{
			"delete user with error",
			http.MethodDelete,
			"/v1/users/:id",
			"/v1/users/4d05b5f6-17c2-475e-87fe-3fc8b9567179",
			[]string{"id"},
			[]string{"4d05b5f6-17c2-475e-87fe-3fc8b9567179"},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users"`)).
					WithArgs("4d05b5f6-17c2-475e-87fe-3fc8b9567179").
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteUser
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

func TestGetUsers(t *testing.T) {
	tests := []Test{
		{
			"get users with empty table",
			http.MethodGet,
			"/v1/users",
			"/v1/users",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
					WillReturnRows(sqlmock.NewRows(userCols))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUsers
			},
			http.StatusOK,
			"[]",
		},
		{
			"get users with 2 rows",
			http.MethodGet,
			"/v1/users",
			"/v1/users",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
					WillReturnRows(getMockedUserRows())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUsers
			},
			http.StatusOK,
			usersResponseBody,
		},
		{
			"get users with error",
			http.MethodGet,
			"/v1/users",
			"/v1/users",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUsers
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

func TestCreateUser(t *testing.T) {
	tests := []Test{
		{
			"create user with bad payload",
			http.MethodPost,
			"/v1/users",
			"/v1/users",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateUser
			},
			http.StatusBadRequest,
			"invalid user",
		},
		{
			"create user with no payload",
			http.MethodPost,
			"/v1/users",
			"/v1/users",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateUser
			},
			http.StatusBadRequest,
			"invalid user",
		},
		{
			"create user with success",
			http.MethodPost,
			"/v1/users",
			"/v1/users",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"name":"name","username":"username","password":"password","features":"{\"albums\":true}"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
					WithArgs(sqlmock.AnyArg(), "name", "username", sqlmock.AnyArg(), "{\"albums\":true}", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateUser
			},
			http.StatusCreated,
			`"name":"name","username":"username"`,
		},
		{
			"create user with error",
			http.MethodPost,
			"/v1/users",
			"/v1/users",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"name":"name","username":"username","password":"password","features":"{\"albums\":true}"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
					WithArgs(sqlmock.AnyArg(), "name", "username", sqlmock.AnyArg(), "{\"albums\":true}", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
				mock.ExpectRollback()
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateUser
			},
			http.StatusInternalServerError,
			"some db error",
		},
	}
	executeTests(t, tests)
}

func getMockedUserRow() *sqlmock.Rows {
	return sqlmock.NewRows(userCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "username", "password", sampleTime, sampleTime)
}

func getMockedUserRows() *sqlmock.Rows {
	return sqlmock.NewRows(userCols).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567179", "name", "username", "password", sampleTime, sampleTime).
		AddRow("4d05b5f6-17c2-475e-87fe-3fc8b9567180", "name", "username", "password", sampleTime, sampleTime)
}

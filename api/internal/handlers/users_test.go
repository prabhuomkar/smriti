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
)

var (
	userCols = []string{
		"id", "name", "username", "password", "features", "created_at", "updated_at",
	}
	userResponseBody = `{"id":"019b7796-6072-76ee-8be3-485ff2b32fd7","name":"name",` +
		`"username":"username",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30"}`
	usersResponseBody = `[{"id":"019b7796-6072-76ee-8be3-485ff2b32fd7","name":"name",` +
		`"username":"username",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30"},` +
		`{"id":"019b7796-6072-76ee-8be3-485ff2b33fd7","name":"name",` +
		`"username":"username",` +
		`"createdAt":"2022-09-22T11:22:33+05:30","updatedAt":"2022-09-22T11:22:33+05:30"}]`
)

func TestGetUser(t *testing.T) {
	tests := []Test{
		{
			"get user bad request", http.MethodGet, "/v1/users/:id", "/v1/users/bad-uuid", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUser
			}, http.StatusBadRequest, "invalid user id",
		},
		{
			"get user not found", http.MethodGet, "/v1/users/:id", "/v1/users/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users`)).
					WithArgs(pgxmock.AnyArg()).WillReturnError(pgx.ErrNoRows)
			}, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUser
			}, http.StatusNotFound, "user not found",
		},
		{
			"get user with error", http.MethodGet, "/v1/users/:id", "/v1/users/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users`)).
					WithArgs(pgxmock.AnyArg()).WillReturnError(errors.New("some db error"))
			}, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUser
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get user with error in scanning", http.MethodGet, "/v1/users/:id", "/v1/users/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users`)).
					WithArgs(pgxmock.AnyArg()).WillReturnRows(pgxmock.NewRows(userCols).AddRow("invalid", "name", "username", "password", "", "invalid", "invalid"))
			}, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUser
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get user with success", http.MethodGet, "/v1/users/:id", "/v1/users/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users`)).
					WithArgs(pgxmock.AnyArg()).WillReturnRows(getMockedUserRow())
			}, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUser
			}, http.StatusOK, userResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestUpdateUser(t *testing.T) {
	tests := []Test{
		{
			"update user bad request", http.MethodPut, "/v1/users/:id", "/v1/users/bad-uuid", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateUser
			}, http.StatusBadRequest, "invalid user id",
		},
		{
			"update user with no payload", http.MethodPut, "/v1/users/:id", "/v1/users/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateUser
			}, http.StatusBadRequest, "invalid user",
		},
		{
			"update user with bad payload", http.MethodPut, "/v1/users/:id", "/v1/users/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"bad":"request}`), nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateUser
			}, http.StatusBadRequest, "invalid user",
		},
		{
			"update user with error", http.MethodPut, "/v1/users/:id", "/v1/users/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"name":"name","username":"username","password":"password"}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE users`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateUser
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"update user with success", http.MethodPut, "/v1/users/:id", "/v1/users/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"name":"name","username":"username","password":"password"}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE users`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			}, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.UpdateUser
			}, http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestDeleteUser(t *testing.T) {
	tests := []Test{
		{
			"delete user bad request", http.MethodDelete, "/v1/users/:id", "/v1/users/bad-uuid", []string{"id"}, []string{"bad-uuid"}, map[string]string{}, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteUser
			}, http.StatusBadRequest, "invalid user id",
		},
		{
			"delete user with error", http.MethodDelete, "/v1/users/:id", "/v1/users/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteUser
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"delete user with success", http.MethodDelete, "/v1/users/:id", "/v1/users/019b7796-6072-76ee-8be3-485ff2b32fd7", []string{"id"}, []string{"019b7796-6072-76ee-8be3-485ff2b32fd7"}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users`)).
					WithArgs(pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			}, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.DeleteUser
			}, http.StatusNoContent, "",
		},
	}
	executeTests(t, tests)
}

func TestGetUsers(t *testing.T) {
	tests := []Test{
		{
			"get users with error", http.MethodGet, "/v1/users", "/v1/users", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUsers
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"get users with error in scanning", http.MethodGet, "/v1/users", "/v1/users", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(userCols).
						AddRow("invalid", "name", "username", "password", "", "invalid", "invalid"))
			}, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUsers
			}, http.StatusInternalServerError, "Scanning value error",
		},
		{
			"get users with empty table", http.MethodGet, "/v1/users", "/v1/users", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows(userCols))
			}, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUsers
			}, http.StatusOK, "[]",
		},
		{
			"get users with 2 rows", http.MethodGet, "/v1/users", "/v1/users", []string{}, []string{}, map[string]string{}, nil, func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnRows(getMockedUserRows())
			}, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetUsers
			}, http.StatusOK, usersResponseBody,
		},
	}
	executeTests(t, tests)
}

func TestCreateUser(t *testing.T) {
	tests := []Test{
		{
			"create user with bad payload", http.MethodPost, "/v1/users", "/v1/users", []string{}, []string{}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"bad":"request"}`), nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateUser
			}, http.StatusBadRequest, "invalid user",
		},
		{
			"create user with no payload", http.MethodPost, "/v1/users", "/v1/users", []string{}, []string{}, map[string]string{}, nil, nil, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateUser
			}, http.StatusBadRequest, "invalid user",
		},
		{
			"create user with error", http.MethodPost, "/v1/users", "/v1/users", []string{}, []string{}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"name":"name","username":"username","password":"password","features":"{\"albums\":true}"}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users`)).
					WithArgs(pgxmock.AnyArg(), "name", "username", pgxmock.AnyArg(), "{\"albums\":true}", pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("some db error"))
			}, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateUser
			}, http.StatusInternalServerError, "some db error",
		},
		{
			"create user with success", http.MethodPost, "/v1/users", "/v1/users", []string{}, []string{}, map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			}, strings.NewReader(`{"name":"name","username":"username","password":"password","features":"{\"albums\":true}"}`), func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users`)).
					WithArgs(pgxmock.AnyArg(), "name", "username", pgxmock.AnyArg(), "{\"albums\":true}", pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			}, nil, func(handler *Handler) func(ctx echo.Context) error {
				return handler.CreateUser
			}, http.StatusCreated, `"name":"name","username":"username"`,
		},
	}
	executeTests(t, tests)
}

func getMockedUserRow() *pgxmock.Rows {
	return pgxmock.NewRows(userCols).
		AddRow("019b7796-6072-76ee-8be3-485ff2b32fd7", "name", "username", "password", "", sampleTime, sampleTime)
}

func getMockedUserRows() *pgxmock.Rows {
	return pgxmock.NewRows(userCols).
		AddRow("019b7796-6072-76ee-8be3-485ff2b32fd7", "name", "username", "password", "", sampleTime, sampleTime).
		AddRow("019b7796-6072-76ee-8be3-485ff2b33fd7", "name", "username", "password", "", sampleTime, sampleTime)
}

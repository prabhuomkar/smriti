package handlers

import (
	"api/config"
	"api/internal/auth"
	"api/internal/models"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

func TestLogin(t *testing.T) {
	tests := []Test{
		{
			"login with bad payload",
			http.MethodPost,
			"/v1/auth/login",
			"/v1/auth/login",
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
				return handler.Login
			},
			http.StatusBadRequest,
			"invalid username or password",
		},
		{
			"login with no payload",
			http.MethodPost,
			"/v1/auth/login",
			"/v1/auth/login",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Login
			},
			http.StatusBadRequest,
			"invalid username or password",
		},
		{
			"login with incomplete credentials",
			http.MethodPost,
			"/v1/auth/login",
			"/v1/auth/login",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"username":"username"}`),
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Login
			},
			http.StatusBadRequest,
			"invalid username or password",
		},
		{
			"login with success",
			http.MethodPost,
			"/v1/auth/login",
			"/v1/auth/login",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"username":"username","password":"password"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
					WillReturnRows(getMockedUserRow())
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Login
			},
			http.StatusOK,
			`"accessToken"`,
		},
		{
			"login with no user found",
			http.MethodPost,
			"/v1/auth/login",
			"/v1/auth/login",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"username":"username","password":"password"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
					WillReturnRows(sqlmock.NewRows(userCols))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Login
			},
			http.StatusNotFound,
			"incorrect username or password",
		},
		{
			"login with error",
			http.MethodPost,
			"/v1/auth/login",
			"/v1/auth/login",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"username":"username","password":"password"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Login
			},
			http.StatusInternalServerError,
			"some db error",
		},
		{
			"login with error getting tokens",
			http.MethodPost,
			"/v1/auth/login",
			"/v1/auth/login",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"username":"username","password":"password"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
					WillReturnRows(getMockedUserRow())
			},
			[]func(interface{}, interface{}) (interface{}, error){
				func(a interface{}, b interface{}) (interface{}, error) {
					val, ok := b.(bool)
					if ok && val == true {
						return b, nil
					}
					return nil, errors.New("some cache error")
				},
				nil,
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Login
			},
			http.StatusInternalServerError,
			"error getting tokens",
		},
	}
	executeTests(t, tests)
}

func TestRefresh(t *testing.T) {
	_, rtoken := auth.GetAccessAndRefreshTokens(&config.Config{Auth: config.Auth{RefreshTTL: 60}},
		models.User{ID: uuid.FromStringOrNil("4d05b5f6-17c2-475e-87fe-3fc8b9567179"), Username: "username"})
	tests := []Test{
		{
			"refresh with success",
			http.MethodPost,
			"/v1/auth/refresh",
			"/v1/auth/refresh",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderAuthorization: rtoken,
			},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Refresh
			},
			http.StatusOK,
			`"accessToken"`,
		},
		{
			"refresh with error",
			http.MethodPost,
			"/v1/auth/refresh",
			"/v1/auth/refresh",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Refresh
			},
			http.StatusInternalServerError,
			"error refreshing tokens",
		},
	}
	executeTests(t, tests)
}

func TestLogout(t *testing.T) {
	_, atoken := auth.GetAccessAndRefreshTokens(&config.Config{Auth: config.Auth{RefreshTTL: 60}},
		models.User{ID: uuid.FromStringOrNil("4d05b5f6-17c2-475e-87fe-3fc8b9567179"), Username: "username"})
	tests := []Test{
		{
			"logout with success",
			http.MethodPost,
			"/v1/auth/logout",
			"/v1/auth/logout",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderAuthorization: atoken,
			},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Logout
			},
			http.StatusNoContent,
			``,
		},
	}
	executeTests(t, tests)
}

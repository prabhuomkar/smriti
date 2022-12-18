package handlers

import (
	"api/config"
	"api/internal/auth"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"
)

func TestLogin(t *testing.T) {
	tests := []Test{
		{
			"login with bad payload",
			http.MethodPost,
			"/v1/auth/login",
			"/v1/auth/login",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"bad":"request"}`),
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Login
			},
			http.StatusBadRequest,
			`{"message":"invalid username or password"}`,
		},
		{
			"login with no payload",
			http.MethodPost,
			"/v1/auth/login",
			"/v1/auth/login",
			map[string]string{},
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Login
			},
			http.StatusBadRequest,
			`{"message":"invalid username or password"}`,
		},
		{
			"login with incomplete credentials",
			http.MethodPost,
			"/v1/auth/login",
			"/v1/auth/login",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"username":"username"}`),
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Login
			},
			http.StatusBadRequest,
			`{"message":"invalid username or password"}`,
		},
		{
			"login with success",
			http.MethodPost,
			"/v1/auth/login",
			"/v1/auth/login",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"username":"username","password":"password"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
					WillReturnRows(getMockedUserRow())
			},
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
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"username":"username","password":"password"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
					WillReturnRows(sqlmock.NewRows(userCols))
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Login
			},
			http.StatusNotFound,
			`{"message":"incorrect username or password"}`,
		},
		{
			"login with error",
			http.MethodPost,
			"/v1/auth/login",
			"/v1/auth/login",
			map[string]string{
				echo.HeaderContentType: echo.MIMEApplicationJSON,
			},
			strings.NewReader(`{"username":"username","password":"password"}`),
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
					WillReturnError(errors.New("some db error"))
			},
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Login
			},
			http.StatusInternalServerError,
			`{"message":"some db error"}`,
		},
		{
			"login with error getting tokens",
			http.MethodPost,
			"/v1/auth/login",
			"/v1/auth/login",
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
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Login
			},
			http.StatusInternalServerError,
			`{"message":"error getting tokens"}`,
		},
	}
	executeTests(t, tests)
}

func TestRefresh(t *testing.T) {
	_, rtoken := auth.GetAccessAndRefreshTokens(&config.Config{Auth: config.Auth{RefreshTTL: 60}}, "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "username")
	tests := []Test{
		{
			"refresh with success",
			http.MethodPost,
			"/v1/auth/refresh",
			"/v1/auth/refresh",
			map[string]string{
				echo.HeaderAuthorization: rtoken,
			},
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
			map[string]string{},
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.Refresh
			},
			http.StatusInternalServerError,
			`{"message":"error refreshing tokens"}`,
		},
	}
	executeTests(t, tests)
}

func TestLogout(t *testing.T) {
	_, atoken := auth.GetAccessAndRefreshTokens(&config.Config{Auth: config.Auth{RefreshTTL: 60}}, "4d05b5f6-17c2-475e-87fe-3fc8b9567179", "username")
	tests := []Test{
		{
			"logout with success",
			http.MethodPost,
			"/v1/auth/logout",
			"/v1/auth/logout",
			map[string]string{
				echo.HeaderAuthorization: atoken,
			},
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

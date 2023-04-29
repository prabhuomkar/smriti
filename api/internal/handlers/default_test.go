package handlers

import (
	"net/http"
	"testing"

	"github.com/labstack/echo"
)

func TestGetFeatures(t *testing.T) {
	tests := []Test{
		{
			"get features with error",
			http.MethodGet,
			"/v1/features",
			"/v1/features",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetFeatures
			},
			http.StatusOK,
			"{}",
		},
		{
			"get features successfully",
			http.MethodGet,
			"/v1/features",
			"/v1/features",
			[]string{},
			[]string{},
			map[string]string{
				echo.HeaderAuthorization: "atoken",
			},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetFeatures
			},
			http.StatusOK,
			`{"albums":true,"explore":true,"places":true}`,
		},
	}
	executeTests(t, tests)
}

func TestGetVersion(t *testing.T) {
	tests := []Test{
		{
			"get version successfully",
			http.MethodGet,
			"/v1/version",
			"/v1/version",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetVersion
			},
			http.StatusOK,
			``,
		},
	}
	executeTests(t, tests)
}

func TestGetDisk(t *testing.T) {
	tests := []Test{
		{
			"get disk successfully",
			http.MethodGet,
			"/v1/disk",
			"/v1/disk",
			[]string{},
			[]string{},
			map[string]string{},
			nil,
			nil,
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetDisk
			},
			http.StatusOK,
			``,
		},
	}
	executeTests(t, tests)
}

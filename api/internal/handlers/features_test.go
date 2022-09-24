package handlers

import (
	"net/http"
	"testing"

	"github.com/labstack/echo"
)

func TestGetFeatures(t *testing.T) {
	tests := []Test{
		{
			"get features",
			"/v1/features",
			"/v1/features",
			nil,
			nil,
			func(handler *Handler) func(ctx echo.Context) error {
				return handler.GetFeatures
			},
			http.StatusOK,
			`{"favourites":false,"hidden":false,"trash":false,"albums":false,` +
				`"explore":false,"places":false,"things":false,"people":false,"sharing":false}`,
		},
	}
	executeTests(t, tests)
}
package middlewares

import (
	"api/config"
	"api/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestFeature(t *testing.T) {
	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/albums", nil)
	rec := httptest.NewRecorder()

	// handler
	cfg := &config.Config{Feature: config.Feature{
		Albums:        false,
		Favourites:    false,
		Hidden:        false,
		Trash:         false,
		Explore:       false,
		ExplorePlaces: false,
		ExploreThings: false,
		ExplorePeople: false,
		Sharing:       false,
	}}
	handler := &handlers.Handler{
		Config: cfg,
		DB:     nil,
	}
	featureHandlerMap := map[string]interface{}{
		"albums":     handler.GetAlbums,
		"favourites": handler.GetFavouriteMediaItems,
		"hidden":     handler.GetHiddenMediaItems,
		"trash":      handler.GetDeletedMediaItems,
		"explore":    handler.GetPlaces,
		"places":     handler.GetPlaces,
		"things":     handler.GetThings,
		"people":     handler.GetPeople,
	}
	for feature, handler := range featureHandlerMap {
		checkFeature := FeatureCheck(cfg, feature)
		server.GET("/v1/albums", checkFeature(handler.(func(ctx echo.Context) error)))
		server.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	}
}

package server

import (
	"api/config"
	"api/internal/handlers"
	"api/internal/middlewares"
	"api/pkg/cache"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
)

const httpTimeout = 10

// StartHTTPServer ...
//
//nolint:funlen
func StartHTTPServer(handler *handlers.Handler) *http.Server {
	srvHandler := echo.New()
	//nolint:gosec
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", handler.Config.API.Host, handler.Config.API.Port),
		Handler: srvHandler,
	}
	// metrics
	srvHandler.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Namespace:  "http",
		Subsystem:  "server",
		Registerer: prometheus.DefaultRegisterer,
	}))
	srvHandler.GET("/metrics", echoprometheus.NewHandler())
	// file server
	if handler.Config.Storage.Provider == "disk" {
		fileRoute := getFileRoute(handler.Config.Storage.DiskRoot)
		slog.Info("starting file server on: " + fileRoute)
		srvHandler.Static(fileRoute, handler.Config.Storage.DiskRoot)
	}
	// routes
	srvHandler.GET("/version", handler.GetVersion)
	srvHandler.GET("/disk", handler.GetDisk)
	version1 := srvHandler.Group("/v1")
	version1.GET("/features", handler.GetFeatures, getMiddlewareFuncs(handler.Config, handler.Cache, true)...)
	// search
	version1.GET("/search", handler.Search, getMiddlewareFuncs(handler.Config, handler.Cache, true)...)
	// mediaitems
	mediaItems := version1.Group("/mediaItems")
	mediaItems.GET("/:id/places", handler.GetMediaItemPlaces,
		getMiddlewareFuncs(handler.Config, handler.Cache, true, "places")...)
	mediaItems.GET("/:id/things", handler.GetMediaItemThings,
		getMiddlewareFuncs(handler.Config, handler.Cache, true, "things")...)
	mediaItems.GET("/:id/people", handler.GetMediaItemPeople,
		getMiddlewareFuncs(handler.Config, handler.Cache, true, "people")...)
	mediaItems.GET("/:id/albums", handler.GetMediaItemAlbums,
		getMiddlewareFuncs(handler.Config, handler.Cache, true, "albums")...)
	mediaItems.GET("/:id", handler.GetMediaItem,
		getMiddlewareFuncs(handler.Config, handler.Cache, true)...)
	mediaItems.PUT("/:id", handler.UpdateMediaItem,
		getMiddlewareFuncs(handler.Config, handler.Cache, true)...)
	mediaItems.DELETE("/:id", handler.DeleteMediaItem,
		getMiddlewareFuncs(handler.Config, handler.Cache, true)...)
	mediaItems.GET("", handler.GetMediaItems,
		getMiddlewareFuncs(handler.Config, handler.Cache, true)...)
	mediaItems.POST("", handler.UploadMediaItems,
		getMiddlewareFuncs(handler.Config, handler.Cache, true)...)
	// library
	favourites := version1.Group("/favourites")
	favourites.Use(getMiddlewareFuncs(handler.Config, handler.Cache, true, "favourites")...)
	favourites.GET("", handler.GetFavouriteMediaItems)
	favourites.POST("", handler.AddFavouriteMediaItems)
	favourites.DELETE("", handler.RemoveFavouriteMediaItems)
	hidden := version1.Group("/hidden")
	hidden.Use(getMiddlewareFuncs(handler.Config, handler.Cache, true, "hidden")...)
	hidden.GET("", handler.GetHiddenMediaItems)
	hidden.POST("", handler.AddHiddenMediaItems)
	hidden.DELETE("", handler.RemoveHiddenMediaItems)
	trash := version1.Group("/trash")
	trash.Use(getMiddlewareFuncs(handler.Config, handler.Cache, true, "trash")...)
	trash.GET("", handler.GetDeletedMediaItems)
	trash.POST("", handler.AddDeletedMediaItems)
	trash.DELETE("", handler.RemoveDeletedMediaItems)
	// explore
	explore := version1.Group("/explore")
	explore.Use(getMiddlewareFuncs(handler.Config, handler.Cache, true, "explore")...)
	explore.GET("/yearsAgo/:monthDate/mediaItems", handler.GetYearsAgoMediaItems)
	places := explore.Group("/places")
	places.Use(getMiddlewareFuncs(handler.Config, handler.Cache, true, "places")...)
	places.GET("/:id/mediaItems", handler.GetPlaceMediaItems)
	places.GET("/:id", handler.GetPlace)
	places.GET("", handler.GetPlaces)
	things := explore.Group("/things")
	things.Use(getMiddlewareFuncs(handler.Config, handler.Cache, true, "things")...)
	things.GET("/:id/mediaItems", handler.GetThingMediaItems)
	things.GET("/:id", handler.GetThing)
	things.GET("", handler.GetThings)
	people := explore.Group("/people")
	people.Use(getMiddlewareFuncs(handler.Config, handler.Cache, true, "people")...)
	people.GET("/:id/mediaItems", handler.GetPeopleMediaItems)
	people.GET("/:id", handler.GetPerson)
	people.PUT("/:id", handler.UpdatePerson)
	people.GET("", handler.GetPeople)
	// albums
	albums := version1.Group("/albums")
	albums.Use(getMiddlewareFuncs(handler.Config, handler.Cache, true, "albums")...)
	albums.GET("/:id/mediaItems", handler.GetAlbumMediaItems)
	albums.POST("/:id/mediaItems", handler.AddAlbumMediaItems)
	albums.DELETE("/:id/mediaItems", handler.RemoveAlbumMediaItems)
	albums.GET("/:id", handler.GetAlbum)
	albums.PUT("/:id", handler.UpdateAlbum)
	albums.DELETE("/:id", handler.DeleteAlbum)
	albums.GET("", handler.GetAlbums)
	albums.POST("", handler.CreateAlbum)
	// jobs
	jobs := version1.Group("/jobs")
	jobs.Use(getMiddlewareFuncs(handler.Config, handler.Cache, true, "jobs")...)
	jobs.GET("/:id", handler.GetJob)
	jobs.PUT("/:id", handler.UpdateJob)
	jobs.GET("", handler.GetJobs)
	jobs.POST("", handler.CreateJob)
	// sharing
	sharing := version1.Group("/sharing")
	sharing.Use(getMiddlewareFuncs(handler.Config, handler.Cache, false, "sharing")...)
	sharing.GET("/:id/mediaItems", handler.GetSharedAlbumMediaItems)
	sharing.GET("/:id", handler.GetSharedAlbum)
	// authentication
	auth := version1.Group("/auth")
	auth.POST("/login", handler.Login)
	auth.POST("/refresh", handler.Refresh)
	auth.POST("/logout", handler.Logout)
	// user management
	users := version1.Group("/users")
	users.Use(middlewares.BasicAuthCheck(handler.Config))
	users.GET("/:id", handler.GetUser)
	users.PUT("/:id", handler.UpdateUser)
	users.DELETE("/:id", handler.DeleteUser)
	users.GET("", handler.GetUsers)
	users.POST("", handler.CreateUser)

	go func() {
		slog.Info(fmt.Sprintf("starting http api server on: %d", handler.Config.API.Port))
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return httpServer
}

// StartHTTPServer ...
func StopHTTPServer(httpServer *http.Server) {
	slog.Info("stopping http api server")
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		panic(err)
	}
}

func getMiddlewareFuncs(cfg *config.Config, cache cache.Provider, jwtCheck bool, features ...string) []echo.MiddlewareFunc {
	middlewareFuncs := []echo.MiddlewareFunc{}
	if jwtCheck {
		middlewareFuncs = append(middlewareFuncs, middlewares.JWTCheck(cfg, cache))
	}
	for _, feature := range features {
		middlewareFuncs = append(middlewareFuncs, middlewares.FeatureCheck(cfg, feature))
	}
	return middlewareFuncs
}

func getFileRoute(storageDiskRoot string) string {
	fileRoute := strings.ReplaceAll(storageDiskRoot, "..", "")
	return fileRoute
}

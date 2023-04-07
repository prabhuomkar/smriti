package server

import (
	"api/config"
	"api/internal/handlers"
	"api/internal/middlewares"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bluele/gcache"
	"github.com/labstack/echo"
)

const httpTimeout = 10

// nolint:funlen
// StartHTTPServer ...
func StartHTTPServer(handler *handlers.Handler) *http.Server {
	srvHandler := echo.New()
	// nolint:gosec
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", handler.Config.API.Host, handler.Config.API.Port),
		Handler: srvHandler,
	}
	// routes
	srvHandler.GET("/version", handler.GetVersion)
	srvHandler.GET("/features", handler.GetFeatures)
	version1 := srvHandler.Group("/v1")
	// mediaitems
	mediaItems := version1.Group("/mediaItems")
	mediaItems.GET("/:id/places", handler.GetMediaItemPlaces,
		getMiddlewareFuncs(handler.Config, handler.Cache, "places")...)
	mediaItems.GET("/:id/things", handler.GetMediaItemThings,
		getMiddlewareFuncs(handler.Config, handler.Cache, "things")...)
	mediaItems.GET("/:id/people", handler.GetMediaItemPeople,
		getMiddlewareFuncs(handler.Config, handler.Cache, "people")...)
	mediaItems.GET("/:id/albums", handler.GetMediaItemAlbums,
		getMiddlewareFuncs(handler.Config, handler.Cache, "albums")...)
	mediaItems.GET("/:id", handler.GetMediaItem,
		getMiddlewareFuncs(handler.Config, handler.Cache)...)
	mediaItems.PUT("/:id", handler.UpdateMediaItem,
		getMiddlewareFuncs(handler.Config, handler.Cache)...)
	mediaItems.DELETE("/:id", handler.DeleteMediaItem,
		getMiddlewareFuncs(handler.Config, handler.Cache)...)
	mediaItems.GET("", handler.GetMediaItems,
		getMiddlewareFuncs(handler.Config, handler.Cache)...)
	mediaItems.POST("", handler.UploadMediaItems,
		getMiddlewareFuncs(handler.Config, handler.Cache)...)
	// library
	favourites := version1.Group("/favourites")
	favourites.Use(getMiddlewareFuncs(handler.Config, handler.Cache, "favourites")...)
	favourites.GET("", handler.GetFavouriteMediaItems)
	favourites.POST("", handler.AddFavouriteMediaItems)
	favourites.DELETE("", handler.RemoveFavouriteMediaItems)
	hidden := version1.Group("/hidden")
	hidden.Use(getMiddlewareFuncs(handler.Config, handler.Cache, "hidden")...)
	hidden.GET("", handler.GetHiddenMediaItems)
	hidden.POST("", handler.AddHiddenMediaItems)
	hidden.DELETE("", handler.RemoveHiddenMediaItems)
	trash := version1.Group("/trash")
	trash.Use(getMiddlewareFuncs(handler.Config, handler.Cache, "trash")...)
	trash.GET("", handler.GetDeletedMediaItems)
	trash.POST("", handler.AddDeletedMediaItems)
	trash.DELETE("", handler.RemoveDeletedMediaItems)
	// explore
	explore := version1.Group("/explore")
	explore.Use(middlewares.FeatureCheck(handler.Config, "explore"))
	places := explore.Group("/places")
	places.Use(getMiddlewareFuncs(handler.Config, handler.Cache, "places")...)
	places.GET("/:id/mediaItems", handler.GetPlaceMediaItems)
	places.GET("/:id", handler.GetPlace)
	places.GET("", handler.GetPlaces)
	things := explore.Group("/things")
	things.Use(getMiddlewareFuncs(handler.Config, handler.Cache, "things")...)
	things.GET("/:id/mediaItems", handler.GetThingMediaItems)
	things.GET("/:id", handler.GetThing)
	things.GET("", handler.GetThings)
	people := explore.Group("/people")
	people.Use(getMiddlewareFuncs(handler.Config, handler.Cache, "people")...)
	explore.GET("/:id/mediaItems", handler.GetPeopleMediaItems)
	explore.GET("/:id", handler.GetPerson)
	explore.PUT("/:id", handler.UpdatePerson)
	explore.GET("", handler.GetPeople)
	// albums
	albums := version1.Group("/albums")
	albums.Use(getMiddlewareFuncs(handler.Config, handler.Cache, "albums")...)
	albums.GET("/:id/mediaItems", handler.GetAlbumMediaItems)
	albums.POST("/:id/mediaItems", handler.AddAlbumMediaItems)
	albums.DELETE("/:id/mediaItems", handler.RemoveAlbumMediaItems)
	albums.GET("/:id", handler.GetAlbum)
	albums.PUT("/:id", handler.UpdateAlbum)
	albums.DELETE("/:id", handler.DeleteAlbum)
	albums.GET("", handler.GetAlbums)
	albums.POST("", handler.CreateAlbum)
	// authentication
	auth := version1.Group("/auth")
	auth.POST("/login", handler.Login)
	auth.POST("/refresh", handler.Refresh)
	auth.POST("/logout", handler.Logout)
	// user management
	users := version1.Group("/users")
	users.Use(middlewares.FeatureCheck(handler.Config, "users"), middlewares.BasicAuthCheck(handler.Config))
	users.GET("/:id", handler.GetUser)
	users.PUT("/:id", handler.UpdateUser)
	users.DELETE("/:id", handler.DeleteUser)
	users.GET("", handler.GetUsers)
	users.POST("", handler.CreateUser)

	go func() {
		log.Printf("starting http api server on: %d", handler.Config.API.Port)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return httpServer
}

// StartHTTPServer ...
func StopHTTPServer(httpServer *http.Server) {
	log.Println("stopping http api server")
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		panic(err)
	}
}

func getMiddlewareFuncs(cfg *config.Config, cache gcache.Cache, features ...string) []echo.MiddlewareFunc {
	middlewareFuncs := []echo.MiddlewareFunc{
		middlewares.JWTCheck(cfg, cache),
	}
	for _, feature := range features {
		middlewareFuncs = append(middlewareFuncs, middlewares.FeatureCheck(cfg, feature))
	}
	return middlewareFuncs
}

package server

import (
	"api/config"
	"api/internal/handlers"
	"api/internal/middlewares"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/bluele/gcache"
	"github.com/labstack/echo"
)

// nolint:funlen
// InitHTTPServer ...
func InitHTTPServer(handler *handlers.Handler) {
	srvHandler := echo.New()
	// nolint:gosec
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", handler.Config.API.Host, handler.Config.API.Port),
		Handler: srvHandler,
	}
	// routes
	// work(omkar): do this in a better way
	version1 := srvHandler.Group("/v1")
	version1.GET("/features", handler.GetFeatures)
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
	mediaItems.GET("/:id", handler.GetMediaItem)
	mediaItems.PUT("/:id", handler.UpdateMediaItem)
	mediaItems.DELETE("/:id", handler.DeleteMediaItem)
	mediaItems.GET("", handler.GetMediaItems)
	mediaItems.POST("", handler.UploadMediaItems,
		getMiddlewareFuncs(handler.Config, handler.Cache)...)
	// library
	version1.GET("/favourites", handler.GetFavouriteMediaItems,
		getMiddlewareFuncs(handler.Config, handler.Cache, "favourites")...)
	version1.POST("/favourites", handler.AddFavouriteMediaItems,
		getMiddlewareFuncs(handler.Config, handler.Cache, "favourites")...)
	version1.DELETE("/favourites", handler.RemoveFavouriteMediaItems,
		getMiddlewareFuncs(handler.Config, handler.Cache, "favourites")...)
	version1.GET("/hidden", handler.GetHiddenMediaItems,
		getMiddlewareFuncs(handler.Config, handler.Cache, "hidden")...)
	version1.POST("/hidden", handler.AddHiddenMediaItems,
		getMiddlewareFuncs(handler.Config, handler.Cache, "hidden")...)
	version1.DELETE("/hidden", handler.RemoveHiddenMediaItems,
		getMiddlewareFuncs(handler.Config, handler.Cache, "hidden")...)
	version1.GET("/trash", handler.GetDeletedMediaItems,
		getMiddlewareFuncs(handler.Config, handler.Cache, "trash")...)
	version1.POST("/trash", handler.AddDeletedMediaItems,
		getMiddlewareFuncs(handler.Config, handler.Cache, "trash")...)
	version1.DELETE("/trash", handler.RemoveDeletedMediaItems,
		getMiddlewareFuncs(handler.Config, handler.Cache, "trash")...)
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

	log.Printf("starting http api server on: %d", handler.Config.API.Port)
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
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

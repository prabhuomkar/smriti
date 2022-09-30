package server

import (
	"api/config"
	"api/internal/handlers"
	"api/internal/middlewares"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// nolint:funlen
// InitHTTPServer ...
func InitHTTPServer(cfg *config.Config, handler *handlers.Handler) {
	srvHandler := echo.New()
	// nolint:gosec
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port),
		Handler: srvHandler,
	}
	// routes
	// todo(omkar): do this in a better way
	version1 := srvHandler.Group("/v1")
	version1.GET("/features", handler.GetFeatures)
	// mediaitems
	mediaItems := version1.Group("/mediaItems")
	mediaItems.GET("/:id/places", handler.GetMediaItemPlaces, middlewares.FeatureCheck(cfg, "places"))
	mediaItems.GET("/:id/things", handler.GetMediaItemThings, middlewares.FeatureCheck(cfg, "things"))
	mediaItems.GET("/:id/people", handler.GetMediaItemPeople, middlewares.FeatureCheck(cfg, "people"))
	mediaItems.GET("/:id/albums", handler.GetMediaItemAlbums, middlewares.FeatureCheck(cfg, "albums"))
	mediaItems.GET("/:id", handler.GetMediaItem)
	mediaItems.PUT("/:id", handler.UpdateMediaItem)
	mediaItems.DELETE("/:id", handler.DeleteMediaItem)
	mediaItems.GET("", handler.GetMediaItems)
	mediaItems.POST("", handler.UploadMediaItems)
	// library
	version1.GET("/favourites", handler.GetFavouriteMediaItems, middlewares.FeatureCheck(cfg, "favourites"))
	version1.POST("/favourites", handler.AddFavouriteMediaItems, middlewares.FeatureCheck(cfg, "favourites"))
	version1.DELETE("/favourites", handler.RemoveFavouriteMediaItems, middlewares.FeatureCheck(cfg, "favourites"))
	version1.GET("/hidden", handler.GetHiddenMediaItems, middlewares.FeatureCheck(cfg, "hidden"))
	version1.POST("/hidden", handler.AddHiddenMediaItems, middlewares.FeatureCheck(cfg, "hidden"))
	version1.DELETE("/hidden", handler.RemoveHiddenMediaItems, middlewares.FeatureCheck(cfg, "hidden"))
	version1.GET("/trash", handler.GetDeletedMediaItems, middlewares.FeatureCheck(cfg, "trash"))
	version1.POST("/trash", handler.AddDeletedMediaItems, middlewares.FeatureCheck(cfg, "trash"))
	version1.DELETE("/trash", handler.RemoveDeletedMediaItems, middlewares.FeatureCheck(cfg, "trash"))
	// explore
	explore := version1.Group("/explore")
	explore.Use(middlewares.FeatureCheck(cfg, "explore"))
	explore.GET("/places/:id/mediaItems", handler.GetPlaceMediaItems, middlewares.FeatureCheck(cfg, "places"))
	explore.GET("/places/:id", handler.GetPlace, middlewares.FeatureCheck(cfg, "places"))
	explore.GET("/places", handler.GetPlaces, middlewares.FeatureCheck(cfg, "places"))
	explore.GET("/things/:id/mediaItems", handler.GetThingMediaItems, middlewares.FeatureCheck(cfg, "things"))
	explore.GET("/things/:id", handler.GetThing, middlewares.FeatureCheck(cfg, "things"))
	explore.GET("/things", handler.GetThings, middlewares.FeatureCheck(cfg, "things"))
	explore.GET("/people/:id/mediaItems", handler.GetPeopleMediaItems, middlewares.FeatureCheck(cfg, "people"))
	explore.GET("/people/:id", handler.GetPerson, middlewares.FeatureCheck(cfg, "people"))
	explore.PUT("/people/:id", handler.UpdatePerson, middlewares.FeatureCheck(cfg, "people"))
	explore.GET("/people", handler.GetPeople, middlewares.FeatureCheck(cfg, "people"))
	// albums
	albums := version1.Group("/albums")
	albums.Use(middlewares.FeatureCheck(cfg, "albums"))
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

	log.Printf("starting http api server on: %d", cfg.API.Port)
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

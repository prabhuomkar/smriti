package server

import (
	"api/config"
	"api/internal/handlers"
	"api/internal/middlewares"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// InitHTTPServer ...
func InitHTTPServer(cfg *config.Config, handler *handlers.Handler) {
	e := echo.New()
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port),
		Handler: e,
	}
	// routes
	// todo(omkar): do this in a better way
	v1 := e.Group("/v1")
	v1.GET("/features", handler.GetFeatures)
	// mediaitems
	mediaItems := v1.Group("/mediaItems")
	mediaItems.GET("/:id/places", handler.GetMediaItemPlaces, middlewares.IsFeatureEnabled(cfg, "places"))
	mediaItems.GET("/:id/things", handler.GetMediaItemThings, middlewares.IsFeatureEnabled(cfg, "things"))
	mediaItems.GET("/:id/people", handler.GetMediaItemPeople, middlewares.IsFeatureEnabled(cfg, "people"))
	mediaItems.GET("/:id", handler.GetMediaItem)
	mediaItems.PUT("/:id", handler.UpdateMediaItem)
	mediaItems.DELETE("/:id", handler.DeleteMediaItem)
	mediaItems.GET("", handler.GetMediaItems)
	mediaItems.POST("", handler.UploadMediaItems)
	// library
	v1.GET("/favourites", handler.GetFavouriteMediaItems, middlewares.IsFeatureEnabled(cfg, "favourites"))
	v1.GET("/hidden", handler.GetHiddenMediaItems, middlewares.IsFeatureEnabled(cfg, "hidden"))
	v1.GET("/trash", handler.GetDeletedMediaItems, middlewares.IsFeatureEnabled(cfg, "trash"))
	// explore
	explore := v1.Group("/explore")
	explore.Use(middlewares.IsFeatureEnabled(cfg, "explore"))
	explore.GET("/places/:id/mediaItems", handler.GetPlaceMediaItems, middlewares.IsFeatureEnabled(cfg, "places"))
	explore.GET("/places/:id", handler.GetPlace, middlewares.IsFeatureEnabled(cfg, "places"))
	explore.GET("/places", handler.GetPlaces, middlewares.IsFeatureEnabled(cfg, "places"))
	explore.GET("/things/:id/mediaItems", handler.GetThingMediaItems, middlewares.IsFeatureEnabled(cfg, "things"))
	explore.GET("/things/:id", handler.GetThing, middlewares.IsFeatureEnabled(cfg, "things"))
	explore.GET("/things", handler.GetThings, middlewares.IsFeatureEnabled(cfg, "things"))
	explore.GET("/people/:id/mediaItems", handler.GetPeopleMediaItems, middlewares.IsFeatureEnabled(cfg, "people"))
	explore.GET("/people/:id", handler.GetPerson, middlewares.IsFeatureEnabled(cfg, "people"))
	explore.GET("/people", handler.GetPeople, middlewares.IsFeatureEnabled(cfg, "people"))
	// albums
	albums := v1.Group("/albums")
	albums.Use(middlewares.IsFeatureEnabled(cfg, "albums"))
	albums.GET("/:id/mediaItems", handler.GetAlbumMediaItems)
	albums.POST("/:id/mediaItems", handler.AddAlbumMediaItems)
	albums.DELETE("/:id/mediaItems", handler.RemoveAlbumMediaItems)
	albums.GET("/:id", handler.GetAlbum)
	albums.PUT("/:id", handler.UpdateAlbum)
	albums.DELETE("/:id", handler.DeleteAlbum)
	albums.GET("", handler.GetAlbums)
	albums.POST("", handler.CreateAlbum)
	// authentication
	auth := v1.Group("/auth")
	auth.POST("/login", handler.Login)
	auth.POST("/refresh", handler.Refresh)
	auth.POST("/logout", handler.Logout)

	log.Printf("starting http api server on: %d", cfg.API.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}
}

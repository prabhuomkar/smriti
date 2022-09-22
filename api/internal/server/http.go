package server

import (
	"api/config"
	"api/internal/handlers"
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
	// todo(omkar): add middleware for feature check
	v1 := e.Group("/v1")
	// mediaitems
	mediaItems := v1.Group("/mediaItems")
	mediaItems.GET("/:id/places", handler.GetMediaItemPlaces)
	mediaItems.GET("/:id/things", handler.GetMediaItemThings)
	mediaItems.GET("/:id/people", handler.GetMediaItemPeople)
	mediaItems.GET("/:id", handler.GetMediaItem)
	mediaItems.PUT("/:id", handler.UpdateMediaItem)
	mediaItems.DELETE("/:id", handler.DeleteMediaItem)
	mediaItems.GET("", handler.GetMediaItems)
	mediaItems.POST("", handler.UploadMediaItems)
	// library
	v1.GET("/favourites", handler.GetFavouriteMediaItems)
	v1.GET("/hidden", handler.GetHiddenMediaItems)
	v1.GET("/trash", handler.GetDeletedMediaItems)
	// explore
	explore := v1.Group("/explore")
	explore.GET("/places/:placeId/mediaItems", handler.GetPlaceMediaItems)
	explore.GET("/places/:placeId", handler.GetPlace)
	explore.GET("/places", handler.GetPlaces)
	explore.GET("/things/:thingId/mediaItems", handler.GetThingMediaItems)
	explore.GET("/things/:thingId", handler.GetThing)
	explore.GET("/things", handler.GetThings)
	explore.GET("/people/:peopleId/mediaItems", handler.GetPeopleMediaItems)
	explore.GET("/people/:peopleId", handler.GetPerson)
	explore.GET("/people", handler.GetPeople)
	// albums
	albums := v1.Group("/albums")
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

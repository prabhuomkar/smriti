package handlers

import (
	"api/internal/models"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pgvector/pgvector-go"
)

const (
	minSearchQueryLen  = 3
	searchDefaultLimit = 50
)

// GetVersion ...
func (h *Handler) GetVersion(ctx echo.Context) error {
	version := models.GetVersion()

	return ctx.JSON(http.StatusOK, version)
}

// GetHealth ...
func (h *Handler) GetHealth(ctx echo.Context) error {
	health := models.GetHealth()

	health.Database.Status = models.StatusUp
	err := h.DB.Ping(ctx.Request().Context())
	if err != nil {
		health.Status = models.StatusDown
		health.Database.Status = models.StatusDown
		health.Database.Error = err.Error()
	}

	health.Cache.Status = models.StatusUp
	err = h.Cache.Ping()
	if err != nil {
		health.Status = models.StatusDown
		health.Cache.Status = models.StatusDown
		health.Cache.Error = err.Error()
	}

	if health.Status == models.StatusDown {
		healthBytes, _ := json.Marshal(health)
		return echo.NewHTTPError(http.StatusInternalServerError, string(healthBytes))
	}

	return ctx.JSON(http.StatusOK, health)
}

// GetFeatures ...
func (h *Handler) GetFeatures(ctx echo.Context) error {
	cfgFeatures := models.GetFeatures(h.Config)
	features, _ := ctx.Get("features").(models.Features)

	features.Favourites = features.Favourites && cfgFeatures.Favourites
	features.Hidden = features.Hidden && cfgFeatures.Hidden
	features.Trash = features.Trash && cfgFeatures.Trash
	features.Albums = features.Albums && cfgFeatures.Albums
	features.Explore = features.Explore && cfgFeatures.Explore
	features.Places = features.Places && cfgFeatures.Places
	features.People = features.People && cfgFeatures.People
	features.Sharing = features.Sharing && cfgFeatures.Sharing
	features.Jobs = features.Jobs && cfgFeatures.Jobs

	return ctx.JSON(http.StatusOK, features)
}

// GetDisk ...
func (h *Handler) GetDisk(ctx echo.Context) error {
	disk := models.GetDisk(h.Config)

	return ctx.JSON(http.StatusOK, disk)
}

// Search ...
func (h *Handler) Search(ctx echo.Context) error {
	searchQuery := ctx.QueryParam("q")
	if len(searchQuery) < minSearchQueryLen {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid search query")
	}
	mediaItems := []models.MediaItem{}
	if h.Config.Search {
		searchEmbedding := []float32{}
		var err error
		if err != nil {
			slog.Error("error getting search query embedding", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		rows, err := h.DB.Query(ctx.Request().Context(),
			"SELECT * FROM mediaitems WHERE id IN (SELECT id from mediaitem_embeddings ORDER BY embedding <-> $1) LIMIT $2",
			pgvector.NewVector(searchEmbedding), searchDefaultLimit)
		if err != nil {
			slog.Error("error searching mediaitems", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			mediaItem, err := models.ScanRowsToMediaItem(rows)
			if err != nil {
				slog.Error("error scanning album mediaitem", "error", err)

				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			mediaItems = append(mediaItems, mediaItem)
		}

		return ctx.JSON(http.StatusOK, mediaItems)
	}
	rows, err := h.DB.Query(ctx.Request().Context(),
		"SELECT * FROM mediaitems WHERE to_tsvector('english', caption) @@ plainto_tsquery('english', $1) LIMIT $2",
		searchQuery, searchDefaultLimit)
	if err != nil {
		slog.Error("error searching mediaitems", "error", err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		mediaItem, err := models.ScanRowsToMediaItem(rows)
		if err != nil {
			slog.Error("error scanning mediaitem", "error", err)

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		mediaItems = append(mediaItems, mediaItem)
	}

	return ctx.JSON(http.StatusOK, mediaItems)
}

package handlers

import (
	"api/internal/models"
	"api/pkg/services/worker"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pgvector/pgvector-go"
	"golang.org/x/exp/slog"
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
	features.Things = features.Things && cfgFeatures.Things
	features.People = features.People && cfgFeatures.People
	features.Sharing = features.Sharing && cfgFeatures.Sharing

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
	if h.Config.ML.Search {
		searchEmbedding, err := h.Worker.GenerateEmbedding(ctx.Request().Context(), &worker.GenerateEmbeddingRequest{Text: searchQuery})
		if err != nil {
			slog.Error("error getting search query embedding", slog.Any("error", err))
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		result := h.DB.Raw("SELECT * from mediaitems WHERE id IN (SELECT id from mediaitem_embeddings ORDER BY embedding <-> ?)", pgvector.NewVector(searchEmbedding.Embedding)).
			Find(&mediaItems).Limit(searchDefaultLimit)
		if result.Error != nil {
			slog.Error("error searching mediaitems", slog.Any("error", result.Error))
			return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
		}
		return ctx.JSON(http.StatusOK, mediaItems)
	}
	result := h.DB.Raw("SELECT * FROM mediaitems WHERE to_tsvector('english', keywords) @@ plainto_tsquery('english', ?)", searchQuery).
		Find(&mediaItems).Limit(searchDefaultLimit)
	if result.Error != nil {
		slog.Error("error searching mediaitems", slog.Any("error", result.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	return ctx.JSON(http.StatusOK, mediaItems)
}

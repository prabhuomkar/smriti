package middlewares

import (
	"api/config"
	"api/internal/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

// FeatureCheck ...
//
//nolint:cyclop
func FeatureCheck(cfg *config.Config, feature string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			features, _ := ctx.Get("features").(models.Features)
			if (feature == "favourites" && cfg.Feature.Favourites && features.Favourites) ||
				(feature == "hidden" && cfg.Feature.Hidden && features.Hidden) ||
				(feature == "trash" && cfg.Feature.Trash && features.Trash) ||
				(feature == "albums" && cfg.Feature.Albums && features.Albums) ||
				(feature == "explore" && cfg.Feature.Explore && features.Explore) ||
				(feature == "places" && cfg.Feature.Places && features.Places) ||
				(feature == "things" && cfg.Feature.Things && features.Things) ||
				(feature == "people" && cfg.Feature.People && features.People) ||
				(feature == "sharing" && cfg.Feature.Sharing) {
				return next(ctx)
			}
			slog.Error("feature disabled or not accessible", slog.Any("config", cfg.Feature), slog.Any("features", features))
			return echo.ErrForbidden
		}
	}
}

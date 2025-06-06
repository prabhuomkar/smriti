package middlewares

import (
	"api/config"
	"api/internal/models"
	"log/slog"

	"github.com/labstack/echo/v4"
)

// FeatureCheck ...
//
//nolint:cyclop
func FeatureCheck(cfg *config.Config, feature string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			features, _ := ctx.Get("features").(models.Features)
			if (feature == "favourites" && cfg.Favourites && features.Favourites) ||
				(feature == "hidden" && cfg.Hidden && features.Hidden) ||
				(feature == "trash" && cfg.Trash && features.Trash) ||
				(feature == "albums" && cfg.Albums && features.Albums) ||
				(feature == "explore" && cfg.Explore && features.Explore) ||
				(feature == "places" && cfg.Feature.Places && features.Places) ||
				(feature == "things" && cfg.Things && features.Things) ||
				(feature == "people" && cfg.People && features.People) ||
				(feature == "jobs" && cfg.Jobs && features.Jobs) ||
				(feature == "sharing" && cfg.Sharing) {
				return next(ctx)
			}
			slog.Error(
				"feature disabled or not accessible",
				"config",
				cfg.Feature,
				"features",
				features,
			)
			return echo.ErrForbidden
		}
	}
}

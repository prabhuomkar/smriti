package middlewares

import (
	"api/config"
	"api/internal/models"

	"github.com/labstack/echo"
)

// nolint:cyclop,gocognit
// FeatureCheck ...
func FeatureCheck(cfg *config.Config, feature string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			features, _ := ctx.Get("features").(models.Features)
			if (feature == "favourites" && cfg.Feature.Favourites && features.Favourites) ||
				(feature == "hidden" && cfg.Feature.Hidden && features.Hidden) ||
				(feature == "trash" && cfg.Feature.Trash && features.Trash) ||
				(feature == "albums" && cfg.Feature.Albums && features.Albums) ||
				(feature == "explore" && cfg.Feature.Explore && features.Explore) ||
				(feature == "places" && cfg.Feature.ExplorePlaces && features.Places) ||
				(feature == "things" && cfg.Feature.ExploreThings && features.Things) ||
				(feature == "people" && cfg.Feature.ExplorePeople && features.People) ||
				(feature == "sharing" && cfg.Feature.Sharing && features.Sharing) {
				return next(ctx)
			}
			return echo.ErrForbidden
		}
	}
}

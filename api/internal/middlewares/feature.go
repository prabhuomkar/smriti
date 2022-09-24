package middlewares

import (
	"api/config"

	"github.com/labstack/echo"
)

// FeatureCheck ...
func FeatureCheck(cfg *config.Config, feature string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// todo(omkar): do this in a better way
			if (feature == "favourites" && cfg.Feature.Favourites) ||
				(feature == "hidden" && cfg.Feature.Hidden) ||
				(feature == "trash" && cfg.Feature.Trash) ||
				(feature == "albums" && cfg.Feature.Albums) ||
				(feature == "explore" && cfg.Feature.Explore) ||
				(feature == "places" && cfg.Feature.ExplorePlaces) ||
				(feature == "things" && cfg.Feature.ExploreThings) ||
				(feature == "people" && cfg.Feature.ExplorePeople) ||
				(feature == "sharing" && cfg.Feature.Sharing) {
				return next(c)
			}

			return echo.ErrForbidden
		}
	}
}

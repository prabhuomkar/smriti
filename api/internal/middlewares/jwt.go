package middlewares

import (
	"api/config"
	"api/internal/auth"
	"strings"

	"github.com/bluele/gcache"
	"github.com/labstack/echo"
)

// JWTCheck ...
func JWTCheck(cfg *config.Config, cache gcache.Cache) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			accessToken := ctx.Request().Header.Get("Authorization")
			accessToken = strings.ReplaceAll(accessToken, "Bearer ", "")
			claims, err := auth.VerifyToken(cfg, cache, accessToken)
			if err == nil && claims != nil {
				ctx.Set("userID", claims.ID)
				ctx.Set("username", claims.Username)
				return next(ctx)
			}
			return echo.ErrForbidden
		}
	}
}

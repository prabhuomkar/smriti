package middlewares

import (
	"api/config"
	"api/internal/auth"
	"api/internal/models"
	"api/pkg/cache"
	"encoding/json"
	"strings"

	"github.com/labstack/echo/v4"
)

// JWTCheck ...
func JWTCheck(cfg *config.Config, cache cache.Provider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			accessToken := ctx.Request().Header.Get("Authorization")
			accessToken = strings.ReplaceAll(accessToken, "Bearer ", "")
			claims, err := auth.VerifyToken(cfg, cache, accessToken)
			if err == nil && claims != nil {
				ctx.Set("userID", claims.ID)
				ctx.Set("username", claims.Username)
				var features models.Features
				_ = json.Unmarshal([]byte(claims.Features), &features)
				ctx.Set("features", features)
				return next(ctx)
			}
			return echo.ErrUnauthorized
		}
	}
}

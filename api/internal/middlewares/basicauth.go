package middlewares

import (
	"api/config"

	"github.com/labstack/echo"
)

// BasicAuthCheck ...
func BasicAuthCheck(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			username, password, ok := ctx.Request().BasicAuth()
			if ok && username == cfg.Admin.Username && password == cfg.Admin.Password {
				return next(ctx)
			}
			return echo.ErrUnauthorized
		}
	}
}

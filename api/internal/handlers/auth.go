package handlers

import (
	"api/internal/auth"
	"api/internal/models"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type (
	// LoginRequest ...
	LoginRequest struct {
		Username *string `json:"username"`
		Password *string `json:"password"`
	}

	// AuthResponse ...
	AuthResponse struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
)

const (
	queryAuthLogin = `SELECT * FROM users WHERE username=$1 AND password=$2`
)

// Login ...
func (h *Handler) Login(ctx echo.Context) error {
	loginRequest, err := getUsernameAndPassword(ctx)
	if err != nil {
		return err
	}
	user := models.User{}
	err = h.DB.QueryRow(ctx.Request().Context(), queryAuthLogin, *loginRequest.Username, *loginRequest.Password).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Password,
		&user.Features,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "incorrect username or password")
		}
		slog.Error("error getting user", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	accessToken, refreshToken, err := auth.GetTokens(h.Config, h.Cache, user)
	if err != nil {
		slog.Error("error getting tokens", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting tokens")
	}
	authResponse := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return ctx.JSON(http.StatusOK, authResponse)
}

// Refresh ...
func (h *Handler) Refresh(ctx echo.Context) error {
	refreshToken := ctx.Request().Header.Get("Authorization")
	refreshToken = strings.ReplaceAll(refreshToken, "Bearer ", "")
	newAccessToken, newRefreshToken, err := auth.RefreshTokens(h.Config, h.Cache, refreshToken)
	if err != nil {
		slog.Error("error refreshing tokens", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error refreshing tokens")
	}
	authResponse := AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}
	return ctx.JSON(http.StatusOK, authResponse)
}

// Logout ...
func (h *Handler) Logout(ctx echo.Context) error {
	accessToken := ctx.Request().Header.Get("Authorization")
	accessToken = strings.ReplaceAll(accessToken, "Bearer ", "")
	_ = auth.RemoveTokens(h.Cache, accessToken)
	return ctx.JSON(http.StatusNoContent, nil)
}

func getUsernameAndPassword(ctx echo.Context) (*LoginRequest, error) {
	loginRequest := new(LoginRequest)
	err := ctx.Bind(loginRequest)
	if err != nil {
		slog.Error("error getting username and password", "error", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid username or password")
	}
	if loginRequest.Username == nil || loginRequest.Password == nil {
		slog.Error("error getting username and password", "error", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid username or password")
	}
	return loginRequest, nil
}

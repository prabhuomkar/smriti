package handlers

import (
	"api/internal/auth"
	"api/internal/models"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

// Login ...
func (h *Handler) Login(ctx echo.Context) error {
	loginRequest, err := getUsernameAndPassword(ctx)
	if err != nil {
		return err
	}
	user := models.User{}
	result := h.DB.Model(&models.User{}).
		Where("username=? AND password=?", &loginRequest.Username, &loginRequest.Password).
		First(&user)
	if result.Error != nil {
		log.Printf("error getting user: %+v", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "incorrect username or password")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
	}
	accessToken, refreshToken, err := auth.GetTokens(h.Config, h.Cache, user)
	if err != nil {
		log.Printf("error getting tokens: %+v", err)
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
		log.Printf("error refreshing tokens: %+v", err)
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
		log.Printf("error getting username and password: %+v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid username or password")
	}
	if loginRequest.Username == nil || loginRequest.Password == nil {
		log.Printf("error getting username and password: %+v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid username or password")
	}
	return loginRequest, nil
}

package auth

import (
	"api/config"
	"api/internal/models"
	"api/pkg/cache"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/exp/slog"
)

type (
	// TokenClaims ...
	TokenClaims struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Features string `json:"features"`
		jwt.RegisteredClaims
	}
)

// GetTokens ...
func GetTokens(cfg *config.Config, cache cache.Provider, user models.User) (string, string, error) {
	accessToken, refreshToken := GetAccessAndRefreshTokens(cfg, user)

	setRefreshErr := cache.SetWithExpire(refreshToken, true, time.Duration(cfg.Auth.RefreshTTL)*time.Second)
	if setRefreshErr != nil {
		slog.Error("error caching refresh token", "error", setRefreshErr)
		return "", "", setRefreshErr
	}
	setAccessErr := cache.SetWithExpire(accessToken, refreshToken, time.Duration(cfg.Auth.AccessTTL)*time.Second)
	if setAccessErr != nil {
		slog.Error("error caching refresh token", "error", setAccessErr)
		return "", "", setAccessErr
	}

	return accessToken, refreshToken, nil
}

// RefreshTokens ...
func RefreshTokens(cfg *config.Config, cache cache.Provider, refreshToken string) (string, string, error) {
	if _, err := cache.Get(refreshToken); err != nil {
		slog.Error("error getting refresh token from cache", "error", err)
		return "", "", err
	}

	claims, err := getClaimsFromToken(cfg, refreshToken)
	if err != nil {
		slog.Error("error getting claims from token", "error", err)
		return "", "", err
	}

	userID, err := uuid.FromString(claims.ID)
	if err != nil || userID == uuid.Nil {
		if err == nil {
			err = errors.New("got nil user id")
		}
		slog.Error("error getting user id from claims", "userId", userID, "error", err)
		return "", "", err
	}

	return GetTokens(cfg, cache, models.User{ID: userID, Username: claims.Username})
}

// RemoveTokens ...
func RemoveTokens(cache cache.Provider, accessToken string) error {
	refreshToken, err := cache.Get(accessToken)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			slog.Error("error getting access token from cache", "error", err)
		}
		return err
	}

	refreshTokenStr, _ := refreshToken.(string)
	_ = cache.Remove(refreshTokenStr)
	_ = cache.Remove(accessToken)

	return nil
}

// VerifyToken ...
func VerifyToken(cfg *config.Config, cache cache.Provider, accessToken string) (*TokenClaims, error) {
	_, err := cache.Get(accessToken)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			slog.Error("error getting access token from cache", "error", err)
		}
		return nil, err
	}

	claims, err := getClaimsFromToken(cfg, accessToken)
	if err != nil {
		slog.Error("error getting claims from token", "error", err)
		return nil, err
	}

	return claims, nil
}

func GetAccessAndRefreshTokens(cfg *config.Config, user models.User) (string, string) {
	return getSignedToken(cfg, user, "access"), getSignedToken(cfg, user, "refresh")
}

func getClaimsFromToken(cfg *config.Config, token string) (*TokenClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(cfg.Auth.Secret), nil
	})
	if err != nil || !parsedToken.Valid {
		slog.Error("error parsing claims from token", "error", err)
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*TokenClaims)
	if !ok {
		slog.Error("error getting claims from token", "error", err)
		return nil, err
	}

	return claims, nil
}

func getSignedToken(cfg *config.Config, user models.User, subject string) string {
	ttl := cfg.Auth.AccessTTL
	if subject == "refresh" {
		ttl = cfg.Auth.RefreshTTL
	}
	creationTime := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		user.ID.String(),
		user.Username,
		user.Features,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(creationTime.Add(time.Duration(ttl) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(creationTime),
			NotBefore: jwt.NewNumericDate(creationTime),
			Issuer:    cfg.Auth.Issuer,
			Audience:  []string{cfg.Auth.Audience},
			Subject:   subject,
			ID:        user.ID.String(),
		},
	})
	signedToken, _ := token.SignedString([]byte(cfg.Auth.Secret))
	return signedToken
}

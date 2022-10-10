package auth

import (
	"api/config"
	"api/internal/models"
	"log"
	"time"

	"github.com/bluele/gcache"
	"github.com/golang-jwt/jwt/v4"
)

type (
	// TokenClaims ...
	TokenClaims struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		jwt.RegisteredClaims
	}
)

// GetTokens ...
func GetTokens(cfg *config.Config, cache gcache.Cache, user models.User) (string, string, error) {
	accessToken, refreshToken, err := getAccessAndRefreshTokens(cfg, user.ID.String(), user.Username)
	if err != nil {
		log.Printf("error getting tokens: %+v", err)
		return "", "", err
	}

	_ = cache.Set(refreshToken, true)
	_ = cache.Set(accessToken, refreshToken)

	return accessToken, refreshToken, err
}

// RefreshTokens ...
func RefreshTokens(cfg *config.Config, cache gcache.Cache, refreshToken string) (string, string, error) {
	_, err := cache.Get(refreshToken)
	if err != nil {
		log.Printf("error getting refresh token from cache: %+v", err)
		return "", "", err
	}

	claims, err := getClaimsFromToken(cfg, refreshToken)
	if err != nil {
		log.Printf("error getting claims from token: %+v", err)
		return "", "", err
	}

	newAccessToken, newRefreshToken, err := getAccessAndRefreshTokens(cfg, claims.ID, claims.Username)
	if err != nil {
		log.Printf("error getting tokens: %+v", err)
		return "", "", err
	}

	_ = cache.Set(newRefreshToken, true)
	_ = cache.Set(newAccessToken, newRefreshToken)

	return newAccessToken, newRefreshToken, err
}

// RemoveTokens ...
func RemoveTokens(cfg *config.Config, cache gcache.Cache, accessToken string) error {
	refreshToken, err := cache.Get(accessToken)
	if err != nil {
		log.Printf("error getting access token from cache: %+v", err)
		return err
	}

	_ = cache.Remove(refreshToken)
	_ = cache.Remove(accessToken)

	return nil
}

func getAccessAndRefreshTokens(cfg *config.Config, userID, username string) (string, string, error) {
	accessToken, err := getSignedToken(cfg, userID, username, "access")
	if err != nil {
		log.Printf("error generating access token: %+v", err)
		return "", "", err
	}

	refreshToken, err := getSignedToken(cfg, userID, username, "refresh")
	if err != nil {
		log.Printf("error generating refresh token: %+v", err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func getClaimsFromToken(cfg *config.Config, token string) (*TokenClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Auth.Secret), nil
	})
	if err != nil || !parsedToken.Valid {
		log.Printf("error parsing claims from token: %+v", err)
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*TokenClaims)
	if !ok {
		log.Printf("error getting claims from token: %+v", err)
		return nil, err
	}

	return claims, nil
}

func getSignedToken(cfg *config.Config, userID, username, subject string) (string, error) {
	ttl := cfg.Auth.AccessTTL
	if subject == "refresh" {
		ttl = cfg.Auth.RefreshTTL
	}
	creationTime := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		userID,
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(creationTime.Add(time.Duration(ttl) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(creationTime),
			NotBefore: jwt.NewNumericDate(creationTime),
			Issuer:    cfg.Auth.Issuer,
			Audience:  []string{cfg.Auth.Audience},
			Subject:   subject,
			ID:        userID,
		},
	})
	return token.SignedString([]byte(cfg.Auth.Secret))
}

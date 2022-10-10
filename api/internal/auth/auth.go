package auth

import (
	"api/internal/models"

	"github.com/bluele/gcache"
)

// GetTokens ...
func GetTokens(cache gcache.Cache, user models.User) (string, string, error) {
	return "", "", nil
}

// RefreshTokens ...
func RefreshTokens(cache gcache.Cache, refreshToken string) (string, string, error) {
	return "", "", nil
}

// RemoveTokens ...
func RemoveTokens(cache gcache.Cache, accessToken string) error {
	return nil
}

package auth

import (
	"api/config"
	"api/internal/models"
	"log"
	"time"

	"github.com/bluele/gcache"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
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
	accessToken, refreshToken := GetAccessAndRefreshTokens(cfg, user.ID.String(), user.Username)

	setRefreshErr := cache.SetWithExpire(refreshToken, true, time.Duration(cfg.Auth.RefreshTTL)*time.Second)
	if setRefreshErr != nil {
		log.Printf("error caching refresh token: %+v", setRefreshErr)
		return "", "", setRefreshErr
	}
	setAccessErr := cache.SetWithExpire(accessToken, refreshToken, time.Duration(cfg.Auth.AccessTTL)*time.Second)
	if setAccessErr != nil {
		log.Printf("error caching refresh token: %+v", setAccessErr)
		return "", "", setAccessErr
	}

	return accessToken, refreshToken, nil
}

// RefreshTokens ...
func RefreshTokens(cfg *config.Config, cache gcache.Cache, refreshToken string) (string, string, error) {
	if _, err := cache.Get(refreshToken); err != nil {
		log.Printf("error getting refresh token from cache: %+v", err)
		return "", "", err
	}

	claims, err := getClaimsFromToken(cfg, refreshToken)
	if err != nil {
		log.Printf("error getting claims from token: %+v", err)
		return "", "", err
	}

	userID, err := uuid.FromString(claims.ID)
	if err != nil {
		log.Printf("error converting user id from claims: %+v", err)
		return "", "", err
	}

	return GetTokens(cfg, cache, models.User{ID: userID, Username: claims.Username})
}

// RemoveTokens ...
func RemoveTokens(cache gcache.Cache, accessToken string) error {
	refreshToken, err := cache.Get(accessToken)
	if err != nil {
		log.Printf("error getting access token from cache: %+v", err)
		return err
	}

	_ = cache.Remove(refreshToken)
	_ = cache.Remove(accessToken)

	return nil
}

// VerifyToken ...
func VerifyToken(cfg *config.Config, cache gcache.Cache, accessToken string) (*TokenClaims, error) {
	if _, err := cache.Get(accessToken); err != nil {
		log.Printf("error getting access token from cache: %+v", err)
		return nil, err
	}

	claims, err := getClaimsFromToken(cfg, accessToken)
	if err != nil {
		log.Printf("error getting claims from token: %+v", err)
		return nil, err
	}

	return claims, nil
}

func GetAccessAndRefreshTokens(cfg *config.Config, userID, username string) (string, string) {
	return getSignedToken(cfg, userID, username, "access"), getSignedToken(cfg, userID, username, "refresh")
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

func getSignedToken(cfg *config.Config, userID, username, subject string) string {
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
	signedToken, _ := token.SignedString([]byte(cfg.Auth.Secret))
	return signedToken
}

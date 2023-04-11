package auth

import (
	"api/config"
	"api/internal/models"
	"errors"
	"testing"

	"github.com/bluele/gcache"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetTokens(t *testing.T) {
	tests := []struct {
		Name          string
		SerializeFunc func(interface{}, interface{}) (interface{}, error)
		WantErr       bool
	}{
		{
			"success",
			nil,
			false,
		},
		{
			"error caching refresh token",
			func(a interface{}, b interface{}) (interface{}, error) {
				val, ok := b.(bool)
				if ok && val == true {
					return nil, errors.New("some cache error")
				}
				return b, nil
			},
			true,
		},
		{
			"error caching access token",
			func(a interface{}, b interface{}) (interface{}, error) {
				val, ok := b.(bool)
				if ok && val == true {
					return b, nil
				}
				return nil, errors.New("some cache error")
			},
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			cfg := &config.Config{}
			cache := gcache.New(1024).LRU().SerializeFunc(test.SerializeFunc).Build()
			atoken, rtoken, err := GetTokens(cfg, cache, models.User{ID: uuid.FromStringOrNil("4d05b5f6-17c2-475e-87fe-3fc8b9567179")})
			if test.WantErr {
				assert.Empty(t, atoken)
				assert.Empty(t, rtoken)
				assert.NotNil(t, err)
			} else {
				assert.NotEmpty(t, atoken)
				assert.NotEmpty(t, rtoken)
				assert.Nil(t, err)
			}
		})
	}
}

func TestRefreshTokens(t *testing.T) {
	tests := []struct {
		Name            string
		Token           func(*config.Config, gcache.Cache) string
		SerializeFunc   func(interface{}, interface{}) (interface{}, error)
		DeserializeFunc func(interface{}, interface{}) (interface{}, error)
		WantErr         bool
	}{
		{
			"success",
			func(cfg *config.Config, cache gcache.Cache) string {
				_, oldRToken := GetAccessAndRefreshTokens(cfg, models.User{ID: uuid.FromStringOrNil("4d05b5f6-17c2-475e-87fe-3fc8b9567179"), Username: "username"})
				_ = cache.Set(oldRToken, true)
				return oldRToken
			},
			nil,
			nil,
			false,
		},
		{
			"error getting refresh token",
			func(cfg *config.Config, cache gcache.Cache) string {
				return "badToken"
			},
			nil,
			nil,
			true,
		},
		{
			"error parsing claims from token",
			func(cfg *config.Config, cache gcache.Cache) string {
				_ = cache.Set("badToken", true)
				return "badToken"
			},
			nil,
			nil,
			true,
		},
		{
			"error converting user id from claims",
			func(cfg *config.Config, cache gcache.Cache) string {
				_, oldRToken := GetAccessAndRefreshTokens(cfg, models.User{ID: uuid.FromStringOrNil("invalid-user-id"), Username: "username"})
				_ = cache.Set(oldRToken, true)
				return oldRToken
			},
			nil,
			nil,
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			cfg := &config.Config{Auth: config.Auth{
				RefreshTTL: 60,
			}}
			cache := gcache.New(1024).LRU().
				SerializeFunc(test.SerializeFunc).
				DeserializeFunc(test.DeserializeFunc).
				Build()
			oldRToken := test.Token(cfg, cache)
			atoken, rtoken, err := RefreshTokens(cfg, cache, oldRToken)
			if test.WantErr {
				assert.Empty(t, atoken)
				assert.Empty(t, rtoken)
				assert.NotNil(t, err)
			} else {
				assert.NotEmpty(t, atoken)
				assert.NotEmpty(t, rtoken)
				assert.Nil(t, err)
			}
		})
	}
}

func TestRemoveTokens(t *testing.T) {
	tests := []struct {
		Name            string
		Token           func(*config.Config, gcache.Cache) string
		SerializeFunc   func(interface{}, interface{}) (interface{}, error)
		DeserializeFunc func(interface{}, interface{}) (interface{}, error)
		WantErr         bool
	}{
		{
			"success",
			func(cfg *config.Config, cache gcache.Cache) string {
				oldAToken, _ := GetAccessAndRefreshTokens(cfg, models.User{ID: uuid.FromStringOrNil("4d05b5f6-17c2-475e-87fe-3fc8b9567179"), Username: "username"})
				_ = cache.Set(oldAToken, true)
				return oldAToken
			},
			nil,
			nil,
			false,
		},
		{
			"error getting refresh token",
			func(cfg *config.Config, cache gcache.Cache) string {
				return "badToken"
			},
			nil,
			nil,
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			cfg := &config.Config{Auth: config.Auth{
				RefreshTTL: 60,
			}}
			cache := gcache.New(1024).LRU().
				SerializeFunc(test.SerializeFunc).
				DeserializeFunc(test.DeserializeFunc).
				Build()
			oldAToken := test.Token(cfg, cache)
			err := RemoveTokens(cache, oldAToken)
			if test.WantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	tests := []struct {
		Name            string
		Token           func(*config.Config, gcache.Cache) string
		SerializeFunc   func(interface{}, interface{}) (interface{}, error)
		DeserializeFunc func(interface{}, interface{}) (interface{}, error)
		WantErr         bool
	}{
		{
			"success",
			func(cfg *config.Config, cache gcache.Cache) string {
				oldAToken, _ := GetAccessAndRefreshTokens(cfg, models.User{ID: uuid.FromStringOrNil("4d05b5f6-17c2-475e-87fe-3fc8b9567179"), Username: "username"})
				_ = cache.Set(oldAToken, true)
				return oldAToken
			},
			nil,
			nil,
			false,
		},
		{
			"error getting access token",
			func(cfg *config.Config, cache gcache.Cache) string {
				return "badToken"
			},
			nil,
			nil,
			true,
		},
		{
			"error parsing claims from token",
			func(cfg *config.Config, cache gcache.Cache) string {
				_ = cache.Set("badToken", true)
				return "badToken"
			},
			nil,
			nil,
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			cfg := &config.Config{Auth: config.Auth{
				AccessTTL: 60,
			}}
			cache := gcache.New(1024).LRU().
				SerializeFunc(test.SerializeFunc).
				DeserializeFunc(test.DeserializeFunc).
				Build()
			oldAToken := test.Token(cfg, cache)
			claims, err := VerifyToken(cfg, cache, oldAToken)
			if test.WantErr {
				assert.Nil(t, claims)
				assert.NotNil(t, err)
			} else {
				assert.NotNil(t, claims)
				assert.Nil(t, err)
			}
		})
	}
}

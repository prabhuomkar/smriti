package auth

import (
	"api/config"
	"api/internal/models"
	"errors"
	"testing"

	"github.com/bluele/gcache"
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
			atoken, rtoken, err := GetTokens(cfg, cache, models.User{})
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
				_, oldRToken := GetAccessAndRefreshTokens(cfg, "userID", "username")
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
			"error caching refresh token",
			func(cfg *config.Config, cache gcache.Cache) string {
				_, oldRToken := GetAccessAndRefreshTokens(cfg, "userID", "username")
				_ = cache.Set(oldRToken, false)
				return oldRToken
			},
			func(a interface{}, b interface{}) (interface{}, error) {
				val, ok := b.(bool)
				if ok && val == true {
					return nil, errors.New("some cache error")
				}
				return b, nil
			},
			nil,
			true,
		},
		{
			"error caching access token",
			func(cfg *config.Config, cache gcache.Cache) string {
				_, oldRToken := GetAccessAndRefreshTokens(cfg, "userID", "username")
				_ = cache.Set(oldRToken, true)
				return oldRToken
			},
			func(a interface{}, b interface{}) (interface{}, error) {
				val, ok := b.(bool)
				if ok && val == true {
					return b, nil
				}
				return nil, errors.New("some cache error")
			},
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
				oldAToken, _ := GetAccessAndRefreshTokens(cfg, "userID", "username")
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
			err := RemoveTokens(cfg, cache, oldAToken)
			if test.WantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

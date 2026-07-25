package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		Name     string
		MockFunc func() func()
		WantErr  bool
	}{
		{
			"error due to invalid feature boolean value", func() func() {
				os.Setenv("SMRITI_FEATURE_FAVOURITES", "invalid")
				return func() {
					os.Setenv("SMRITI_FEATURE_FAVOURITES", "")
				}
			}, true,
		},
		{
			"success", func() func() {
				os.Setenv("SMRITI_FEATURE_FAVOURITES", "true")
				return func() {
					os.Setenv("SMRITI_FEATURE_FAVOURITES", "")
				}
			}, false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if test.MockFunc != nil {
				defer test.MockFunc()()
			}
			cfg, err := Init()
			if test.WantErr {
				assert.Nil(t, cfg)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, cfg)
				assert.NoError(t, err)
			}
		})
	}
}

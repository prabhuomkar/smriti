package models

import (
	"api/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	trueVal  = true
	falseVal = false
)

func TestGetFeatures(t *testing.T) {
	features := GetFeatures(&config.Config{})
	assert.Equal(t, &Features{
		Albums:     falseVal,
		Favourites: falseVal,
		Hidden:     falseVal,
		Trash:      falseVal,
		Explore:    falseVal,
		Places:     falseVal,
		People:     falseVal,
		Things:     falseVal,
		Sharing:    falseVal,
		ML: &MLFeatures{
			Places:         falseVal,
			Classification: falseVal,
			Detection:      falseVal,
			Faces:          falseVal,
			OCR:            falseVal,
			Speech:         falseVal,
		},
	}, features)

	features = GetFeatures(&config.Config{
		Feature: config.Feature{
			Albums:     true,
			Favourites: true,
			Hidden:     true,
			Trash:      true,
			Explore:    true,
			Places:     true,
			People:     true,
			Things:     true,
			Sharing:    true,
			ML: config.MLFeatures{
				Places:         true,
				Classification: true,
				Detection:      true,
				Faces:          true,
				OCR:            true,
				Speech:         true,
			},
		},
	})
	assert.Equal(t, &Features{
		Albums:     trueVal,
		Favourites: trueVal,
		Hidden:     trueVal,
		Trash:      trueVal,
		Explore:    trueVal,
		Places:     trueVal,
		People:     trueVal,
		Things:     trueVal,
		Sharing:    trueVal,
		ML: &MLFeatures{
			Places:         trueVal,
			Classification: trueVal,
			Detection:      trueVal,
			Faces:          trueVal,
			OCR:            trueVal,
			Speech:         trueVal,
		},
	}, features)
}

package models

import "api/config"

type (
	// Features ...
	Features struct {
		Favourites *bool `json:"favourites,omitempty"`
		Hidden     *bool `json:"hidden,omitempty"`
		Trash      *bool `json:"trash,omitempty"`
		Albums     *bool `json:"albums,omitempty"`
		Explore    *bool `json:"explore,omitempty"`
		Places     *bool `json:"places,omitempty"`
		Things     *bool `json:"things,omitempty"`
		People     *bool `json:"people,omitempty"`
		Sharing    *bool `json:"sharing,omitempty"`
	}
)

// GetFeatures ...
func GetFeatures(cfg *config.Config) *Features {
	return &Features{
		Favourites: &cfg.Feature.Favourites,
		Hidden:     &cfg.Feature.Hidden,
		Trash:      &cfg.Feature.Trash,
		Albums:     &cfg.Feature.Albums,
		Explore:    &cfg.Feature.Explore,
		Places:     &cfg.Feature.ExplorePlaces,
		Things:     &cfg.Feature.ExploreThings,
		People:     &cfg.Feature.ExplorePeople,
		Sharing:    &cfg.Feature.Sharing,
	}
}

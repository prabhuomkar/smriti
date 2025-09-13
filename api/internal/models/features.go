package models

import "api/config"

type ( // Features ...
	Features struct {
		Favourites bool `json:"favourites,omitempty"`
		Hidden     bool `json:"hidden,omitempty"`
		Trash      bool `json:"trash,omitempty"`
		Albums     bool `json:"albums,omitempty"`
		Explore    bool `json:"explore,omitempty"`
		Places     bool `json:"places,omitempty"`
		Things     bool `json:"things,omitempty"`
		People     bool `json:"people,omitempty"`
		Sharing    bool `json:"sharing,omitempty"`
		Jobs       bool `json:"jobs,omitempty"`
	}
)

// GetFeatures ...
func GetFeatures(cfg *config.Config) *Features {
	return &Features{
		Favourites: cfg.Favourites, Hidden: cfg.Hidden, Trash: cfg.Trash,
		Albums: cfg.Albums, Explore: cfg.Explore, Places: cfg.Feature.Places,
		Things: cfg.Things, People: cfg.People, Sharing: cfg.Sharing,
	}
}

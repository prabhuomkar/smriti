package models

import "api/config"

type (
	// MLFeatures ...
	MLFeatures struct {
		Places         bool `json:"places,omitempty"`
		Classification bool `json:"classification,omitempty"`
		Detection      bool `json:"detection,omitempty"`
		Faces          bool `json:"faces,omitempty"`
		OCR            bool `json:"ocr,omitempty"`
		Speech         bool `json:"speech,omitempty"`
	}

	// Features ...
	Features struct {
		Favourites bool        `json:"favourites,omitempty"`
		Hidden     bool        `json:"hidden,omitempty"`
		Trash      bool        `json:"trash,omitempty"`
		Albums     bool        `json:"albums,omitempty"`
		Explore    bool        `json:"explore,omitempty"`
		Places     bool        `json:"places,omitempty"`
		Things     bool        `json:"things,omitempty"`
		People     bool        `json:"people,omitempty"`
		Sharing    bool        `json:"sharing,omitempty"`
		ML         *MLFeatures `json:"ml,omitempty"`
	}
)

// GetMLFeaturesList ...
func (m *MLFeatures) GetMLFeaturesList() []string {
	result := []string{}
	if m.Places {
		result = append(result, "places")
	}
	if m.Classification {
		result = append(result, "classification")
	}
	if m.Detection {
		result = append(result, "detection")
	}
	if m.Faces {
		result = append(result, "faces")
	}
	if m.OCR {
		result = append(result, "ocr")
	}
	if m.Speech {
		result = append(result, "speech")
	}
	return result
}

// GetFeatures ...
func GetFeatures(cfg *config.Config) *Features {
	return &Features{
		Favourites: cfg.Feature.Favourites,
		Hidden:     cfg.Feature.Hidden,
		Trash:      cfg.Feature.Trash,
		Albums:     cfg.Feature.Albums,
		Explore:    cfg.Feature.Explore,
		Places:     cfg.Feature.Places,
		Things:     cfg.Feature.Things,
		People:     cfg.Feature.People,
		Sharing:    cfg.Feature.Sharing,
		ML: &MLFeatures{
			Places:         cfg.Feature.ML.Places,
			Classification: cfg.Feature.ML.Classification,
			Detection:      cfg.Feature.ML.Detection,
			Faces:          cfg.Feature.ML.Faces,
			OCR:            cfg.Feature.ML.OCR,
			Speech:         cfg.Feature.ML.Speech,
		},
	}
}

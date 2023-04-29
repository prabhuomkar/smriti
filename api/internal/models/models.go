package models

// GetModels ...
func GetModels() []interface{} {
	return []interface{}{
		User{},
		Album{},
		Place{},
		Thing{},
		People{},
		MediaItem{},
	}
}

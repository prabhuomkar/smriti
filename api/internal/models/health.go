package models

//nolint:gochecknoglobals
var (
	DefaultStatus = "UP"
	StatusUp      = "UP"
	StatusDown    = "DOWN"
)

type (
	// ResourceHealth ...
	ResourceHealth struct {
		Status string `json:"status"`
		Error  string `json:"error,omitempty"`
	}

	// Health ...
	Health struct {
		Status   string         `json:"status"`
		Database ResourceHealth `json:"database"`
		Cache    ResourceHealth `json:"cache"`
	}
)

// GetHealth ...
func GetHealth() *Health {
	return &Health{Status: DefaultStatus}
}

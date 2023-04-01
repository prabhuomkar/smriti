package models

var (
	DefaultVersion = "dev"
	DefaultGitSHA  = "-"
)

type (
	// Version ...
	Version struct {
		Version string `json:"version"`
		GitSHA  string `json:"gitSHA"`
	}
)

// GetVersion ...
func GetVersion() *Version {
	return &Version{
		Version: DefaultVersion,
		GitSHA:  DefaultGitSHA,
	}
}

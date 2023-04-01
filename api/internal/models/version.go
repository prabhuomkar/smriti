package models

// nolint:gochecknoglobals
var (
	DefaultVersion = "dev"
	DefaultGitSHA  = "-"
)

type (
	// Version ...
	Version struct {
		Version string `json:"version"`
		GitSHA  string `json:"gitSha"`
	}
)

// GetVersion ...
func GetVersion() *Version {
	return &Version{
		Version: DefaultVersion,
		GitSHA:  DefaultGitSHA,
	}
}

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	ver := GetVersion()
	assert.Equal(t, &Version{Version: DefaultVersion, GitSHA: DefaultGitSHA}, ver)
}

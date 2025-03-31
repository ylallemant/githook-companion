package binary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInfo(t *testing.T) {
	info := GetInfo()
	assert.Contains(t, info, semverVersion)
	assert.Contains(t, info, gitCommitHash)
}

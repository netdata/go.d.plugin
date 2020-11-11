package du

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.Equal(t, len(New().Config.Paths), 0)
}

func TestDu_Init(t *testing.T) {
	assert.False(t, New().Init())
	du := &Du{
		Config: Config{
			Paths: []string{
				"/tmp/dummy",
			},
		},
		collectedDims: make(map[string]bool),
	}
	assert.True(t, du.Init())
}

func TestDu_Collect(t *testing.T) {
	du := &Du{
		Config: Config{
			Paths: []string{
				"/tmp/dummy",
			},
		},
		collectedDims: make(map[string]bool),
	}
	du.Init()
	assert.Equal(t, du.Collect()["/tmp/dummy"], int64(-1))
}

func TestFileSize(t *testing.T) {
	size, _ := fileSize("testdata/file.txt")
	assert.Greater(t, size, int64(0))
}

func TestValidateConfig(t *testing.T) {
	du := Du{
		Config: Config{
			Paths: []string{
				"/tmp/fileA",
				"/tmp/fileB",
			},
		},
		collectedDims: make(map[string]bool),
	}

	assert.Nil(t, du.validateConfig())
}

func TestDu_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

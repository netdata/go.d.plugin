package modules

import (
	"github.com/netdata/go.d.plugin/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase_Init(t *testing.T) {
	assert.True(t, (&Base{}).Init())
}

func TestBase(t *testing.T) {
	base := Base{}

	assert.Implements(t, (*interface {
		Init() bool
		SetLogger(logger *logger.Logger)
		Cleanup()
	})(nil), &base)

}

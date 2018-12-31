package activemq

import (
	"testing"

	"github.com/netdata/go.d.plugin/modules"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	mod := New()

	assert.Implements(t, (*modules.Module)(nil), mod)
	assert.Equal(t, defURL, mod.URL)
	assert.Equal(t, defHTTPTimeout, mod.Client.Timeout.Duration)
	assert.Equal(t, defMaxQueues, mod.MaxQueues)
	assert.Equal(t, defMaxTopics, mod.MaxTopics)
}

func TestActivemq_Init(t *testing.T) {
	mod := New()

	// NG case
	assert.False(t, mod.Init())

	// OK case
	mod.Webadmin = "webadmin"
	assert.True(t, mod.Init())
	assert.NotNil(t, mod.reqQueues)
	assert.NotNil(t, mod.reqTopics)
	assert.NotNil(t, mod.client)

}

func TestActivemq_Check(t *testing.T) {

}

func TestActivemq_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestActivemq_Cleanup(t *testing.T) {

}

func TestActivemq_Collect(t *testing.T) {

}

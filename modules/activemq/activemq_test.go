package activemq

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestActivemq_Collect_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	mod := New()
	mod.Webadmin = "webadmin"
	mod.HTTP.Request = web.Request{URL: ts.URL}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

func TestActivemq_Collect_InvalidData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello and goodbye!"))
	}))
	defer ts.Close()

	mod := New()
	mod.Webadmin = "webadmin"
	mod.HTTP.Request = web.Request{URL: ts.URL}

	require.True(t, mod.Init())
	assert.False(t, mod.Check())
}

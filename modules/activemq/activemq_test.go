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

var (
	testQueues = `<queues>
<queue name="sandra">
<stats size="0" consumerCount="0" enqueueCount="0" dequeueCount="0"/>
<feed>
<atom>queueBrowse/sandra?view=rss&amp;feedType=atom_1.0</atom>
<rss>queueBrowse/sandra?view=rss&amp;feedType=rss_2.0</rss>
</feed>
</queue>
<queue name="Test">
<stats size="0" consumerCount="0" enqueueCount="0" dequeueCount="0"/>
<feed>
<atom>queueBrowse/Test?view=rss&amp;feedType=atom_1.0</atom>
<rss>queueBrowse/Test?view=rss&amp;feedType=rss_2.0</rss>
</feed>
</queue>
</queues>`

	testTopics = `<topics>
<topic name="ActiveMQ.Advisory.MasterBroker ">
<stats size="0" consumerCount="0" enqueueCount="1" dequeueCount="0"/>
</topic>
<topic name="AAA ">
<stats size="0" consumerCount="1" enqueueCount="0" dequeueCount="0"/>
</topic>
<topic name="ActiveMQ.Advisory.Topic ">
<stats size="0" consumerCount="0" enqueueCount="1" dequeueCount="0"/>
</topic>
<topic name="ActiveMQ.Advisory.Queue ">
<stats size="0" consumerCount="0" enqueueCount="2" dequeueCount="0"/>
</topic>
<topic name="AAAA ">
<stats size="0" consumerCount="0" enqueueCount="0" dequeueCount="0"/>
</topic>
</topics>`
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
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/webadmin/xml/queues.jsp":
					_, _ = w.Write([]byte(testQueues))
				case "/webadmin/xml/topics.jsp":
					_, _ = w.Write([]byte(testTopics))
				}
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL}
	mod.Webadmin = "webadmin"

	require.True(t, mod.Init())
	require.True(t, mod.Check())
}

func TestActivemq_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestActivemq_Cleanup(t *testing.T) {
	New().Cleanup()
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

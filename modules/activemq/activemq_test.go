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
	queuesData = []string{
		`<queues>
<queue name="sandra">
<stats size="1" consumerCount="1" enqueueCount="2" dequeueCount="1"/>
<feed>
<atom>queueBrowse/sandra?view=rss&amp;feedType=atom_1.0</atom>
<rss>queueBrowse/sandra?view=rss&amp;feedType=rss_2.0</rss>
</feed>
</queue>
<queue name="Test">
<stats size="1" consumerCount="1" enqueueCount="2" dequeueCount="1"/>
<feed>
<atom>queueBrowse/Test?view=rss&amp;feedType=atom_1.0</atom>
<rss>queueBrowse/Test?view=rss&amp;feedType=rss_2.0</rss>
</feed>
</queue>
</queues>`,
		`<queues>
<queue name="sandra">
<stats size="2" consumerCount="2" enqueueCount="3" dequeueCount="2"/>
<feed>
<atom>queueBrowse/sandra?view=rss&amp;feedType=atom_1.0</atom>
<rss>queueBrowse/sandra?view=rss&amp;feedType=rss_2.0</rss>
</feed>
</queue>
<queue name="Test">
<stats size="2" consumerCount="2" enqueueCount="3" dequeueCount="2"/>
<feed>
<atom>queueBrowse/Test?view=rss&amp;feedType=atom_1.0</atom>
<rss>queueBrowse/Test?view=rss&amp;feedType=rss_2.0</rss>
</feed>
</queue>
<queue name="Test2">
<stats size="0" consumerCount="0" enqueueCount="0" dequeueCount="0"/>
<feed>
<atom>queueBrowse/Test?view=rss&amp;feedType=atom_1.0</atom>
<rss>queueBrowse/Test?view=rss&amp;feedType=rss_2.0</rss>
</feed>
</queue>
</queues>`,
		`<queues>
<queue name="sandra">
<stats size="3" consumerCount="3" enqueueCount="4" dequeueCount="3"/>
<feed>
<atom>queueBrowse/sandra?view=rss&amp;feedType=atom_1.0</atom>
<rss>queueBrowse/sandra?view=rss&amp;feedType=rss_2.0</rss>
</feed>
</queue>
<queue name="Test">
<stats size="3" consumerCount="3" enqueueCount="4" dequeueCount="3"/>
<feed>
<atom>queueBrowse/Test?view=rss&amp;feedType=atom_1.0</atom>
<rss>queueBrowse/Test?view=rss&amp;feedType=rss_2.0</rss>
</feed>
</queue>
</queues>`,
	}

	topicsData = []string{
		`<topics>
<topic name="ActiveMQ.Advisory.MasterBroker ">
<stats size="0" consumerCount="0" enqueueCount="1" dequeueCount="0"/>
</topic>
<topic name="AAA ">
<stats size="1" consumerCount="1" enqueueCount="2" dequeueCount="1"/>
</topic>
<topic name="ActiveMQ.Advisory.Topic ">
<stats size="0" consumerCount="0" enqueueCount="1" dequeueCount="0"/>
</topic>
<topic name="ActiveMQ.Advisory.Queue ">
<stats size="0" consumerCount="0" enqueueCount="2" dequeueCount="0"/>
</topic>
<topic name="AAAA ">
<stats size="1" consumerCount="1" enqueueCount="2" dequeueCount="1"/>
</topic>
</topics>`,
		`<topics>
<topic name="ActiveMQ.Advisory.MasterBroker ">
<stats size="0" consumerCount="0" enqueueCount="1" dequeueCount="0"/>
</topic>
<topic name="AAA ">
<stats size="2" consumerCount="2" enqueueCount="3" dequeueCount="2"/>
</topic>
<topic name="ActiveMQ.Advisory.Topic ">
<stats size="0" consumerCount="0" enqueueCount="1" dequeueCount="0"/>
</topic>
<topic name="ActiveMQ.Advisory.Queue ">
<stats size="0" consumerCount="0" enqueueCount="2" dequeueCount="0"/>
</topic>
<topic name="AAAA ">
<stats size="2" consumerCount="2" enqueueCount="3" dequeueCount="2"/>
</topic>
<topic name="BBB ">
<stats size="1" consumerCount="1" enqueueCount="2" dequeueCount="1"/>
</topic>
</topics>`,
		`<topics>
<topic name="ActiveMQ.Advisory.MasterBroker ">
<stats size="0" consumerCount="0" enqueueCount="1" dequeueCount="0"/>
</topic>
<topic name="AAA ">
<stats size="3" consumerCount="3" enqueueCount="4" dequeueCount="3"/>
</topic>
<topic name="ActiveMQ.Advisory.Topic ">
<stats size="0" consumerCount="0" enqueueCount="1" dequeueCount="0"/>
</topic>
<topic name="ActiveMQ.Advisory.Queue ">
<stats size="0" consumerCount="0" enqueueCount="2" dequeueCount="0"/>
</topic>
<topic name="AAAA ">
<stats size="3" consumerCount="3" enqueueCount="4" dequeueCount="3"/>
</topic>
</topics>`,
	}
)

func TestNew(t *testing.T) {
	mod := New()

	assert.Implements(t, (*modules.Module)(nil), mod)
	assert.Equal(t, defaultURL, mod.URL)
	assert.Equal(t, defaultHTTPTimeout, mod.Client.Timeout.Duration)
	assert.Equal(t, defaultMaxQueues, mod.MaxQueues)
	assert.Equal(t, defaultMaxTopics, mod.MaxTopics)
}

func TestActivemq_Init(t *testing.T) {
	mod := New()

	// NG case
	assert.False(t, mod.Init())

	// OK case
	mod.Webadmin = "webadmin"
	assert.True(t, mod.Init())
	assert.NotNil(t, mod.apiClient)
}

func TestActivemq_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/webadmin/xml/queues.jsp":
					_, _ = w.Write([]byte(queuesData[0]))
				case "/webadmin/xml/topics.jsp":
					_, _ = w.Write([]byte(topicsData[0]))
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
	var collectNum int
	getQueues := func() string { return queuesData[collectNum] }
	getTopics := func() string { return topicsData[collectNum] }

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/webadmin/xml/queues.jsp":
					_, _ = w.Write([]byte(getQueues()))
				case "/webadmin/xml/topics.jsp":
					_, _ = w.Write([]byte(getTopics()))
				}
			}))
	defer ts.Close()

	mod := New()
	mod.HTTP.Request = web.Request{URL: ts.URL}
	mod.Webadmin = "webadmin"

	require.True(t, mod.Init())
	require.True(t, mod.Check())

	cases := []struct {
		expected  map[string]int64
		numQueues int
		numTopics int
		numCharts int
	}{
		{
			expected: map[string]int64{
				"queues_sandra_consumers":   1,
				"queues_sandra_dequeued":    1,
				"queues_Test_enqueued":      2,
				"queues_Test_unprocessed":   1,
				"topics_AAA_dequeued":       1,
				"topics_AAAA_unprocessed":   1,
				"queues_Test_dequeued":      1,
				"topics_AAA_enqueued":       2,
				"topics_AAA_unprocessed":    1,
				"topics_AAAA_consumers":     1,
				"topics_AAAA_dequeued":      1,
				"queues_Test_consumers":     1,
				"queues_sandra_enqueued":    2,
				"queues_sandra_unprocessed": 1,
				"topics_AAA_consumers":      1,
				"topics_AAAA_enqueued":      2,
			},
			numQueues: 2,
			numTopics: 2,
			numCharts: 12,
		},
		{
			expected: map[string]int64{
				"queues_sandra_enqueued":    3,
				"queues_Test_enqueued":      3,
				"queues_Test_unprocessed":   1,
				"queues_Test2_dequeued":     0,
				"topics_BBB_enqueued":       2,
				"queues_sandra_dequeued":    2,
				"queues_sandra_unprocessed": 1,
				"queues_Test2_enqueued":     0,
				"topics_AAAA_enqueued":      3,
				"topics_AAAA_dequeued":      2,
				"topics_BBB_unprocessed":    1,
				"topics_AAA_dequeued":       2,
				"topics_AAAA_unprocessed":   1,
				"queues_Test_consumers":     2,
				"queues_Test_dequeued":      2,
				"queues_Test2_consumers":    0,
				"queues_Test2_unprocessed":  0,
				"topics_AAA_consumers":      2,
				"topics_AAA_enqueued":       3,
				"topics_BBB_dequeued":       1,
				"queues_sandra_consumers":   2,
				"topics_AAA_unprocessed":    1,
				"topics_AAAA_consumers":     2,
				"topics_BBB_consumers":      1,
			},
			numQueues: 3,
			numTopics: 3,
			numCharts: 18,
		},
		{
			expected: map[string]int64{
				"queues_sandra_unprocessed": 1,
				"queues_Test_unprocessed":   1,
				"queues_sandra_consumers":   3,
				"topics_AAAA_enqueued":      4,
				"queues_sandra_dequeued":    3,
				"queues_Test_consumers":     3,
				"queues_Test_enqueued":      4,
				"queues_Test_dequeued":      3,
				"topics_AAA_consumers":      3,
				"topics_AAA_unprocessed":    1,
				"topics_AAAA_consumers":     3,
				"topics_AAAA_unprocessed":   1,
				"queues_sandra_enqueued":    4,
				"topics_AAA_enqueued":       4,
				"topics_AAA_dequeued":       3,
				"topics_AAAA_dequeued":      3,
			},
			numQueues: 2,
			numTopics: 2,
			numCharts: 18,
		},
	}

	for _, c := range cases {
		require.Equal(t, c.expected, mod.Collect())
		assert.Len(t, mod.activeQueues, c.numQueues)
		assert.Len(t, mod.activeTopics, c.numTopics)
		assert.Len(t, *mod.charts, c.numCharts)
		collectNum++
	}
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

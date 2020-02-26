package vernemq

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	metricsV1100MQTT5, _ = ioutil.ReadFile("testdata/metrics-v1.10.0-mqtt5.txt")
)

func Test_readTestData(t *testing.T) {
	assert.NotNil(t, metricsV1100MQTT5)
}

func TestVerneMQ_Collect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(metricsV1100MQTT5)
		}))

	verneMQ := New()
	verneMQ.UserURL = ts.URL
	require.True(t, verneMQ.Init())

	m := verneMQ.Collect()
	l := make([]string, 0)
	for k := range m {
		l = append(l, k)
	}
	sort.Strings(l)
	for i, value := range l {
		//if !strings.Contains(value, "disconnect") {
		//	continue
		//}
		fmt.Println(fmt.Sprintf("%d \"%s\": %d,", i+1, value, m[value]))
	}
}

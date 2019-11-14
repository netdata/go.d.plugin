package unbound

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
)

var (
	statsData, _    = ioutil.ReadFile("testdata/stats.txt")
	extStatsData, _ = ioutil.ReadFile("testdata/extended_stats.txt")
)

func Test_readTestData(t *testing.T) {
	assert.NotNil(t, statsData)
	assert.NotNil(t, extStatsData)
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestUnbound_Init(t *testing.T) {

}

func TestUnbound_Check(t *testing.T) {

}

func TestUnbound_Cleanup(t *testing.T) {

}

func TestUnbound_Charts(t *testing.T) {

}

func TestUnbound_Collect(t *testing.T) {
	v := New()
	v.Init()
	v.client = newClient(clientConfig{
		address: "192.168.88.223:8953",
		timeout: time.Second * 2,
		useTLS:  false,
		tlsConf: nil,
	})
	m := v.Collect()

	for k, v := range m {
		fmt.Println(k, v)
	}
}

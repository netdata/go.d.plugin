package dnsquery

import (
	"errors"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestDNSQuery_Init(t *testing.T) {
	mod := New()

	// NG case
	assert.False(t, mod.Init())

	// OK case
	mod.Domains = []string{"google.com"}
	mod.Servers = []string{"8.8.8.8", "8.8.4.4"}
	require.True(t, mod.Init())
	assert.Len(t, mod.servers, len(mod.Servers))
	assert.Len(t, mod.workers, len(mod.Servers))
}

func TestDNSQuery_Check(t *testing.T) {
	assert.True(t, New().Check())
}

func TestDNSQuery_Charts(t *testing.T) {
	mod := New()

	assert.NotNil(t, mod.Charts())

	mod.Domains = []string{"google.com"}
	mod.Servers = []string{"8.8.8.8"}
	require.True(t, mod.Init())
	charts := mod.Charts()
	assert.True(t, charts.Get("query_time").HasDim("8_8_8_8"))
}

func TestDNSQuery_Cleanup(t *testing.T) {
	mod := New()

	mod.Domains = []string{"google.com"}
	mod.Servers = []string{"8.8.8.8"}
	require.True(t, mod.Init())

	time.Sleep(time.Second)
	require.Len(t, mod.servers, len(mod.Servers))
	require.Len(t, mod.workers, len(mod.Servers))

	mod.Cleanup()
	time.Sleep(time.Second)
	assert.Len(t, mod.workers, 0)

	wait := time.NewTimer(time.Second)
	defer wait.Stop()

	select {
	case <-wait.C:
		t.Error("cleanup failed, task channel is not closed")
	case <-mod.task:
	}
}

func TestDNSQuery_Collect(t *testing.T) {
	mod := New()
	defer mod.Cleanup()

	mod.Domains = []string{"google.com"}
	mod.Servers = []string{"8.8.8.8"}
	mod.exchangerFactory = func(network string, duration time.Duration) exchanger {
		return okMockExchanger{}
	}

	require.True(t, mod.Init())
	require.True(t, mod.Check())

	assert.Equal(
		t,
		map[string]int64{"8_8_8_8": 1000000000},
		mod.Collect(),
	)
}

func TestDNSQuery_Collect_Error(t *testing.T) {
	mod := New()
	defer mod.Cleanup()

	mod.Domains = []string{"google.com"}
	mod.Servers = []string{"8.8.8.8"}
	mod.exchangerFactory = func(network string, duration time.Duration) exchanger {
		return errMockExchanger{}
	}

	require.True(t, mod.Init())
	require.True(t, mod.Check())

	assert.Len(
		t,
		mod.Collect(),
		0,
	)
}

type okMockExchanger struct{}

func (m okMockExchanger) Exchange(msg *dns.Msg, address string) (response *dns.Msg, rtt time.Duration, err error) {
	return nil, time.Second, nil
}

type errMockExchanger struct{}

func (m errMockExchanger) Exchange(msg *dns.Msg, address string) (response *dns.Msg, rtt time.Duration, err error) {
	return nil, time.Second, errors.New("mock error")
}

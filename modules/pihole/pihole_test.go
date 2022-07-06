// SPDX-License-Identifier: GPL-3.0-or-later

package pihole

import (
	"errors"
	"testing"

	"github.com/netdata/go.d.plugin/modules/pihole/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSetupVarsPathOK    = "testdata/setupVars.conf"
	testSetupVarsPathWrong = "testdata/wrong.conf"
	testWebPassword        = "1ebd33f882f9aa5fac26a7cb74704742f91100228eb322e41b7bd6e6aeb8f74b"
)

func newTestJob() *Pihole {
	job := New()
	job.SetupVarsPath = testSetupVarsPathOK

	return job
}

func TestNew(t *testing.T) {
	job := New()

	assert.IsType(t, (*Pihole)(nil), job)
	assert.Equal(t, defaultURL, job.URL)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
	assert.Equal(t, defaultSetupVarsPath, job.SetupVarsPath)
	assert.Equal(t, defaultTopClients, job.TopClientsEntries)
	assert.Equal(t, defaultTopItems, job.TopItemsEntries)
}

func TestPihole_Init(t *testing.T) {
	job := newTestJob()

	assert.True(t, job.Init())
	assert.Equal(t, job.Password, testWebPassword)
	assert.NotNil(t, job.client)

	job = newTestJob()

	job.SetupVarsPath = testSetupVarsPathWrong
	assert.True(t, job.Init())
}

func TestPihole_Check(t *testing.T) {
	job := newTestJob()

	assert.True(t, job.Init())

	job.client = newOKTestPiholeClient()
	assert.True(t, job.Check())
}

func TestPihole_Check_WrongVersion(t *testing.T) {
	job := newTestJob()

	require.True(t, job.Init())

	job.client = &testPiholeAPIClient{
		version: func() (i int, e error) { return supportedAPIVersion + 1, nil },
	}
	assert.False(t, job.Check())
}

func TestPihole_Check_NoData(t *testing.T) {
	job := newTestJob()

	require.True(t, job.Init())

	job.client = &testPiholeAPIClient{
		version: func() (i int, e error) { return supportedAPIVersion, nil },
	}
	assert.False(t, job.Check())
}

func TestPihole_Charts(t *testing.T) {
	job := newTestJob()

	require.True(t, job.Init())
	job.client = newOKTestPiholeClient()
	require.True(t, job.Check())
	assert.Len(t, *job.Charts(), len(charts)+len(authCharts)+3) // 3 top* charts, added during check

	job = newTestJob()

	job.SetupVarsPath = testSetupVarsPathWrong
	require.True(t, job.Init())
	job.client = newOKTestPiholeClient()
	require.True(t, job.Check())
	assert.Len(t, *job.Charts(), len(charts))
}

func TestPihole_Cleanup(t *testing.T) { assert.NotPanics(t, newTestJob().Cleanup) }

func TestPihole_Collect(t *testing.T) {
	job := newTestJob()

	require.True(t, job.Init())
	job.client = newOKTestPiholeClient()
	require.True(t, job.Check())
	require.NotNil(t, job.Charts())

	expected := map[string]int64{
		"A":                    0,
		"AAAA":                 0,
		"ANY":                  0,
		"PTR":                  0,
		"SOA":                  0,
		"SRV":                  0,
		"TXT":                  0,
		"ads_blocked_today":    0,
		"ads_percentage_today": 0,
		// "blocklist_last_update": 1561019970,
		"destination_d1":        3329,
		"destination_d2":        6659,
		"dns_queries_today":     0,
		"domains_being_blocked": 0,
		"file_exists":           0,
		"queries_cached":        0,
		"queries_forwarded":     0,
		"status":                0,
		"top_blocked_domain_a1": 33,
		"top_blocked_domain_a2": 66,
		"top_client_c1":         33,
		"top_client_c2":         66,
		"top_perm_domain_q1":    33,
		"top_perm_domain_q2":    66,
		"unique_clients":        0,
	}

	//collected := job.Collect()
	// expected["blocklist_last_update"] = collected["blocklist_last_update"]

	assert.Equal(t, expected, job.Collect())
}

func TestPihole_Collect_OnlySummary(t *testing.T) {
	job := newTestJob()

	require.True(t, job.Init())

	c := newOKTestPiholeClient()
	c.queryTypes = nil
	c.forwardDest = nil
	c.topClients = nil
	c.topItems = nil
	job.client = c

	require.True(t, job.Check())
	require.NotNil(t, job.Charts())

	expected := map[string]int64{
		"ads_blocked_today":    0,
		"ads_percentage_today": 0,
		// "blocklist_last_update": 1561019970,
		"dns_queries_today":     0,
		"domains_being_blocked": 0,
		"file_exists":           0,
		"queries_cached":        0,
		"queries_forwarded":     0,
		"status":                0,
		"unique_clients":        0,
	}

	//collected := job.Collect()
	//expected["blocklist_last_update"] = collected["blocklist_last_update"]

	assert.Equal(t, expected, job.Collect())
}

func TestPihole_Collect_NoData(t *testing.T) {
	job := newTestJob()

	require.True(t, job.Init())

	job.client = newOKTestPiholeClient()
	require.True(t, job.Check())
	require.NotNil(t, job.Charts())

	job.client = &testPiholeAPIClient{}
	assert.Nil(t, job.Collect())
}

func newOKTestPiholeClient() *testPiholeAPIClient {
	return &testPiholeAPIClient{
		version: func() (int, error) {
			return supportedAPIVersion, nil
		},
		summary: func() (*client.SummaryRaw, error) {
			return &client.SummaryRaw{}, nil
		},
		queryTypes: func() (*client.QueryTypes, error) {
			return &client.QueryTypes{}, nil
		},
		forwardDest: func() (*[]client.ForwardDestination, error) {
			return &[]client.ForwardDestination{
				{Name: "d1", Percent: 33.3},
				{Name: "d2", Percent: 66.6},
			}, nil
		},
		topClients: func() (*[]client.TopClient, error) {
			return &[]client.TopClient{
				{Name: "c1", Requests: 33},
				{Name: "c2", Requests: 66},
			}, nil
		},
		topItems: func() (*client.TopItems, error) {
			return &client.TopItems{
				TopQueries: []client.TopQuery{
					{Name: "q1", Hits: 33},
					{Name: "q2", Hits: 66},
				},
				TopAds: []client.TopAdvertisement{
					{Name: "a1", Hits: 33},
					{Name: "a2", Hits: 66},
				},
			}, nil
		},
	}
}

type testPiholeAPIClient struct {
	version     func() (int, error)
	summary     func() (*client.SummaryRaw, error)
	queryTypes  func() (*client.QueryTypes, error)
	forwardDest func() (*[]client.ForwardDestination, error)
	topClients  func() (*[]client.TopClient, error)
	topItems    func() (*client.TopItems, error)
}

func (t testPiholeAPIClient) Version() (int, error) {
	if t.version == nil {
		return 0, errors.New("version is <nil>")
	}
	return t.version()
}

func (t testPiholeAPIClient) SummaryRaw() (*client.SummaryRaw, error) {
	if t.summary == nil {
		return nil, errors.New("summary is <nil>")
	}
	return t.summary()
}

func (t testPiholeAPIClient) QueryTypes() (*client.QueryTypes, error) {
	if t.queryTypes == nil {
		return nil, errors.New("queryTypes is <nil>")
	}
	return t.queryTypes()
}

func (t testPiholeAPIClient) ForwardDestinations() (*[]client.ForwardDestination, error) {
	if t.forwardDest == nil {
		return nil, errors.New("forwardDest is <nil>")
	}
	return t.forwardDest()
}

func (t testPiholeAPIClient) TopClients(_ int) (*[]client.TopClient, error) {
	if t.topClients == nil {
		return nil, errors.New("topClients is <nil>")
	}
	return t.topClients()
}

func (t testPiholeAPIClient) TopItems(_ int) (*client.TopItems, error) {
	if t.topItems == nil {
		return nil, errors.New("topItems is <nil>")
	}
	return t.topItems()
}

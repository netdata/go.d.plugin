package pihole

import (
	"errors"
	"testing"

	"github.com/netdata/go.d.plugin/modules/pihole/client"

	"github.com/stretchr/testify/assert"
)

const (
	testSetupVarsPath      = "testdata/setupVars.conf"
	testSetupVarsPathWrong = "testdata/wrong.conf"
	testWebPassword        = "1ebd33f882f9aa5fac26a7cb74704742f91100228eb322e41b7bd6e6aeb8f74b"
)

func TestNew(t *testing.T) {
	job := New()

	assert.IsType(t, (*Pihole)(nil), job)
	assert.Equal(t, defaultURL, job.UserURL)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
	assert.Equal(t, defaultSetupVarsPath, job.SetupVarsPath)
	assert.Equal(t, defaultTopClients, job.TopClientsEntries)
	assert.Equal(t, defaultTopItems, job.TopItemsEntries)
}

func TestPihole_Init(t *testing.T) {
	job := New()
	job.SetupVarsPath = testSetupVarsPath

	assert.True(t, job.Init())
	assert.Equal(t, job.Password, testWebPassword)
	assert.NotNil(t, job.client)

	job = New()
	job.SetupVarsPath = testSetupVarsPathWrong
	assert.True(t, job.Init())
}

func TestPihole_Check(t *testing.T) {
	job := New()
	assert.True(t, job.Init())
	job.client = newOKTestPiholeClient()

	assert.True(t, job.Check())
}

func TestPihole_Charts(t *testing.T) {

}

func TestPihole_Cleanup(t *testing.T) {

}

func TestPihole_Collect(t *testing.T) {

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
			return &[]client.ForwardDestination{}, nil
		},
		topClients: func() (*[]client.TopClient, error) {
			return &[]client.TopClient{}, nil
		},
		topItems: func() (*client.TopItems, error) {
			return &client.TopItems{}, nil
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

func (t testPiholeAPIClient) TopClients(top int) (*[]client.TopClient, error) {
	if t.topClients == nil {
		return nil, errors.New("topClients is <nil>")
	}
	return t.topClients()
}

func (t testPiholeAPIClient) TopItems(top int) (*client.TopItems, error) {
	if t.topItems == nil {
		return nil, errors.New("topItems is <nil>")
	}
	return t.topItems()
}

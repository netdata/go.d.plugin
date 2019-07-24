package client

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testWebPassword = "12345678"

var (
	testEmptyData                  = []byte("[]")
	testVersionData                = []byte(`{"version": 3}`)
	testSummaryRawData, _          = ioutil.ReadFile("testdata/summaryRaw.json")
	testQueryTypesData, _          = ioutil.ReadFile("testdata/getQueryTypes.json")
	testForwardDestinationsData, _ = ioutil.ReadFile("testdata/getForwardDestinations.json")
	testTopClientsData, _          = ioutil.ReadFile("testdata/topClients.json")
	testTopItemsData, _            = ioutil.ReadFile("testdata/topItems.json")
)

func Test_data(t *testing.T) {
	assert.NotEmpty(t, testSummaryRawData)
	assert.NotEmpty(t, testQueryTypesData)
	assert.NotEmpty(t, testForwardDestinationsData)
	assert.NotEmpty(t, testTopClientsData)
	assert.NotEmpty(t, testTopItemsData)
}

func TestNew(t *testing.T) {
	client := New(Configuration{})
	assert.NotNil(t, client)
	assert.NotNil(t, client.HTTPClient)
	assert.Empty(t, client.URL)
	assert.Empty(t, client.WebPassword)

	client = New(Configuration{URL: "url", WebPassword: "pass"})
	assert.NotNil(t, client)
	assert.NotNil(t, client.HTTPClient)
	assert.Equal(t, "url", client.URL)
	assert.Equal(t, "pass", client.WebPassword)

}

func TestClient_Version(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	client := New(Configuration{})
	client.URL = ts.URL

	ver, err := client.Version()
	require.NoError(t, err)

	expected := 3

	assert.Equal(t, expected, ver)
}

func TestClient_SummaryRaw(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	client := New(Configuration{})
	client.URL = ts.URL

	rv, err := client.SummaryRaw()
	require.NoError(t, err)

	var absolute int64 = 1560443834

	expected := &SummaryRaw{
		DomainsBeingBlocked: 1,
		DNSQueriesToday:     1,
		AdsBlockedToday:     1,
		AdsPercentageToday:  1,
		UniqueDomains:       1,
		QueriesForwarded:    1,
		QueriesCached:       1,
		ClientsEverSeen:     1,
		UniqueClients:       1,
		DNSQueriesAllTypes:  1,
		ReplyNODATA:         1,
		ReplyNXDOMAIN:       1,
		ReplyCNAME:          1,
		ReplyIP:             1,
		PrivacyLevel:        1,
		Status:              "enabled",
		GravityLastUpdated: struct {
			FileExists bool `json:"file_exists"`
			Absolute   *int64
		}{FileExists: true, Absolute: &absolute},
	}

	assert.Equal(t, expected, rv)
}

func TestClient_QueryTypes(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	client := New(Configuration{})
	client.URL = ts.URL
	client.WebPassword = testWebPassword

	rv, err := client.QueryTypes()
	require.NoError(t, err)

	expected := &QueryTypes{
		A:    12.29,
		AAAA: 12.29,
		ANY:  1,
		SRV:  1,
		SOA:  1,
		PTR:  71.43,
		TXT:  1,
	}

	assert.Equal(t, expected, rv)
}

func TestClient_ForwardDestinations(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	client := New(Configuration{})
	client.URL = ts.URL
	client.WebPassword = testWebPassword

	rv, err := client.ForwardDestinations()
	require.NoError(t, err)

	sort.Slice(*rv, func(i, j int) bool {
		return (*rv)[i].Percent < (*rv)[j].Percent
	})

	expected := &[]ForwardDestination{
		{Name: "blocklist", Percent: 0},
		{Name: "resolver1.opendns.com", Percent: 2.78},
		{Name: "resolver2.opendns.com", Percent: 8.33},
		{Name: "cache", Percent: 88.89},
	}

	assert.Equal(t, expected, rv)
}

func TestClient_TopClients(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	client := New(Configuration{})
	client.URL = ts.URL
	client.WebPassword = testWebPassword

	rv, err := client.TopClients(5)
	require.NoError(t, err)

	expected := &[]TopClient{
		{Name: "localhost", Requests: 36},
	}

	assert.Equal(t, expected, rv)
}

func TestClient_TopItems(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	client := New(Configuration{})
	client.URL = ts.URL
	client.WebPassword = testWebPassword

	rv, err := client.TopItems(5)
	require.NoError(t, err)

	sort.Slice(rv.TopAds, func(i, j int) bool {
		return rv.TopAds[i].Hits < rv.TopAds[j].Hits
	})
	sort.Slice(rv.TopQueries, func(i, j int) bool {
		return rv.TopQueries[i].Hits < rv.TopQueries[j].Hits
	})

	expected := &TopItems{
		TopQueries: []TopQuery{
			{Name: "api.github.com", Hits: 10},
			{Name: "220.220.67.208.in-addr.arpa", Hits: 11},
			{Name: "222.222.67.208.in-addr.arpa", Hits: 12},
		},
	}

	assert.Equal(t, expected, rv)
}

func Test_NoResponse(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testEmptyData)
			}))
	defer ts.Close()

	client := New(Configuration{})
	client.URL = ts.URL
	client.WebPassword = testWebPassword

	_, err := client.Version()
	assert.Error(t, err)
	_, err = client.SummaryRaw()
	assert.Error(t, err)
	_, err = client.QueryTypes()
	assert.Error(t, err)
	_, err = client.ForwardDestinations()
	assert.Error(t, err)
	_, err = client.TopClients(5)
	assert.Error(t, err)
	_, err = client.TopItems(5)
	assert.Error(t, err)
}

func Test_404(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))
	defer ts.Close()

	client := New(Configuration{})
	client.URL = ts.URL
	client.WebPassword = testWebPassword

	_, err := client.Version()
	assert.Error(t, err)
	_, err = client.SummaryRaw()
	assert.Error(t, err)
	_, err = client.QueryTypes()
	assert.Error(t, err)
	_, err = client.ForwardDestinations()
	assert.Error(t, err)
	_, err = client.TopClients(5)
	assert.Error(t, err)
	_, err = client.TopItems(5)
	assert.Error(t, err)
}

func newTestServer() *httptest.Server {
	handle := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/admin/api.php" {
			w.WriteHeader(http.StatusBadRequest)
		}

		qs := r.URL.Query()
		if len(qs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
		}

		_, isVersion := qs[string(QueryVersion)]
		_, isSummaryRaw := qs[string(QuerySummaryRaw)]

		switch {
		case isVersion:
			_, _ = w.Write(testVersionData)
			return
		case isSummaryRaw:
			_, _ = w.Write(testSummaryRawData)
			return
		}

		authOK := len(qs["auth"]) == 1 && qs["auth"][0] == testWebPassword
		if !authOK {
			_, _ = w.Write(testEmptyData)
			return
		}

		_, isQueryTypes := qs[string(QueryGetQueryTypes)]
		_, isForwardDestinations := qs[string(QueryGetForwardDestinations)]
		_, isTopClients := qs[string(QueryTopClients)]
		_, isTopItems := qs[string(QueryTopItems)]

		switch {
		default:
			_, _ = w.Write(testEmptyData)
		case isQueryTypes:
			_, _ = w.Write(testQueryTypesData)
		case isForwardDestinations:
			_, _ = w.Write(testForwardDestinationsData)
		case isTopClients:
			_, _ = w.Write(testTopClientsData)
		case isTopItems:
			_, _ = w.Write(testTopItemsData)
		}
	}

	return httptest.NewServer(http.HandlerFunc(handle))
}

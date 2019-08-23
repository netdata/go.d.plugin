package zookeeper

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testMntrData, _               = ioutil.ReadFile("testdata/mntr.txt")
	testMntrNotInWhiteListData, _ = ioutil.ReadFile("testdata/mntr_notinwhitelist.txt")
)

func Test_testDataLoad(t *testing.T) {
	assert.NotNil(t, testMntrData)
	assert.NotNil(t, testMntrNotInWhiteListData)
}

func TestNew(t *testing.T) {

}

func TestZookeeper_Init(t *testing.T) {

}

func TestZookeeper_Check(t *testing.T) {

}

func TestZookeeper_Charts(t *testing.T) {

}

func TestZookeeper_Cleanup(t *testing.T) {

}

func TestZookeeper_Collect(t *testing.T) {
	job := New()
	job.zookeeperFetcher = &mockZookeeperFetcher{data: testMntrData}
	require.True(t, job.Init())

	expected := map[string]int64{
		"approximate_data_size":      44,
		"avg_latency":                0,
		"ephemerals_count":           0,
		"max_file_descriptor_count":  1048576,
		"max_latency":                0,
		"min_latency":                0,
		"num_alive_connections":      1,
		"open_file_descriptor_count": 46,
		"outstanding_requests":       0,
		"packets_received":           464,
		"packets_sent":               463,
		"server_state":               4,
		"watch_count":                0,
		"znode_count":                5,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestZookeeper_CollectMntrNotInWhiteList(t *testing.T) {
	job := New()
	job.zookeeperFetcher = &mockZookeeperFetcher{data: testMntrNotInWhiteListData}
	require.True(t, job.Init())

	assert.Nil(t, job.Collect())
}

func TestZookeeper_CollectMntrEmptyResponse(t *testing.T) {
	job := New()
	job.zookeeperFetcher = &mockZookeeperFetcher{}
	require.True(t, job.Init())

	assert.Nil(t, job.Collect())
}

func TestZookeeper_CollectMntrInvalidData(t *testing.T) {
	job := New()
	job.zookeeperFetcher = &mockZookeeperFetcher{data: []byte("hello \nand good buy\n")}
	require.True(t, job.Init())

	assert.Nil(t, job.Collect())
}

func TestZookeeper_CollectMntrReceiveError(t *testing.T) {
	job := New()
	job.zookeeperFetcher = &mockZookeeperFetcher{err: true}
	require.True(t, job.Init())

	assert.Nil(t, job.Collect())
}

type mockZookeeperFetcher struct {
	data []byte
	err  bool
}

func (m mockZookeeperFetcher) fetch(command string) ([]string, error) {
	if m.err {
		return nil, errors.New("mock fetch error")
	}

	var rv []string

	s := bufio.NewScanner(bytes.NewReader(m.data))
	for s.Scan() {
		rv = append(rv, s.Text())
	}

	return rv, nil
}

package docker_engine

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testMetrics, _             = ioutil.ReadFile("testdata/metrics.txt")
	testSwarmManagerMetrics, _ = ioutil.ReadFile("testdata/swarm_manager_metrics.txt")
)

func TestNew(t *testing.T) {
	job := New()

	assert.IsType(t, (*DockerEngine)(nil), job)
	assert.Equal(t, defaultURL, job.UserURL)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
}

func TestDockerEngine_Charts(t *testing.T) { assert.NotNil(t, New().Charts()) }

func TestDockerEngine_Cleanup(t *testing.T) { New().Cleanup() }

func TestDockerEngine_Init(t *testing.T) { assert.True(t, New().Init()) }

func TestDockerEngine_InitNG(t *testing.T) {
	job := New()
	job.UserURL = ""
	assert.False(t, job.Init())
}

func TestDockerEngine_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testMetrics)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL + "/metrics"
	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestDockerEngine_CheckNG(t *testing.T) {
	job := New()
	job.UserURL = "http://127.0.0.1:38001/metrics"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestDockerEngine_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testMetrics)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL + "/metrics"
	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"container_actions_changes":                      1,
		"container_actions_start":                        1,
		"container_actions_commit":                       1,
		"container_actions_delete":                       1,
		"container_actions_create":                       1,
		"container_states_paused":                        11,
		"container_states_running":                       12,
		"container_states_stopped":                       13,
		"builder_fails_dockerfile_empty_error":           4,
		"builder_fails_dockerfile_syntax_error":          5,
		"builder_fails_error_processing_commands_error":  6,
		"builder_fails_build_canceled":                   1,
		"builder_fails_build_target_not_reachable_error": 2,
		"builder_fails_command_not_supported_error":      3,
		"builder_fails_missing_onbuild_arguments_error":  7,
		"builder_fails_unknown_instruction_error":        8,
		"health_checks_failed":                           33,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestDockerEngine_CollectSwarmManager(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testSwarmManagerMetrics)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL + "/metrics"
	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"builder_fails_build_canceled":                   1,
		"builder_fails_build_target_not_reachable_error": 2,
		"builder_fails_command_not_supported_error":      3,
		"builder_fails_dockerfile_empty_error":           4,
		"builder_fails_dockerfile_syntax_error":          5,
		"builder_fails_error_processing_commands_error":  6,
		"builder_fails_missing_onbuild_arguments_error":  7,
		"builder_fails_unknown_instruction_error":        8,
		"container_actions_changes":                      1,
		"container_actions_commit":                       1,
		"container_actions_create":                       1,
		"container_actions_delete":                       1,
		"container_actions_start":                        1,
		"container_states_paused":                        11,
		"container_states_running":                       12,
		"container_states_stopped":                       13,
		"health_checks_failed":                           33,
		"swarm_manager_configs_total":                    1,
		"swarm_manager_leader":                           1,
		"swarm_manager_networks_total":                   3,
		"swarm_manager_nodes_state_disconnected":         1,
		"swarm_manager_nodes_state_down":                 2,
		"swarm_manager_nodes_state_ready":                3,
		"swarm_manager_nodes_state_unknown":              4,
		"swarm_manager_nodes_total":                      10,
		"swarm_manager_secrets_total":                    1,
		"swarm_manager_services_total":                   1,
		"swarm_manager_tasks_state_accepted":             1,
		"swarm_manager_tasks_state_assigned":             2,
		"swarm_manager_tasks_state_complete":             3,
		"swarm_manager_tasks_state_failed":               4,
		"swarm_manager_tasks_state_new":                  5,
		"swarm_manager_tasks_state_orphaned":             6,
		"swarm_manager_tasks_state_pending":              7,
		"swarm_manager_tasks_state_preparing":            8,
		"swarm_manager_tasks_state_ready":                9,
		"swarm_manager_tasks_state_rejected":             10,
		"swarm_manager_tasks_state_remove":               11,
		"swarm_manager_tasks_state_running":              12,
		"swarm_manager_tasks_state_shutdown":             13,
		"swarm_manager_tasks_state_starting":             14,
		"swarm_manager_tasks_total":                      105,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestDockerEngine_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hello and goodbye"))
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL + "/metrics"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestDockerEngine_404(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL + "/metrics"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

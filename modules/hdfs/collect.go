package hdfs

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/netdata/go.d.plugin/pkg/stm"
)

type metrics struct {
	*jvmMetrics `stm:"jvm"`
}

type jvmMetrics struct {
	TagProcessName             string  `json:"tag.ProcessName"`
	TagHostName                string  `json:"tag.Hostname"`
	MemNonHeapUsedM            float64 `stm:"mem_non_heap_used"`
	MemNonHeapCommittedM       float64 `stm:"mem_non_heap_committed"`
	MemNonHeapMaxM             float64 `stm:"mem_non_heap_max"`
	MemHeapUsedM               float64 `stm:"mem_heap_used"`
	MemHeapCommittedM          float64 `stm:"mem_heap_committed"`
	MemHeapMaxM                float64 `stm:"mem_heap_max"`
	MemMaxM                    float64 `stm:"mem_max"`
	GcCount                    float64 `stm:"gc_count"`
	GcTimeMillis               float64 `stm:"gc_time_millis"`
	GcNumWarnThresholdExceeded float64 `stm:"gc_num_warn_threshold_exceeded"`
	GcNumInfoThresholdExceeded float64 `stm:"gc_num_info_threshold_exceeded"`
	GcTotalExtraSleepTime      float64 `stm:"gc_total_extra_sleep_time"`
	ThreadsNew                 float64 `stm:"threads_new"`
	ThreadsRunnable            float64 `stm:"threads_runnable"`
	ThreadsBlocked             float64 `stm:"threads_blocked"`
	ThreadsWaiting             float64 `stm:"threads_waiting"`
	ThreadsTimedWaiting        float64 `stm:"threads_timed_waiting"`
	ThreadsTerminated          float64 `stm:"threads_terminated"`
	LogFatal                   float64 `stm:"log_fatal"`
	LogError                   float64 `stm:"log_error"`
	LogWarn                    float64 `stm:"log_warn"`
	LogInfo                    float64 `stm:"log_info"`
}

type (
	rawData map[string]json.RawMessage
	rawJMX  struct {
		Beans []rawData
	}
)

func (h *HDFS) collect() (map[string]int64, error) {
	var raw rawJMX
	err := h.client.doOKWithDecodeJSON(&raw)
	if err != nil {
		return nil, err
	}

	var mx metrics
	err = h.collectJVM(&mx, raw)
	if err != nil {
		return nil, err
	}

	return stm.ToMap(mx), nil
}

func (h HDFS) collectJVM(mx *metrics, raw rawJMX) error {
	rawJvm := findJvm(raw)
	if rawJvm == nil {
		return errors.New("")
	}

	b, err := json.Marshal(rawJvm)
	if err != nil {
		return errors.New("")
	}

	var jvm jvmMetrics
	err = json.Unmarshal(b, &jvm)
	if err != nil {
		return errors.New("")
	}

	mx.jvmMetrics = &jvm
	return nil
}

func isJvm(data rawData) bool {
	v, ok := data["modelerType"]
	return ok && bytes.Equal(v, []byte("\"JvmMetrics\""))

}

func findJvm(raw rawJMX) (jvm rawData) {
	for _, v := range raw.Beans {
		if !isJvm(v) {
			continue
		}
		return v
	}
	return nil
}

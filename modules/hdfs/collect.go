package hdfs

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/netdata/go.d.plugin/pkg/stm"
)

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

package hdfs

import (
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

func (r rawJMX) isEmpty() bool {
	return len(r.Beans) == 0
}

func (r rawJMX) find(f func(rawData) bool) rawData {
	for _, v := range r.Beans {
		if f(v) {
			return v
		}
	}
	return nil
}

func isJvm(data rawData) bool { return string(data["modelerType"]) == "JvmMetrics" }

func isFsn(data rawData) bool { return string(data["modelerType"]) == "FSNamesystem" }

func (r rawJMX) findJvm() rawData { return r.find(isJvm) }

func (r rawJMX) findFsn() rawData { return r.find(isFsn) }

func (h *HDFS) collect() (map[string]int64, error) {
	var raw rawJMX
	err := h.client.doOKWithDecodeJSON(&raw)
	if err != nil {
		return nil, err
	}

	if raw.isEmpty() {
		return nil, errors.New("empty response")
	}

	var mx metrics

	switch h.nodeType {
	default:
		h.collectUnknownNode(&mx, raw)
	case nameNodeType:
		h.collectNameNode(&mx, raw)
	case dataNodeType:
		h.collectDataNode(&mx, raw)
	}

	return stm.ToMap(mx), nil
}

func (h HDFS) collectNameNode(mx *metrics, raw rawJMX) {
	err := h.collectJVM(mx, raw)
	if err != nil {
		h.Errorf("error on collecting jvm : %v", err)
	}

	err = h.collectFns(mx, raw)
	if err != nil {
		h.Errorf("error on collecting fsn : %v", err)
	}
}

func (h HDFS) collectDataNode(mx *metrics, raw rawJMX) {
	err := h.collectJVM(mx, raw)
	if err != nil {
		h.Errorf("error on collecting jvm : %v", err)
	}
}

func (h HDFS) collectUnknownNode(mx *metrics, raw rawJMX) {
	h.collectDataNode(mx, raw)
}

func (h HDFS) collectJVM(mx *metrics, raw rawJMX) error {
	rawJvm := raw.findJvm()
	if rawJvm == nil {
		return errors.New("couldn't find jvm data")
	}

	b, err := json.Marshal(rawJvm)
	if err != nil {
		return err
	}

	var jvm jvmMetrics
	err = json.Unmarshal(b, &jvm)
	if err != nil {
		return err
	}

	mx.jvmMetrics = &jvm
	return nil
}

func (h HDFS) collectFns(mx *metrics, raw rawJMX) error {
	rawFns := raw.findFsn()
	if rawFns == nil {
		return nil
	}

	b, err := json.Marshal(rawFns)
	if err != nil {
		return err
	}

	var fns fsnNameSystemMetrics
	err = json.Unmarshal(b, &fns)
	if err != nil {
		return err
	}

	mx.fsnNameSystemMetrics = &fns
	return nil
}

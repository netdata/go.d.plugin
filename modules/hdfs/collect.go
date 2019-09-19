package hdfs

import (
	"encoding/json"
	"errors"
	"fmt"

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

func (r rawJMX) findJvm() rawData {
	f := func(data rawData) bool { return string(data["modelerType"]) == "\"JvmMetrics\"" }
	return r.find(f)
}

func (r rawJMX) findFSNameSystem() rawData {
	f := func(data rawData) bool { return string(data["modelerType"]) == "\"FSNamesystem\"" }
	return r.find(f)
}

func (r rawJMX) findFSDatasetState() rawData {
	f := func(data rawData) bool { return string(data["modelerType"]) == "\"FSDatasetState\"" }
	return r.find(f)
}

func (h *HDFS) collect() (map[string]int64, error) {
	var raw rawJMX
	err := h.client.doOKWithDecodeJSON(&raw)
	if err != nil {
		return nil, err
	}

	if raw.isEmpty() {
		return nil, errors.New("empty response")
	}

	mx := h.collectRawJMX(raw)

	return stm.ToMap(mx), nil
}

func (h HDFS) collectRawJMX(raw rawJMX) *metrics {
	var mx metrics
	switch h.nodeType {
	default:
		panic(fmt.Sprintf("unsupported node type : '%s'", h.nodeType))
	case unknownNodeType:
		h.collectUnknownNode(&mx, raw)
	case nameNodeType:
		h.collectNameNode(&mx, raw)
	case dataNodeType:
		h.collectDataNode(&mx, raw)
	}
	return &mx
}

func (h HDFS) collectNameNode(mx *metrics, raw rawJMX) {
	err := h.collectJVM(mx, raw)
	if err != nil {
		h.Errorf("error on collecting jvm : %v", err)
	}

	err = h.collectFSNameSystem(mx, raw)
	if err != nil {
		h.Errorf("error on collecting fs name system : %v", err)
	}
}

func (h HDFS) collectDataNode(mx *metrics, raw rawJMX) {
	err := h.collectJVM(mx, raw)
	if err != nil {
		h.Errorf("error on collecting jvm : %v", err)
	}

	err = h.collectFSDatasetState(mx, raw)
	if err != nil {
		h.Errorf("error on collecting fs dataset state : %v", err)
	}
}

func (h HDFS) collectUnknownNode(mx *metrics, raw rawJMX) {
	err := h.collectJVM(mx, raw)
	if err != nil {
		h.Errorf("error on collecting jvm : %v", err)
	}
}

func (h HDFS) collectJVM(mx *metrics, raw rawJMX) error {
	v := raw.findJvm()
	if v == nil {
		return errors.New("couldn't find jvm data")
	}

	b, err := json.Marshal(v)
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

func (h HDFS) collectFSNameSystem(mx *metrics, raw rawJMX) error {
	v := raw.findFSNameSystem()
	if v == nil {
		return nil
	}

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	var fs fsNameSystemMetrics
	err = json.Unmarshal(b, &fs)
	if err != nil {
		return err
	}

	mx.fsNameSystemMetrics = &fs
	return nil
}

func (h HDFS) collectFSDatasetState(mx *metrics, raw rawJMX) error {
	v := raw.findFSDatasetState()
	if v == nil {
		return nil
	}

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	var fs fsDatasetStateMetrics
	err = json.Unmarshal(b, &fs)
	if err != nil {
		return err
	}

	fs.Used = fs.Capacity - fs.Remaining

	mx.fsDatasetStateMetrics = &fs
	return nil
}

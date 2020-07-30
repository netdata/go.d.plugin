package prometheus

import (
	"sort"
	"strings"
	"unsafe"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/textparse"
)

type (
	// Metric is a pair of label set and value
	Metric struct {
		Labels labels.Labels
		Value  float64
	}
	MetaEntry struct {
		Help string
		Type textparse.MetricType
	}

	Metadata map[string]*MetaEntry
	// Metrics is a list of Metric
	Metrics []Metric
)

// Name the __name__ label value
func (m Metric) Name() string {
	return m.Labels[0].Value
}

// Add appends a metric.
func (m *Metrics) Add(kv Metric) {
	*m = append(*m, kv)
}

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
func (m *Metrics) Reset() {
	*m = (*m)[:0]
}

// Sort sorts data.
func (m Metrics) Sort() {
	sort.Sort(m)
}

// Len returns metric length.
func (m Metrics) Len() int {
	return len(m)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (m Metrics) Less(i, j int) bool {
	return m[i].Name() < m[j].Name()
}

// Swap swaps the elements with indexes i and j.
func (m Metrics) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// FindByName finds metrics where it's __name__ label matches given name.
// It expects the metrics is sorted.
// Complexity: O(log(N))
func (m Metrics) FindByName(name string) Metrics {
	from := sort.Search(len(m), func(i int) bool {
		return m[i].Name() >= name
	})
	if from == len(m) || m[from].Name() != name { // not found
		return Metrics{}
	}
	until := from + 1
	for until < len(m) && m[until].Name() == name {
		until++
	}
	return m[from:until]
}

// FindByNames finds metrics where it's __name__ label matches given any of names.
// It expects the metrics is sorted.
// Complexity: O(log(N))
func (m Metrics) FindByNames(names ...string) Metrics {
	switch len(names) {
	case 0:
		return Metrics{}
	case 1:
		return m.FindByName(names[0])
	}
	var result Metrics
	for _, name := range names {
		result = append(result, m.FindByName(name)...)
	}
	return result
}

// Match match finds metrics where it's label matches given matcher.
// It does NOT expect the metrics is sorted.
// Complexity: O(N)
func (m Metrics) Match(matcher *labels.Matcher) Metrics {
	res := Metrics{}
	for _, kv := range m {
		value := kv.Labels.Get(matcher.Name)
		if matcher.Matches(value) {
			res.Add(kv)
		}
	}
	return res
}

// Max returns the max value.
// It do NOT expect the metrics is sorted.
// Complexity: O(N)
func (m Metrics) Max() float64 {
	switch len(m) {
	case 0:
		return 0
	case 1:
		return m[0].Value
	}
	max := m[0].Value
	for _, kv := range m[1:] {
		if max < kv.Value {
			max = kv.Value
		}
	}
	return max
}

func (m Metadata) Help(name string) string {
	entry, ok := m[name]
	if !ok {
		if strings.HasSuffix(name, "_bucket") {
			return m.Help(name[:len(name)-len("_bucket")])
		}
		return ""
	}
	return entry.Help
}

func (m Metadata) Type(name string) textparse.MetricType {
	entry, ok := m[name]
	if !ok {
		if strings.HasSuffix(name, "_bucket") {
			return m.Type(name[:len(name)-len("_bucket")])
		}
		return textparse.MetricTypeUnknown
	}
	return entry.Type
}

func (m Metadata) setHelp(metric, help []byte) {
	entry, ok := m[unsafeString(metric)]
	if !ok {
		entry = &MetaEntry{Type: textparse.MetricTypeUnknown}
		m[string(metric)] = entry
	}
	if entry.Help != unsafeString(help) {
		entry.Help = string(help)
	}
}

func (m Metadata) setType(metric []byte, mType textparse.MetricType) {
	entry, ok := m[unsafeString(metric)]
	if !ok {
		entry = &MetaEntry{Type: textparse.MetricTypeUnknown}
		m[string(metric)] = entry
	}
	entry.Type = mType
}

func (m Metadata) reset() {
	for key := range m {
		delete(m, key)
	}
}

func unsafeString(b []byte) string {
	return *((*string)(unsafe.Pointer(&b)))
}

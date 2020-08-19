package selector

import (
	"fmt"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/matcher"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/stretchr/testify/assert"
)

func TestLabelMatcher_Matches(t *testing.T) {

}

func TestParse(t *testing.T) {
	tests := map[string]struct {
		input       string
		expectedSr  Selector
		expectedErr bool
	}{
		"glob format: only metric name": {
			input:      "go_memstats_*",
			expectedSr: mustGlobName("go_memstats_*"),
		},
		"simple patterns format: only metric name": {
			input:      "go_memstats_alloc_bytes !go_memstats_* *",
			expectedSr: mustSPName("go_memstats_alloc_bytes !go_memstats_* *"),
		},
		"string format: metric name with labels": {
			input: fmt.Sprintf(`go_memstats_*{label%s"value"}`, FmtEqual),
			expectedSr: andSelector{
				lhs: mustGlobName("go_memstats_*"),
				rhs: mustString("label", "value"),
			},
		},
		"neg string format: metric name with labels": {
			input: fmt.Sprintf(`go_memstats_*{label%s"value"}`, FmtNegEqual),
			expectedSr: andSelector{
				lhs: mustGlobName("go_memstats_*"),
				rhs: Not(mustString("label", "value")),
			},
		},
		"regexp format: metric name with labels": {
			input: fmt.Sprintf(`go_memstats_*{label%s"valu.+"}`, FmtRegexp),
			expectedSr: andSelector{
				lhs: mustGlobName("go_memstats_*"),
				rhs: mustRegexp("label", "valu.+"),
			},
		},
		"neg regexp format: metric name with labels": {
			input: fmt.Sprintf(`go_memstats_*{label%s"valu.+"}`, FmtNegRegexp),
			expectedSr: andSelector{
				lhs: mustGlobName("go_memstats_*"),
				rhs: Not(mustRegexp("label", "valu.+")),
			},
		},
		"glob format: metric name with labels": {
			input: fmt.Sprintf(`go_memstats_*{label%s"valu*"}`, FmtGlob),
			expectedSr: andSelector{
				lhs: mustGlobName("go_memstats_*"),
				rhs: mustGlob("label", "valu*"),
			},
		},
		"neg glob format: metric name with labels": {
			input: fmt.Sprintf(`go_memstats_*{label%s"valu*"}`, FmtNegGlob),
			expectedSr: andSelector{
				lhs: mustGlobName("go_memstats_*"),
				rhs: Not(mustGlob("label", "valu*")),
			},
		},
		"simple patterns format: metric name with labels": {
			input: fmt.Sprintf(`go_memstats_*{label%s"value !val* *"}`, FmtGlob),
			expectedSr: andSelector{
				lhs: mustGlobName("go_memstats_*"),
				rhs: mustSP("label", "value !val* *"),
			},
		},
		"neg simple patterns format: metric name with labels": {
			input: fmt.Sprintf(`go_memstats_*{label%s"value !val* *"}`, FmtNegGlob),
			expectedSr: andSelector{
				lhs: mustGlobName("go_memstats_*"),
				rhs: Not(mustSP("label", "value !val* *")),
			},
		},
		"metric name with several labels": {
			input: fmt.Sprintf(`go_memstats_*{label1%s"value1",label2%s"value2"}`,
				FmtEqual, FmtEqual),
			expectedSr: andSelector{
				lhs: andSelector{
					lhs: mustGlobName("go_memstats_*"),
					rhs: mustString("label1", "value1"),
				},
				rhs: mustString("label2", "value2"),
			},
		},
		"only labels (unsugar)": {
			input: fmt.Sprintf(`{__name__%s"go_memstats_*",label1%s"value1",label2%s"value2"}`,
				FmtGlob, FmtEqual, FmtEqual),
			expectedSr: andSelector{
				lhs: andSelector{
					lhs: mustGlobName("go_memstats_*"),
					rhs: mustString("label1", "value1"),
				},
				rhs: mustString("label2", "value2"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m, err := Parse(test.input)

			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.expectedSr, m)
			}
		})
	}
}

func mustString(name string, pattern string) Selector {
	return labelSelector{name: name, m: matcher.Must(matcher.NewStringMatcher(pattern, true, true))}
}

func mustRegexp(name string, pattern string) Selector {
	return labelSelector{name: name, m: matcher.Must(matcher.NewRegExpMatcher(pattern))}
}

func mustGlob(name string, pattern string) Selector {
	return labelSelector{name: name, m: matcher.Must(matcher.NewGlobMatcher(pattern))}
}

func mustSP(name string, pattern string) Selector {
	return labelSelector{name: name, m: matcher.Must(matcher.NewSimplePatternsMatcher(pattern))}
}

func mustGlobName(pattern string) Selector { return mustGlob(labels.MetricName, pattern) }
func mustSPName(pattern string) Selector   { return mustSP(labels.MetricName, pattern) }

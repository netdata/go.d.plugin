// SPDX-License-Identifier: GPL-3.0-or-later

package pipeline

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_globAny(t *testing.T) {
	tests := map[string]struct {
		patterns  []string
		value     string
		wantFalse bool
	}{
		"one param, matches": {
			patterns: []string{"*"},
			value:    "value",
		},
		"one param, matches with *": {
			patterns: []string{"**/value"},
			value:    "/one/two/three/value",
		},
		"one param, not matches": {
			patterns:  []string{"Value"},
			value:     "value",
			wantFalse: true,
		},
		"several params, last one matches": {
			patterns: []string{"not", "matches", "*"},
			value:    "value",
		},
		"several params, no matches": {
			patterns:  []string{"not", "matches", "really"},
			value:     "value",
			wantFalse: true,
		},
	}

	for name, test := range tests {
		name := fmt.Sprintf("name: %s, patterns: '%v', value: '%s'", name, test.patterns, test.value)

		if test.wantFalse {
			assert.Falsef(t, globAny(test.value, test.patterns[0], test.patterns[1:]...), name)
		} else {
			assert.Truef(t, globAny(test.value, test.patterns[0], test.patterns[1:]...), name)
		}
	}
}

func Test_regexpAny(t *testing.T) {
	tests := map[string]struct {
		patterns  []string
		value     string
		wantFalse bool
	}{
		"one param, matches": {
			patterns: []string{"^value$"},
			value:    "value",
		},
		"one param, not matches": {
			patterns:  []string{"^Value$"},
			value:     "value",
			wantFalse: true,
		},
		"several params, last one matches": {
			patterns: []string{"not", "matches", "va[lue]{3}"},
			value:    "value",
		},
		"several params, no matches": {
			patterns:  []string{"not", "matches", "val[^l]ue"},
			value:     "value",
			wantFalse: true,
		},
	}

	for name, test := range tests {
		name := fmt.Sprintf("name: %s, patterns: '%v', value: '%s'", name, test.patterns, test.value)

		if test.wantFalse {
			assert.Falsef(t, regexpAny(test.value, test.patterns[0], test.patterns[1:]...), name)
		} else {
			assert.Truef(t, regexpAny(test.value, test.patterns[0], test.patterns[1:]...), name)
		}
	}
}

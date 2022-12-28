// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMongo_Init(t *testing.T) {
	tests := map[string]struct {
		config  Config
		success bool
	}{
		"success on default config": {
			success: true,
			config:  New().Config,
		},
		"fails on unset 'address'": {
			success: true,
			config: Config{
				URI:     "mongodb://localhost:27017",
				Timeout: 10,
			},
		},
		"fails on invalid port": {
			success: false,
			config: Config{
				URI:     "",
				Timeout: 0,
			},
		},
	}

	msg := "Init() result does not match Init()"

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := New()
			m.Config = test.config
			assert.Equal(t, test.success, m.Init(), msg)
		})
	}
}

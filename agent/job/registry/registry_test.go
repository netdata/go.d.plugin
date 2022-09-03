// SPDX-License-Identifier: GPL-3.0-or-later

package registry

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileLockRegistry(t *testing.T) {
	assert.NotNil(t, NewFileLockRegistry(""))
}

func TestFileLockRegistry_Register(t *testing.T) {
	tests := map[string]func(t *testing.T, dir string){
		"register a lock": func(t *testing.T, dir string) {
			reg := NewFileLockRegistry(dir)

			ok, err := reg.Register("name")
			assert.True(t, ok)
			assert.NoError(t, err)
		},
		"register the same lock twice": func(t *testing.T, dir string) {
			reg := NewFileLockRegistry(dir)

			ok, err := reg.Register("name")
			require.True(t, ok)
			require.NoError(t, err)

			ok, err = reg.Register("name")
			assert.True(t, ok)
			assert.NoError(t, err)
		},
		"failed to register locked by other process lock": func(t *testing.T, dir string) {
			reg1 := NewFileLockRegistry(dir)
			reg2 := NewFileLockRegistry(dir)

			ok, err := reg1.Register("name")
			require.True(t, ok)
			require.NoError(t, err)

			ok, err = reg2.Register("name")
			assert.False(t, ok)
			assert.NoError(t, err)
		},
		"failed to register because a directory doesnt exist": func(t *testing.T, dir string) {
			reg := NewFileLockRegistry(dir + dir)

			ok, err := reg.Register("name")
			assert.False(t, ok)
			assert.Error(t, err)
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dir, err := os.MkdirTemp(os.TempDir(), "netdata-go-test-file-lock-registry")
			require.NoError(t, err)
			defer func() { require.NoError(t, os.RemoveAll(dir)) }()

			test(t, dir)
		})
	}
}

func TestFileLockRegistry_Unregister(t *testing.T) {
	tests := map[string]func(t *testing.T, dir string){
		"unregister a lock": func(t *testing.T, dir string) {
			reg := NewFileLockRegistry(dir)

			ok, err := reg.Register("name")
			require.True(t, ok)
			require.NoError(t, err)

			assert.NoError(t, reg.Unregister("name"))
		},
		"unregister not registered lock": func(t *testing.T, dir string) {
			reg := NewFileLockRegistry(dir)

			assert.NoError(t, reg.Unregister("name"))
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dir, err := os.MkdirTemp(os.TempDir(), "netdata-go-test-file-lock-registry")
			require.NoError(t, err)
			defer func() { require.NoError(t, os.RemoveAll(dir)) }()

			test(t, dir)
		})
	}
}

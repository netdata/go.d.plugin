// SPDX-License-Identifier: GPL-3.0-or-later

package functions

import (
	"context"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	mgr := NewManager()

	assert.NotNilf(t, mgr.Input, "Input")
	assert.NotNilf(t, mgr.FunctionRegistry, "FunctionRegistry")
}

func TestManager_Register(t *testing.T) {
	type testInputFn struct {
		name    string
		invalid bool
	}
	tests := map[string]struct {
		input    []testInputFn
		expected []string
	}{
		"valid registration": {
			input: []testInputFn{
				{name: "fn1"},
				{name: "fn2"},
			},
			expected: []string{"fn1", "fn2"},
		},
		"registration with duplicates": {
			input: []testInputFn{
				{name: "fn1"},
				{name: "fn2"},
				{name: "fn1"},
			},
			expected: []string{"fn1", "fn2"},
		},
		"registration with nil functions": {
			input: []testInputFn{
				{name: "fn1"},
				{name: "fn2", invalid: true},
			},
			expected: []string{"fn1"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mgr := NewManager()

			for _, v := range test.input {
				if v.invalid {
					mgr.Register(v.name, nil)
				} else {
					mgr.Register(v.name, func(Function) {})
				}
			}

			var got []string
			for name := range mgr.FunctionRegistry {
				got = append(got, name)
			}
			sort.Strings(got)
			sort.Strings(test.expected)

			assert.Equal(t, test.expected, got)
		})
	}
}

func TestManager_Run(t *testing.T) {
	tests := map[string]struct {
		register []string
		input    string
		expected []Function
	}{
		"valid function: single": {
			register: []string{"fn1"},
			input: `
FUNCTION UID 1 "fn1 arg1 arg2"
`,
			expected: []Function{
				{
					key:     "FUNCTION",
					UID:     "UID",
					Timeout: time.Second,
					Name:    "fn1",
					Args:    []string{"arg1", "arg2"},
					Payload: nil,
				},
			},
		},
		"valid function: multiple": {
			register: []string{"fn1", "fn2"},
			input: `
FUNCTION UID 1 "fn1 arg1 arg2"
FUNCTION UID 1 "fn2 arg1 arg2"
`,
			expected: []Function{
				{
					key:     "FUNCTION",
					UID:     "UID",
					Timeout: time.Second,
					Name:    "fn1",
					Args:    []string{"arg1", "arg2"},
					Payload: nil,
				},
				{
					key:     "FUNCTION",
					UID:     "UID",
					Timeout: time.Second,
					Name:    "fn2",
					Args:    []string{"arg1", "arg2"},
					Payload: nil,
				},
			},
		},
		"valid function: single with payload": {
			register: []string{"fn1", "fn2"},
			input: `
FUNCTION_PAYLOAD UID 1 "fn1 arg1 arg2"
payload line1
payload line2
FUNCTION_PAYLOAD_END
`,
			expected: []Function{
				{
					key:     "FUNCTION_PAYLOAD",
					UID:     "UID",
					Timeout: time.Second,
					Name:    "fn1",
					Args:    []string{"arg1", "arg2"},
					Payload: []byte("payload line1\npayload line2"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mgr := NewManager()

			mgr.Input = strings.NewReader(test.input)

			mock := &mockFunctionExecutor{}
			for _, v := range test.register {
				mgr.Register(v, mock.execute)
			}

			runTime := time.Second * 5
			ctx, cancel := context.WithTimeout(context.Background(), runTime)
			defer cancel()

			done := make(chan struct{})

			go func() { defer close(done); mgr.Run(ctx) }()

			timeout := runTime + time.Second*2
			tk := time.NewTimer(timeout)
			defer tk.Stop()

			select {
			case <-done:
				assert.Equal(t, test.expected, mock.executed)
			case <-tk.C:
				t.Errorf("timed out afteter %s", timeout)
			}
		})
	}
}

type mockFunctionExecutor struct {
	executed []Function
}

func (m *mockFunctionExecutor) execute(fn Function) {
	m.executed = append(m.executed, fn)
}

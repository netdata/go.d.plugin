// SPDX-License-Identifier: GPL-3.0-or-later

package docker

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocker_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"default config": {
			wantFail: false,
			config:   New().Config,
		},
		"unset 'address'": {
			wantFail: false,
			config: Config{
				Address: "",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := New()
			d.Config = test.config

			if test.wantFail {
				assert.False(t, d.Init())
			} else {
				assert.True(t, d.Init())
			}
		})
	}
}

func TestDocker_Charts(t *testing.T) {
	assert.Equal(t, len(charts), len(*New().Charts()))
}

func TestDocker_Cleanup(t *testing.T) {
	tests := map[string]struct {
		prepare   func(d *Docker)
		wantClose bool
	}{
		"after New": {
			wantClose: false,
			prepare:   func(d *Docker) {},
		},
		"after Init": {
			wantClose: false,
			prepare:   func(d *Docker) { d.Init() },
		},
		"after Check": {
			wantClose: true,
			prepare:   func(d *Docker) { d.Init(); d.Check() },
		},
		"after Collect": {
			wantClose: true,
			prepare:   func(d *Docker) { d.Init(); d.Collect() },
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mockClient{}
			d := prepareDockerWithMock(m)
			test.prepare(d)

			require.NotPanics(t, d.Cleanup)

			if test.wantClose {
				assert.True(t, m.closeCalled)
			} else {
				assert.False(t, m.closeCalled)
			}
		})
	}

}

func TestDocker_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func() *Docker
		wantFail bool
	}{
		"success when no errors on all calls": {
			wantFail: false,
			prepare:  func() *Docker { return prepareDockerWithMock(&mockClient{}) },
		},
		"fail when error on creating docker client": {
			wantFail: true,
			prepare:  func() *Docker { return prepareDockerWithMock(nil) },
		},
		"fail when error on DiskUsage()": {
			wantFail: true,
			prepare:  func() *Docker { return prepareDockerWithMock(&mockClient{errOnInfo: true}) },
		},
		"fail when error on ContainerList()": {
			wantFail: true,
			prepare:  func() *Docker { return prepareDockerWithMock(&mockClient{errOnContainerList: true}) },
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := test.prepare()

			require.True(t, d.Init())

			if test.wantFail {
				assert.False(t, d.Check())
			} else {
				assert.True(t, d.Check())
			}
		})
	}
}

func TestDocker_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare  func() *Docker
		expected map[string]int64
	}{
		"success when no errors on all calls": {
			prepare: func() *Docker { return prepareDockerWithMock(&mockClient{}) },
			expected: map[string]int64{
				"healthy_containers":   2,
				"images_active":        1,
				"images_dangling":      1,
				"images_size":          300,
				"paused_containers":    5,
				"running_containers":   4,
				"exited_containers":    6,
				"unhealthy_containers": 3,
			},
		},
		"fail when error on creating docker client": {
			prepare:  func() *Docker { return prepareDockerWithMock(nil) },
			expected: nil,
		},
		"fail when error on Info()": {
			prepare:  func() *Docker { return prepareDockerWithMock(&mockClient{errOnInfo: true}) },
			expected: nil,
		},
		"fail when error on ContainerList()": {
			prepare:  func() *Docker { return prepareDockerWithMock(&mockClient{errOnContainerList: true}) },
			expected: nil,
		},
		"fail when error on ImageList()": {
			prepare:  func() *Docker { return prepareDockerWithMock(&mockClient{errOnImageList: true}) },
			expected: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := test.prepare()

			require.True(t, d.Init())
			_ = d.Check()

			mx := d.Collect()

			assert.Equal(t, test.expected, mx)
		})
	}
}

func prepareDockerWithMock(m *mockClient) *Docker {
	d := New()
	if m == nil {
		d.newClient = func(_ Config) (dockerClient, error) { return nil, errors.New("mock.newClient() error") }
	} else {
		d.newClient = func(_ Config) (dockerClient, error) { return m, nil }
	}
	return d
}

type mockClient struct {
	errOnInfo          bool
	errOnContainerList bool
	errOnImageList     bool
	closeCalled        bool
}

func (m *mockClient) Info(_ context.Context) (types.Info, error) {
	if m.errOnInfo {
		return types.Info{}, errors.New("mockClient.Info() error")
	}

	return types.Info{
		ContainersRunning: 4,
		ContainersPaused:  5,
		ContainersStopped: 6,
	}, nil
}

func (m *mockClient) ContainerList(_ context.Context, opts types.ContainerListOptions) ([]types.Container, error) {
	if m.errOnContainerList {
		return nil, errors.New("mockClient.ContainerList() error")
	}

	v := opts.Filters.Get("health")

	if len(v) == 0 {
		return nil, errors.New("mockClient.ContainerList() error (expect 'health' filter)")
	}

	switch v[0] {
	case "healthy":
		return []types.Container{{}, {}}, nil
	case "unhealthy":
		return []types.Container{{}, {}, {}}, nil
	default:
		return nil, nil
	}
}

func (m *mockClient) ImageList(_ context.Context, _ types.ImageListOptions) ([]types.ImageSummary, error) {
	if m.errOnImageList {
		return nil, errors.New("mockClient.ImageList() error")
	}

	return []types.ImageSummary{
		{
			Containers: 0,
			Size:       100,
		},
		{
			Containers: 1,
			Size:       200,
		},
	}, nil
}

func (m *mockClient) NegotiateAPIVersion(_ context.Context) {}

func (m *mockClient) Close() error {
	m.closeCalled = true
	return nil
}

// SPDX-License-Identifier: GPL-3.0-or-later

package vnode

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/netdata/go.d.plugin/logger"

	"gopkg.in/yaml.v2"
)

var Disabled = false // TODO: remove after Netdata v1.39.0. Fix for "from source" stable-channel installations.

func NewRegistry(confDir string) *Registry {
	r := &Registry{
		confDir: confDir,
		nodes:   make(map[string]*VirtualNode),
		Logger:  logger.New("vnode", "registry"),
	}

	r.readConfDir()

	return r
}

type VirtualNode struct {
	GUID     string            `yaml:"guid"`
	Hostname string            `yaml:"hostname"`
	Labels   map[string]string `yaml:"labels"`
}

type Registry struct {
	confDir string
	nodes   map[string]*VirtualNode
	*logger.Logger
}

func (r *Registry) Len() int {
	return len(r.nodes)
}

func (r *Registry) Lookup(key string) (*VirtualNode, bool) {
	v, ok := r.nodes[key]
	return v, ok
}

func (r *Registry) readConfDir() {
	_ = filepath.WalkDir(r.confDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			r.Warning(err)
			return nil
		}

		if !d.Type().IsRegular() || !isConfigFile(path) {
			return nil
		}

		var cfg []VirtualNode
		if err := loadYAML(&cfg, path); err != nil {
			r.Warning(err)
			return nil
		}

		for _, v := range cfg {
			if v.Hostname == "" || v.GUID == "" {
				r.Warningf("skipping virtual node '%+v': some required fields are missing (%s)", v, path)
				continue
			}
			if _, ok := r.nodes[v.Hostname]; ok {
				r.Warningf("skipping virtual node '%+v': duplicate node (%s)", v, path)
				continue
			}
			v := v
			r.Debugf("adding virtual node'%+v' (%s)", v, path)
			r.nodes[v.Hostname] = &v
		}
		return nil
	})
}

func isConfigFile(path string) bool {
	switch filepath.Ext(path) {
	case ".yaml", ".yml", ".conf":
		return true
	default:
		return false
	}
}

func loadYAML(conf interface{}, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	if err := yaml.NewDecoder(f).Decode(conf); err != nil && err != io.EOF {
		return err
	}
	return nil
}

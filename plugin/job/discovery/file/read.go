package file

import (
	"context"
	"os"
	"path/filepath"

	"github.com/netdata/go.d.plugin/pkg/logger"
	"github.com/netdata/go.d.plugin/plugin/job/confgroup"
)

type (
	staticConfig struct {
		confgroup.Default `yaml:",inline"`
		Jobs              []confgroup.Config `yaml:"jobs"`
	}
	sdConfig []confgroup.Config
)

type Reader struct {
	reg   confgroup.Registry
	paths []string
	*logger.Logger
}

func NewReader(reg confgroup.Registry, paths []string) *Reader {
	return &Reader{
		reg:    reg,
		paths:  paths,
		Logger: logger.New("discovery", "file reader"),
	}
}

func (r Reader) String() string {
	return "file reader"
}

func (r Reader) Run(ctx context.Context, in chan<- []*confgroup.Group) {
	r.Info("instance is started")
	defer func() { r.Info("instance is stopped") }()

	select {
	case <-ctx.Done():
	case in <- r.groups():
	}
	close(in)
}

func (r Reader) groups() (groups []*confgroup.Group) {
	for _, pattern := range r.paths {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		for _, path := range matches {
			if fi, err := os.Stat(path); err != nil || !fi.Mode().IsRegular() {
				continue
			}

			if group, err := parse(r.reg, path); err != nil {
				r.Warningf("parse '%s': %v", path, err)
			} else if group == nil {
				groups = append(groups, &confgroup.Group{Source: path})
			} else {
				groups = append(groups, group)
			}
		}
	}
	for _, group := range groups {
		for _, cfg := range group.Configs {
			cfg.SetSource(group.Source)
			cfg.SetProvider("file reader")
		}
	}
	return groups
}

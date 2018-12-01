package multipath

import (
	"fmt"
	"os"
	"path"

	"github.com/netdata/go.d.plugin/logger"

	"github.com/mitchellh/go-homedir"
)

var log = logger.New("multipath", "")

// MultiPath multi-paths
type MultiPath []string

// New New multi-paths
func New(paths ...string) MultiPath {
	set := map[string]bool{}
	path := make(MultiPath, 0, len(paths))
	for _, dir := range paths {
		if dir == "" {
			continue
		}
		if d, err := homedir.Expand(dir); err != nil {
			dir = d
		}
		if !set[dir] {
			path = append(path, dir)
			set[dir] = true
		}
	}
	return path
}

// Find find a file in given paths
func (p MultiPath) Find(filename string) (string, error) {
	for _, dir := range p {
		if dir == "" {
			continue
		}
		if d, err := homedir.Expand(dir); err != nil {
			dir = d
		}
		file := path.Join(dir, filename)
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return file, nil
		}
	}
	return "", fmt.Errorf("cannot find file %s in any of %v", filename, p)
}

// MustFind find a file in given paths. if not find, exit program
func (p MultiPath) MustFind(filename string) string {
	file, err := p.Find(filename)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return file
}

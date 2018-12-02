package multipath

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

// MultiPath multi-paths
type MultiPath []string

// New New multi-paths
func New(paths ...string) MultiPath {
	set := map[string]bool{}
	mPath := make(MultiPath, 0, len(paths))

	for _, dir := range paths {
		if dir == "" {
			continue
		}
		if d, err := homedir.Expand(dir); err != nil {
			dir = d
		}
		if !set[dir] {
			mPath = append(mPath, dir)
			set[dir] = true
		}
	}

	return mPath
}

// Find find a file in given paths
func (p MultiPath) Find(filename string) (string, error) {
	for _, dir := range p {
		file := path.Join(dir, filename)
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return file, nil
		}
	}
	return "", fmt.Errorf("can't find '%s' in %v", filename, p)
}

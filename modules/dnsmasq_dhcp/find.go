package dnsmasq_dhcp

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type (
	extension string

	extensions []extension

	configDir struct {
		path    string
		include extensions
		exclude extensions
	}
)

func (e extension) match(filename string) bool {
	return strings.HasSuffix(filename, string(e))
}

func (es extensions) match(filename string) bool {
	for _, e := range es {
		if e.match(filename) {
			return true
		}
	}
	return false
}

func parseConfDir(confDirStr string) configDir {
	// # Include all the files in a directory except those ending in .bak
	//#conf-dir=/etc/dnsmasq.d,.bak
	//# Include all files in a directory which end in .conf
	//#conf-dir=/etc/dnsmasq.d/,*.conf

	parts := strings.Split(confDirStr, ",")
	cd := configDir{path: parts[0]}

	for _, arg := range parts[1:] {
		arg = strings.TrimSpace(arg)
		if strings.HasPrefix(arg, "*") {
			cd.include = append(cd.include, extension(arg[1:]))
		} else {
			cd.exclude = append(cd.exclude, extension(arg))
		}
	}
	return cd
}

func (cd configDir) isValidFilename(filename string) bool {
	switch {
	default:
		return true
	case strings.HasPrefix(filename, "."):
	case strings.HasPrefix(filename, "~"):
	case strings.HasPrefix(filename, "#") && strings.HasSuffix(filename, "#"):
	}
	return false
}

func (cd configDir) match(filename string) bool {
	switch {
	default:
		return true
	case !cd.isValidFilename(filename):
	case len(cd.include) > 0 && !cd.include.match(filename):
	case cd.exclude.match(filename):
	}
	return false
}

func (cd configDir) findConfigs() ([]string, error) {
	fis, err := ioutil.ReadDir(cd.path)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, fi := range fis {
		if !fi.Mode().IsRegular() || !cd.match(fi.Name()) {
			continue
		}
		files = append(files, filepath.Join(cd.path, fi.Name()))
	}
	return files, nil
}

func openFile(filepath string) (f *os.File, err error) {
	defer func() {
		if err != nil && f != nil {
			_ = f.Close()
		}
	}()

	f, err = os.Open(filepath)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("'%s' is not a regular file", filepath)
	}
	return f, nil
}

type (
	configOption struct {
		key, value string
	}

	configFile struct {
		path    string
		options []configOption
	}
)

func (cf *configFile) get(name string) []string {
	var options []string
	for _, o := range cf.options {
		if o.key != name {
			continue
		}
		options = append(options, o.value)
	}
	return options
}

func parseConfFile(filename string) (*configFile, error) {
	f, err := openFile(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cf := configFile{path: filename}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}

		if !strings.Contains(line, "=") {
			continue
		}

		line = strings.ReplaceAll(line, " ", "")
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}

		cf.options = append(cf.options, configOption{key: parts[0], value: parts[1]})
	}
	return &cf, nil
}

type ConfigFinder struct {
	entryConfig    string
	entryDir       string
	visitedConfigs map[string]bool
	visitedDirs    map[string]bool
}

func (f *ConfigFinder) find() []*configFile {
	f.visitedConfigs = make(map[string]bool)
	f.visitedDirs = make(map[string]bool)

	configs := f.recursiveFind(f.entryConfig)

	for _, file := range f.entryDirConfigs() {
		configs = append(configs, f.recursiveFind(file)...)
	}
	return configs
}

func (f ConfigFinder) entryDirConfigs() []string {
	if f.entryDir == "" {
		return nil
	}
	files, err := parseConfDir(f.entryDir).findConfigs()
	if err != nil {
		return nil
	}
	return files
}

func (f *ConfigFinder) recursiveFind(filename string) (configs []*configFile) {
	if f.visitedConfigs[filename] {
		return nil
	}

	config, err := parseConfFile(filename)
	if err != nil {
		return nil
	}

	files, dirs := config.get("conf-file"), config.get("conf-dir")

	f.visitedConfigs[filename] = true
	configs = append(configs, config)

	for _, file := range files {
		configs = append(configs, f.recursiveFind(file)...)
	}

	for _, dir := range dirs {
		if dir == "" {
			continue
		}

		d := parseConfDir(dir)

		if f.visitedDirs[d.path] {
			continue
		}
		f.visitedDirs[d.path] = true

		files, err = d.findConfigs()
		if err != nil {
			continue
		}

		for _, file := range files {
			configs = append(configs, f.recursiveFind(file)...)
		}
	}
	return configs
}

func findConfigurationFiles(entryConfig string, entryDir string) []*configFile {
	cf := ConfigFinder{
		entryConfig: entryConfig,
		entryDir:    entryDir,
	}
	return cf.find()
}

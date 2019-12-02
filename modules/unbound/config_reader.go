package unbound

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

/*
Files can be included using the include: directive.
It can appear anywhere, it accepts a single file name as argument.

Processing continues as if the text  from  the included file was copied into the config file at that point.

If also using chroot, using full path names for the included files works, relative pathnames for the included names
work if the directory where the daemon is started equals its chroot/working directory or is specified before
the include statement with  directory:  dir. Wildcards can be used to include multiple files, see glob(7).


Unbound stop processing and exits on any error:
 - syntax error
 - recursive include
*/

var configAttributes = []string{
	"include",
	"statistics-cumulative",
	"control-enable",
	"control-interface",
	"control-port",
	"control-use-cert",
	"control-key-file",
	"control-cert-file",
}

func neededAttribute(line string) bool {
	for _, v := range configAttributes {
		if strings.HasPrefix(line, v) {
			return true
		}
	}
	return false
}

func readConfigFile(entryPath string) ([][]string, error) {
	return configFileReader{visited: make(map[string]bool)}.read(entryPath)
}

type configFileReader struct {
	visited map[string]bool
}

func (c configFileReader) read(filename string) ([][]string, error) {
	var attrs [][]string
	if c.visited[filename] {
		return nil, fmt.Errorf("file '%s' has been visited", filename)
	}
	c.visited[filename] = true

	f, err := c.open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || !neededAttribute(line) {
			continue
		}

		key, value, err := c.parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("file '%s', error on parsing line '%s': %v", filename, line, err)
		}

		if key != "include" {
			attrs = append(attrs, []string{key, value})
			continue
		}

		attrs, err := c.handleInclude(value)
		if err != nil {
			return nil, err
		}
		for _, v := range attrs {
			attrs = append(attrs, v)
		}
	}
	return attrs, nil
}

func (c configFileReader) handleInclude(value string) ([][]string, error) {
	if !isGlobPattern(value) {
		return c.read(value)
	}
	filenames, err := filepath.Glob(value)
	if err != nil {
		return nil, err
	}
	var attrs [][]string
	for _, name := range filenames {
		attrs, err := c.read(name)
		if err != nil {
			return nil, err
		}
		for _, v := range attrs {
			attrs = append(attrs, v)
		}
	}
	return attrs, nil
}

func (configFileReader) parseLine(line string) (string, string, error) {
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return "", "", errors.New("bad syntax")
	}
	key, value := cleanKeyValue(parts[0], parts[1])
	return key, value, nil
}

func (configFileReader) open(filename string) (*os.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("'%s' is not a regular file", filename)
	}
	return f, nil
}

func isGlobPattern(value string) bool {
	magicChars := `*?[`
	if runtime.GOOS != "windows" {
		magicChars = `*?[\`
	}
	return strings.ContainsAny(value, magicChars)
}

func cleanKeyValue(key, value string) (string, string) {
	if i := strings.IndexByte(value, '#'); i > 0 {
		value = value[:i-1]
	}
	key = strings.TrimSpace(key)
	value = strings.Trim(strings.TrimSpace(value), "\"'")
	return key, value
}

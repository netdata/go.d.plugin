package tail

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
)

var maxFails = 10

var (
	ErrNotInited   = errors.New("not inited")
	ErrGlob        = errors.New("glob returns an empty slice")
	ErrBadFile     = errors.New("not a readable file")
	SizeNotChanged = errors.New("size not changed")
)

func New() *Tail {
	return new(Tail)
}

type Tail struct {
	Path string // used as pattern
	path string

	fails int
	pos   int64
}

func (t *Tail) Init(path string) error {
	t.Path = path
	err := t.globPath()
	if err != nil {
		return err
	}

	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	t.pos = fi.Size()
	t.path = path
	return nil
}

func (t *Tail) Tail() (io.ReadCloser, error) {
	if t.path == "" {
		return nil, ErrNotInited
	}

	if t.fails > maxFails {
		err := t.globPath()
		if err != nil {
			return nil, err
		}
	}

	fi, err := os.Stat(t.path)

	if err != nil {
		t.fails++
		return nil, err
	}

	if fi.Size() < t.pos {
		t.pos = 0
	}

	if fi.Size() == t.pos {
		t.fails = 0
		return nil, SizeNotChanged
	}

	f, err := os.Open(t.path)
	if err != nil {
		t.fails++
		return nil, err
	}

	_, err = f.Seek(t.pos, io.SeekStart)
	if err != nil {
		t.fails++
		return nil, err
	}

	t.fails = 0
	t.pos = fi.Size()
	return f, nil
}

func (t *Tail) globPath() error {
	v, err := filepath.Glob(t.Path)
	if err != nil {
		return err
	}

	if len(v) == 0 {
		return ErrGlob
	}

	sort.Strings(v)
	p := v[len(v)-1]

	if !isReadableFile(p) {
		return ErrBadFile
	}

	t.path = p
	return nil
}

func isReadableFile(name string) bool {
	_, e1 := os.Open(name)
	v, e2 := os.Stat(name)
	return e1 == nil && e2 == nil && v.Mode().IsRegular()
}

// ReadLastLine returns the last line of the file and any read error encountered.
func ReadLastLine(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readLastLine(f)
}

func readLastLine(f io.ReadSeeker) ([]byte, error) {
	b := make([]byte, 1)
	var c int

	f.Seek(0, io.SeekEnd)
	for {
		if v, err := f.Seek(-2, io.SeekCurrent); err != nil {
			return nil, err
		} else if v == 0 {
			c += 2
			break
		}

		if _, err := f.Read(b); err != nil {
			return nil, err
		}
		c++
		if b[0] == '\n' {
			break
		}
		continue
	}

	rv := make([]byte, c)
	f.Read(rv)
	return rv, nil
}

package simpletail

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

var maxFails = 10

var (
	// ErrNotInitialized ErrNotInitialized
	ErrNotInitialized = errors.New("not initialized")
	// ErrGlob ErrGlob
	ErrGlob = errors.New("glob returns an empty slice")
	// ErrBadFile ErrBadFile
	ErrBadFile = errors.New("not a readable file")
	// SizeNotChanged SizeNotChanged
	SizeNotChanged = errors.New("size not changed")
)

func New(path string) *Tail {
	return &Tail{
		Path: path,
	}
}

type Tail struct {
	Path string // pattern
	path string

	fails int
	pos   int64
}

func (t *Tail) Init() error {
	err := t.globPath()
	if err != nil {
		return err
	}

	fi, err := os.Stat(t.Path)
	if err != nil {
		return err
	}

	t.pos = fi.Size()
	return nil
}

func (t *Tail) Tail() (io.ReadCloser, error) {
	if t.path == "" {
		return nil, ErrNotInitialized
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

	if fi.Size() == t.pos {
		return nil, SizeNotChanged
	}

	if fi.Size() < t.pos {
		t.pos = 0
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
	// TODO: this is wrong
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
	defer func() {
		_ = f.Close()
	}()
	return readLastLine(f)
}

func readLastLine(f io.ReadSeeker) ([]byte, error) {
	_, _ = f.Seek(0, io.SeekEnd)
	b := make([]byte, 1)

	for b[0] != '\n' {
		if v, err := f.Seek(-2, io.SeekCurrent); err != nil {
			return nil, err
		} else if v == 0 {
			break
		}
		_, _ = f.Read(b)
	}

	line, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return line, nil
}

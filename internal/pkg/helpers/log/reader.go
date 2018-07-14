package log

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
)

var runFindPath = 10

var (
	ErrNotStarted     = errors.New("not started")
	ErrSizeNotChanged = errors.New("size not changed")
)

// TODO the overall design looks bad. But it it works with minimum memory allocation even if log file is huge
// and all it's content have to be read. Anyway should be fixed.

// NewReader returns init'ed instance of Reader and error if any.
func NewReader(path string) (*Reader, error) {
	l := Reader{
		Path: path,
		err:  make(chan error),
		data: make(chan string),
		do:   make(chan bool),
	}
	if err := l.init(); err != nil {
		return nil, err
	}
	return &l, nil
}

type Reader struct {
	Path    string // used as pattern
	path    string
	pos     int64
	fails   int
	started bool

	err  chan error
	data chan string
	do   chan bool
}

// GetRawData returns error from worker if any or data channel
func (r *Reader) GetRawData() (chan string, error) {
	if !r.started {
		return nil, ErrNotStarted
	}
	r.do <- true
	err := <-r.err

	if err != nil {
		if err == ErrSizeNotChanged {
			r.fails = 0
			return nil, err
		}
		r.fails++
		return nil, err
	}
	r.fails = 0
	return r.data, nil
}

// init finds the exact value of path and starts worker
func (r *Reader) init() error {
	p, err := findPath(r.Path)
	if err != nil {
		return err
	}
	r.path = p
	r.seekEnd()
	go worker(r)
	r.started = true
	return nil
}

func (r *Reader) seekEnd() {
	fi, _ := os.Stat(r.path)
	r.pos = fi.Size()
}

// ------------------------------------------------------------------------------
func worker(l *Reader) {
	for {
		<-l.do
		if l.fails > runFindPath {
			if v, err := findPath(l.Path); err != nil {
				l.err <- err
				continue
			} else {
				l.path = v
				l.seekEnd()
			}
		}
		fi, err := os.Stat(l.path)
		if err != nil {
			l.err <- err
			continue
		}

		if fi.Size() < l.pos {
			l.pos = 0
		} else if fi.Size() == l.pos {
			l.err <- ErrSizeNotChanged
			continue
		}

		f, err := os.Open(l.path)
		if err != nil {
			l.err <- err
			continue
		}

		if _, err := f.Seek(l.pos, io.SeekStart); err != nil {
			l.err <- err
			continue
		}

		s := bufio.NewScanner(f)
		l.err <- nil

		for s.Scan() {
			l.data <- s.Text()
		}
		close(l.data)

		l.data = make(chan string)
		l.pos, _ = f.Seek(0, io.SeekCurrent)
		f.Close()
	}
}

func findPath(path string) (string, error) {
	v, err := filepath.Glob(path)
	if err != nil {
		return "", err
	}
	if len(v) == 0 {
		return "", errors.New("glob failed")
	}
	sort.Strings(v)

	p := v[len(v)-1]

	if !isFile(p) {
		return "", errors.New("not a file")
	}
	if !isReadable(p) {
		return "", errors.New("not readable")
	}
	return p, nil
}

func isFile(name string) bool {
	v, e := os.Stat(name)
	return e == nil && v.Mode().IsRegular()
}

func isReadable(name string) bool {
	_, e := os.Open(name)
	return e == nil
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

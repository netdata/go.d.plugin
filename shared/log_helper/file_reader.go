package log_helper

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
	ErrNotStarted    = errors.New("not started")
	ErrNothingToRead = errors.New("nothing to read")
)

// TODO the overall design looks bad. But it it works with minimum memory allocation even if log file is huge
// and all it's content have to be read. Anyway should be fixed.

// NewFileReader returns init'ed instance of FileReader and error if any.
func NewFileReader(path string) (*FileReader, error) {
	l := FileReader{
		Path: path,
		err:  make(chan error),
		data: make(chan []byte),
		do:   make(chan bool),
	}
	if err := l.init(); err != nil {
		return nil, err
	}
	return &l, nil
}

type FileReader struct {
	Path    string // used as pattern
	path    string
	pos     int64
	fails   int
	started bool

	err  chan error
	data chan []byte
	do   chan bool
}

// GetRawData returns error from worker if any or data channel
func (l *FileReader) GetRawData() (chan []byte, error) {
	if !l.started {
		return nil, ErrNotStarted
	}
	l.do <- true
	err := <-l.err

	if err != nil {
		if err == ErrNothingToRead {
			l.fails = 0
			return nil, err
		}
		l.fails++
		return nil, err
	}
	l.fails = 0
	return l.data, nil
}

// init finds the exact value of path and starts worker
func (l *FileReader) init() error {
	p, err := findPath(l.Path)
	if err != nil {
		return err
	}
	l.path = p
	l.seekEnd()
	go worker(l)
	l.started = true
	return nil
}

func (l *FileReader) seekEnd() {
	fi, _ := os.Stat(l.path)
	l.pos = fi.Size()
}

func worker(l *FileReader) {
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
			l.err <- ErrNothingToRead
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
			l.data <- s.Bytes()
		}
		close(l.data)

		l.data = make(chan []byte)
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

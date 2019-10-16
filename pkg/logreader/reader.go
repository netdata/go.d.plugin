package logreader

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/netdata/go-orchestrator/logger"
)

const (
	maxEOF = 600
)

var (
	ErrNoMatchedFile = errors.New("no matched files")
)

// Reader is a log rotate aware Reader
type Reader struct {
	file          *os.File
	path          string
	excludePath   string
	eofCounter    int
	continuousEOF int
	log           *logger.Logger
}

// Open a file and seek to end of the file.
// path: shell file name pattern
// excludePath: shell file name pattern
func Open(path string, excludePath string, log *logger.Logger) (*Reader, error) {
	var err error
	if path, err = filepath.Abs(path); err != nil {
		return nil, err
	}
	if _, err = filepath.Match(path, "/"); err != nil {
		return nil, fmt.Errorf("bad path syntax: %q", path)
	}
	if _, err = filepath.Match(excludePath, "/"); err != nil {
		return nil, fmt.Errorf("bad exclude_path syntax: %q", path)
	}
	r := &Reader{
		path:        path,
		excludePath: excludePath,
		log:         log,
	}

	if err = r.open(); err != nil {
		return nil, err
	}
	return r, nil
}

// CurrentFilename get current opened file name
func (f *Reader) CurrentFilename() string {
	return f.file.Name()
}

func (f *Reader) open() error {
	path := f.findFile()
	f.log.Debug("open log file: ", path)
	if path == "" {
		return ErrNoMatchedFile
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if _, err = file.Seek(stat.Size(), io.SeekStart); err != nil {
		return err
	}
	f.file = file
	return nil
}

func (f *Reader) Read(p []byte) (n int, err error) {
	n, err = f.file.Read(p)
	if err == io.EOF {
		f.eofCounter++
		f.continuousEOF++
		if f.eofCounter >= maxEOF && f.continuousEOF >= 2 {
			if err2 := f.reopen(); err2 != nil {
				err = err2
			}
		}
	} else {
		f.continuousEOF = 0
	}
	return
}

func (f *Reader) Close() (err error) {
	if f == nil || f.file == nil {
		return
	}
	f.log.Debug("close log file: ", f.file.Name())
	err = f.file.Close()
	f.file = nil
	f.eofCounter = 0
	return
}

func (f *Reader) reopen() error {
	f.Close()
	return f.open()
}

func (f *Reader) findFile() string {
	files, _ := filepath.Glob(f.path)
	if len(files) == 0 {
		return ""
	}

	if f.excludePath != "" {
		files2 := make([]string, 0, len(files))
		for _, file := range files {
			if ok, _ := filepath.Match(f.excludePath, file); !ok {
				files2 = append(files2, file)
			}
		}
		if len(files2) == 0 {
			return ""
		}
		files = files2
	}

	sort.Strings(files)
	for i := len(files) - 1; i >= 0; i-- {
		stat, err := os.Stat(files[i])
		if err == nil && !stat.IsDir() {
			return files[i]
		}
	}
	return ""
}

const DefaultMaxLineWidth = 4 * 1024 // assume disk block size is 4K

var ErrTooLongLine = errors.New("too long line")

// ReadLastLine returns the last line of the file and any read error encountered.
// It expect last line width <= maxLineWidth.
// If maxLineWidth <= 0, it defaults to DefaultMaxLineWidth.
func ReadLastLine(filename string, maxLineWidth int64) ([]byte, error) {
	if maxLineWidth <= 0 {
		maxLineWidth = DefaultMaxLineWidth
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	stat, _ := f.Stat()
	endPos := stat.Size()
	if endPos == 0 {
		return []byte{}, nil
	}
	startPos := endPos - maxLineWidth
	if startPos < 0 {
		startPos = 0
	}
	buf := make([]byte, endPos-startPos)
	n, err := f.ReadAt(buf, startPos)
	if err != nil {
		return nil, err
	}
	lnPos := 0
	foundLn := false
	for i := n - 2; i >= 0; i-- {
		ch := buf[i]
		if ch == '\n' {
			foundLn = true
			lnPos = i
			break
		}
	}
	if foundLn {
		return buf[lnPos+1 : n], nil
	}
	if startPos == 0 {
		return buf[0:n], nil
	}

	return nil, ErrTooLongLine
}

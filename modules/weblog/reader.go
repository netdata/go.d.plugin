package weblog

import (
	"io"
	"os"
	"path/filepath"
	"sort"
)

type Reader struct {
	path        string
	currentFile string
	file        *os.File
	size        int64
}

func NewReader(path string) *Reader {
	return &Reader{
		path: path,
	}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if r == nil {
		return 0, os.ErrInvalid
	}
	if r.file == nil {
		if err = r.open(); err != nil {
			return
		}
	}
	n, err = r.file.Read(p)
	r.size += int64(n)
	return
}

func (r *Reader) Reload() error {
	if r == nil || r.file == nil {
		return os.ErrInvalid
	}
	stat, err := os.Stat(r.currentFile)
	if err != nil {
		return err
	}
	if stat.Size() >= r.size { // no need to reload
		return nil
	}
	r.Close()
	r.open()
	return nil
}

func (r *Reader) Close() error {
	if r == nil {
		return os.ErrInvalid
	}
	err := r.file.Close()
	r.file = nil
	return err
}

func (r *Reader) open() error {
	r.file.Close()
	file, err := r.filename()
	if err != nil {
		return err
	}
	stat, err := os.Stat(file)
	if err != nil {
		return err
	}
	r.file, err = os.Open(file)
	if err != nil {
		return err
	}
	r.size = stat.Size()
	r.currentFile = file
	r.file.Seek(stat.Size(), io.SeekStart)
	return nil
}

func (r *Reader) filename() (string, error) {
	files, err := filepath.Glob(r.path)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", nil
	}
	sort.Strings(files)
	return files[0], nil
}

package weblog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

const (
	maxEOF = 2
)

type FileManager struct {
	path        string
	excludePath string
}

type File struct {
	file       *os.File
	filePath   string
	manager    *FileManager
	eofCounter int
}

func (f *File) Read(p []byte) (n int, err error) {
	n, err = f.file.Read(p)
	if err == io.EOF {
		f.eofCounter++
		if f.eofCounter >= maxEOF {
			f.Reopen()
		}
	} else {
		f.eofCounter = 0
	}
	return
}

func (f *File) Close() (err error) {
	if f == nil {
		return
	}
	err = f.file.Close()
	f.file = nil
	f.eofCounter = 0
	return
}

func (f *File) Reopen() error {
	path := f.manager.FindFile()
	if path == "" {
		return errNoMatchedFile
	}
	if f.filePath != path {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		f.Close()
		f.file = file
		f.filePath = path
	}
	f.file.Seek(0, io.SeekEnd)
	return nil
}

var (
	errNoMatchedFile = errors.New("no matched files")
)

func NewFileManager(path string, excludePath string) (*FileManager, error) {
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
	return &FileManager{
		path:        path,
		excludePath: excludePath,
	}, nil
}

func (f *FileManager) OpenFile() (*File, error) {
	file := &File{
		manager: f,
	}
	if err := file.Reopen(); err != nil {
		return nil, err
	}
	return file, nil
}

func (f *FileManager) FindFile() string {
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

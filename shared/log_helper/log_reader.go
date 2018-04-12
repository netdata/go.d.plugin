package log_helper

import (
	"io"
	"os"
)

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

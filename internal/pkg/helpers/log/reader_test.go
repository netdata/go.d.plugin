package log

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewFileReader(t *testing.T) {
	_, err := NewReader("this will fail")

	if err == nil {
		t.Error("expected error, but got nil")
	}

	tmp, err := ioutil.TempFile("", "temp-")

	if err != nil {
		t.Fatal("could not create temporary file")
	}

	defer os.Remove(tmp.Name())

	v, err := NewReader(tmp.Name())

	if err != nil {
		t.Fatalf("could not create Reader: %s", err)
	}

	if !v.started {
		t.Error("Reader not started")
	}

	if v.path == "" {
		t.Error("expected not empty Reader path")
	}
}

func TestFileReader_GetRawData(t *testing.T) {
	tmp, err := ioutil.TempFile("", "temp-")

	if err != nil {
		t.Fatal("could not create temporary file")
	}

	defer os.Remove(tmp.Name())

	v, err := NewReader(tmp.Name())

	if err != nil {
		t.Fatalf("could not create Reader: %s", err)
	}

	data, err := v.GetRawData()
	if err == nil {
		t.Fatal("expected error, but got nil")
	}

	if err != ErrSizeNotChanged {
		t.Errorf("expected '%s' error, but got %s", ErrSizeNotChanged, err)
	}
	lines := [...]string{1: "first", 2: "second", 3: "third", 4: "fourth"}

	tmp.Write([]byte(lines[1] + "\n"))
	data, err = v.GetRawData()

	if err != nil {
		t.Fatalf("expected nil, but got %s", err)
	}
	rv := <-data
	if rv != lines[1] {
		t.Fatalf("expected %s, but got %s", lines[1], rv)
	}

	tmp.Write([]byte(lines[2] + "\n"))
	data, err = v.GetRawData()

	if err != nil {
		t.Fatalf("expected nil, but got %s", err)
	}
	rv = <-data
	if rv != lines[2] {
		t.Fatalf("expected %s, but got %s", lines[2], rv)
	}

	data, err = v.GetRawData()
	if err == nil {
		t.Fatal("expected error, but got nil")
	}

	if err != ErrSizeNotChanged {
		t.Errorf("expected '%s' error, but got %s", ErrSizeNotChanged, err)
	}

	tmp.Write([]byte(lines[3] + "\n" + lines[4] + "\n"))

	data, err = v.GetRawData()

	if err != nil {
		t.Fatalf("expected nil, but got %s", err)
	}
	rv = <-data
	rv += <-data
	if rv != lines[3]+lines[4] {
		t.Fatalf("expected %s, but got %s", lines[3]+lines[4], rv)
	}
}

func TestReadLastLine(t *testing.T) {
	empty := ""
	input := bytes.NewReader([]byte(empty))

	if _, err := readLastLine(input); err == nil {
		t.Error("expected error, but got nil")
	}

	oneLine := "first line"
	input = bytes.NewReader([]byte(oneLine))

	if v, err := readLastLine(input); err != nil {
		t.Errorf("expected nil, but got %s", err)
	} else if string(v) != oneLine {
		t.Errorf("expected %s, but got %s", oneLine, string(v))
	}

	oneLine = "世界"
	input = bytes.NewReader([]byte(oneLine))

	if v, err := readLastLine(input); err != nil {
		t.Errorf("expected nil, but got %s", err)
	} else if string(v) != oneLine {
		t.Errorf("expected %s, but got %s", oneLine, string(v))
	}

	secondLine := "second line"
	input = bytes.NewReader([]byte(oneLine + "\n" + secondLine))

	if v, err := readLastLine(input); err != nil {
		t.Errorf("expected nil, but got %s", err)
	} else if string(v) != secondLine {
		t.Errorf("expected %s, but got %s", secondLine, string(v))
	}

	secondLine = "世界"
	input = bytes.NewReader([]byte(oneLine + "\n" + secondLine))

	if v, err := readLastLine(input); err != nil {
		t.Errorf("expected nil, but got %s", err)
	} else if string(v) != secondLine {
		t.Errorf("expected %s, but got %s", secondLine, string(v))
	}
}

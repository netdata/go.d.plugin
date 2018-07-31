package tail

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	tail := New()
	if _, ok := interface{}(tail).(*Tail); !ok {
		t.Error("expected *Tail type")
	}
}

func TestTail_Init(t *testing.T) {
	tail := New()

	err := tail.Init("fail")
	if err == nil {
		t.Error("expected error, but got nil")
	}
	if err != ErrGlob {
		t.Errorf("expected %s, but got %s", ErrGlob, err)
	}

	_, err = tail.Tail()
	if err != ErrNotInited {
		t.Errorf("expected %s, but got %s", ErrNotInited, err)
	}

	tmp, err := ioutil.TempFile("", "temp-")
	if err != nil {
		t.Fatal("could not create temporary file")
	}
	defer os.Remove(tmp.Name())

	err = tail.Init(tmp.Name())

	if err != nil {
		t.Fatalf("tail init: %s", err)
	}

	if tail.path == "" {
		t.Error("expected not empty tail 'path'")
	}
}

func TestTail_Tail(t *testing.T) {
	tail := New()

	tmp, err := ioutil.TempFile("", "temp-")
	if err != nil {
		t.Fatal("could not create temporary file")
	}
	defer os.Remove(tmp.Name())
	tail.Init(tmp.Name())

	_, err = tail.Tail()
	if err == nil {
		t.Fatal("expected error, but got nil")
	}

	if err != SizeNotChanged {
		t.Errorf("expected %s error, but got %s", SizeNotChanged, err)
	}

	w := func() {
		tmp.WriteString("Donatello\n")
		tmp.WriteString("Leonardo\n")
		tmp.WriteString("Michelangelo\n")
		tmp.WriteString("Raphael\n")
	}
	w()

	data, err := tail.Tail()

	if err != nil {
		t.Fatalf("excpected nil, but got %s", err)
	}

	if _, ok := interface{}(data).(io.ReadCloser); !ok {
		t.Error("excpected io.ReadCloser type")
	}

	var c int
	s := bufio.NewScanner(data)
	for s.Scan() {
		c++
	}
	if c != 4 {
		t.Errorf("excepted 4, but got %d", c)
	}

	w()
	w()
	data, err = tail.Tail()
	if err != nil {
		t.Fatalf("excpected nil, but got %s", err)
	}

	c = 0
	s = bufio.NewScanner(data)
	for s.Scan() {
		c++
	}
	if c != 8 {
		t.Errorf("excepted 4, but got %d", c)
	}

	tail.pos = 999
	data, err = tail.Tail()
	if err != nil {
		t.Fatalf("excpected nil, but got %s", err)
	}
	c = 0
	s = bufio.NewScanner(data)
	for s.Scan() {
		c++
	}
	if c != 12 {
		t.Errorf("excepted 12, but got %d", c)
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

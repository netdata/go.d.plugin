package tail

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.IsType(t, (*Tail)(nil), New(""))
}

func TestTail_Init(t *testing.T) {
	tail := New("fail")

	err := tail.Init()
	if assert.Error(t, err) {
		assert.Equal(t, err, ErrGlob)
	}

	_, err = tail.Tail()

	if assert.Error(t, err) {
		assert.Equal(t, err, ErrNotInited)
	}

	tmp, err := ioutil.TempFile("", "temp-")

	assert.Nilf(t, err, "could not create temporary file")
	defer func() {
		_ = os.Remove(tmp.Name())
	}()

	tail = New(tmp.Name())
	err = tail.Init()

	assert.Nil(t, err)
	assert.NotEqual(t, tail.path, "")
}

func TestTail_Tail(t *testing.T) {
	tmp, err := ioutil.TempFile("", "temp-")
	assert.Nil(t, err)

	defer func() {
		_ = os.Remove(tmp.Name())
	}()

	tail := New(tmp.Name())
	_ = tail.Init()

	_, err = tail.Tail()
	if assert.Error(t, err) {
		assert.Equal(t, err, SizeNotChanged)
	}

	w := func() {
		_, _ = tmp.WriteString("Donatello\n")
		_, _ = tmp.WriteString("Leonardo\n")
		_, _ = tmp.WriteString("Michelangelo\n")
		_, _ = tmp.WriteString("Raphael\n")
	}
	w()

	data, err := tail.Tail()

	assert.Nil(t, err)

	assert.Implements(t, (*io.ReadCloser)(nil), data)

	var c int
	s := bufio.NewScanner(data)
	for s.Scan() {
		c++
	}

	assert.Equal(t, c, 4)

	w()
	w()

	data, err = tail.Tail()

	assert.Nil(t, err)

	c = 0
	s = bufio.NewScanner(data)

	for s.Scan() {
		c++
	}

	assert.Equal(t, c, 8)

	tail.pos = 999
	data, err = tail.Tail()

	assert.Nil(t, err)

	c = 0
	s = bufio.NewScanner(data)
	for s.Scan() {
		c++
	}

	assert.Equal(t, c, 12)
}

func TestReadLastLine(t *testing.T) {
	empty := ""
	input := bytes.NewReader([]byte(empty))
	_, err := readLastLine(input)

	assert.Error(t, err)

	oneLine := "first line"
	input = bytes.NewReader([]byte(oneLine))
	v, err := readLastLine(input)

	assert.Nil(t, err)
	assert.Equal(t, string(v), oneLine)

	oneLine = "世界"
	input = bytes.NewReader([]byte(oneLine))
	v, err = readLastLine(input)

	assert.Nil(t, err)
	assert.Equal(t, string(v), oneLine)

	secondLine := "second line"
	input = bytes.NewReader([]byte(oneLine + "\n" + secondLine))
	v, err = readLastLine(input)

	assert.Nil(t, err)
	assert.Equal(t, string(v), secondLine)

	secondLine = "世界"
	input = bytes.NewReader([]byte(oneLine + "\n" + secondLine))
	v, err = readLastLine(input)

	assert.Nil(t, err)
	assert.Equal(t, string(v), secondLine)
}

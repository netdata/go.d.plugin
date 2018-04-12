package log_helper

import (
	"bytes"
	"testing"
)

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

package shared

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

type Duration struct {
	Duration time.Duration
}

func (m *Duration) UnmarshalTOML(input []byte) error {
	s := string(bytes.Trim(input, "'\""))

	if v, err := time.ParseDuration(s); err == nil {
		m.Duration = v
		return nil
	}
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		m.Duration = time.Duration(v) * time.Second
		return nil
	}
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		m.Duration = time.Duration(v) * time.Second
		return nil
	}
	return fmt.Errorf("unparsable duration format '%s'", s)
}

func (m *Duration) ConvertTo(to time.Duration) int {
	return int(int64(m.Duration) / (int64(to) / int64(time.Nanosecond)))
}

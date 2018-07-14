package utils

import (
	"fmt"
	"strconv"
	"time"
)

type Duration struct {
	Duration time.Duration
}

func (m *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string

	if err := unmarshal(&s); err != nil {
		return err
	}

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

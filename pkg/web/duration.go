package web

import (
	"fmt"
	"strconv"
	"time"
)

// Duration is a time.Duration wrapper
type Duration struct {
	Duration time.Duration
}

// UnmarshalYAML implements yaml.Unmarshaler
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

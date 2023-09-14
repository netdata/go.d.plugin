// SPDX-License-Identifier: GPL-3.0-or-later

package functions

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Function struct {
	key     string
	UID     string
	Timeout time.Duration
	Name    string
	Args    []string
	Payload []byte
}

func (f *Function) String() string {
	return fmt.Sprintf("KEY: %s, UID: %s, TIMEOUT: %s, FUNCTION: %s, ARGS: %v, PAYLOAD: %s",
		f.key, f.UID, f.Timeout, f.Name, f.Args, string(f.Payload))
}

func parseFunctionString(s string) (*Function, error) {
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = ' '

	parts, err := r.Read()
	if err != nil {
		return nil, err
	}
	if len(parts) != 4 {
		return nil, fmt.Errorf("unexpected number of words: want 4, got %d (%v)", len(parts), parts)
	}

	cmd := strings.Split(parts[3], " ")

	timeout, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, err
	}

	fn := &Function{
		key:     parts[0],
		UID:     parts[1],
		Timeout: time.Duration(timeout) * time.Second,
		Name:    cmd[0],
		Args:    cmd[1:],
	}

	return fn, nil
}

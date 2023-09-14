// SPDX-License-Identifier: GPL-3.0-or-later

package functions

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"strings"

	"github.com/netdata/go.d.plugin/logger"
)

const (
	apiKeyFunction           = "FUNCTION"
	apiKeyFunctionPayload    = "FUNCTION_PAYLOAD"
	apiKeyFunctionPayloadEnd = "FUNCTION_PAYLOAD_END"
)

func NewManager() *Manager {
	return &Manager{
		Logger:           logger.New("functions", "manager"),
		Input:            os.Stdin,
		FunctionRegistry: make(map[string]func(Function)),
	}
}

type Manager struct {
	*logger.Logger

	Input            io.Reader
	FunctionRegistry map[string]func(Function)
}

func (m *Manager) Register(name string, fn func(Function)) {
	m.Infof("FUNCTION REGISTRATION: '%s'", name)
	m.FunctionRegistry[name] = fn
}

func (m *Manager) Run(ctx context.Context) {
	m.Info("instance is started")
	defer func() { m.Info("instance is stopped") }()

	go func() {
		sc := bufio.NewScanner(m.Input)

		for sc.Scan() {
			text := sc.Text()

			var fn *Function
			var err error

			switch {
			case strings.HasPrefix(text, apiKeyFunction+" "):
				fn, err = m.parseFunction(text)
			case strings.HasPrefix(text, apiKeyFunctionPayload+" "):
				fn, err = m.parseFunctionWithPayload(text, sc)
			default:
				m.Warningf("unexpected line: '%s'", text)
				continue
			}

			if err != nil {
				m.Warningf("parse function: %v ('%s')", err, text)
				continue
			}

			m.runFunction(fn)
		}
	}()

	<-ctx.Done()
}

func (m *Manager) parseFunction(text string) (*Function, error) {
	return parseFunctionString(text)
}

func (m *Manager) parseFunctionWithPayload(text string, sc *bufio.Scanner) (*Function, error) {
	fn, err := parseFunctionString(text)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	for sc.Scan() && sc.Text() != apiKeyFunctionPayloadEnd {
		buf.WriteString(sc.Text() + "\n")
	}

	fn.Payload = append(fn.Payload, buf.Bytes()...)

	return fn, nil
}

func (m *Manager) runFunction(fn *Function) {
	m.Infof("FUNCTION: '%s'", fn.String())

	regFn, ok := m.FunctionRegistry[fn.Name]
	if !ok {
		m.Infof("UNREGISTERED FUNCTION: '%s'", fn.Name)
		return
	}

	regFn(*fn)
}
